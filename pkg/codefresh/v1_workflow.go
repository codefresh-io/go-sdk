package codefresh

import (
	"errors"
	"fmt"
	"time"
)

type (
	V1WorkflowAPI interface {
		Get(string) (*Workflow, error)
		WaitForStatus(string, string, time.Duration, time.Duration) error
	}

	v1Workflow struct {
		codefresh *codefresh
	}

	Workflow struct {
		ID                 string    `json:"id"`
		Status             string    `json:"status"`
		UserYamlDescriptor string    `json:"userYamlDescriptor"`
		Progress           string    `json:"progress"`
		Created            time.Time `json:"created"`
		Updated            time.Time `json:"updated"`
		Finished           time.Time `json:"finished"`
	}
)

func (w *v1Workflow) Get(id string) (*Workflow, error) {
	resp, err := w.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/builds/%s", id),
		method: "GET",
	})
	// failed in api call
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	wf := &Workflow{}
	err = w.codefresh.decodeResponseInto(resp, wf)
	if err != nil {
		return nil, err
	}

	return wf, nil
}

func (w *v1Workflow) WaitForStatus(id string, status string, interval time.Duration, timeout time.Duration) error {
	return waitFor(interval, timeout, func() (bool, error) {
		resp, err := w.codefresh.requestAPI(&requestOptions{
			path:   fmt.Sprintf("/api/builds/%s", id),
			method: "GET",
		})
		if err != nil {
			return false, err
		}

		defer resp.Body.Close()
		wf := &Workflow{}
		err = w.codefresh.decodeResponseInto(resp, wf)
		if err != nil {
			return false, err
		}

		if wf.Status == status {
			return true, nil
		}

		return false, nil
	})
}

func waitFor(interval time.Duration, timeout time.Duration, execution func() (bool, error)) error {
	t := time.After(timeout)
	tick := time.Tick(interval)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-t:
			return errors.New("timed out")
		case <-tick:
			ok, err := execution()
			if err != nil {
				return err
			}

			if ok {
				return nil
			}
		}
	}
}
