package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/codefresh-io/go-sdk/pkg/client"
)

type (
	WorkflowAPI interface {
		Get(string) (*Workflow, error)
		WaitForStatus(string, string, time.Duration, time.Duration) error
	}

	workflow struct {
		codefresh *client.CfClient
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

func (w *workflow) Get(id string) (*Workflow, error) {
	res, err := w.codefresh.RestAPI(context.TODO(), &client.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/builds/%s", id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting a workflow: %w", err)
	}

	result := &Workflow{}
	return result, json.Unmarshal(res, result)
}

func (w *workflow) WaitForStatus(id string, status string, interval time.Duration, timeout time.Duration) error {
	return waitFor(interval, timeout, func() (bool, error) {
		res, err := w.Get(id)
		if err != nil {
			return false, err
		}

		if res.Status == status {
			return true, nil
		}

		return false, nil
	})
}

func waitFor(interval time.Duration, timeout time.Duration, execution func() (bool, error)) error {
	t := time.After(timeout)
	ticker := time.NewTicker(interval)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-t:
			return errors.New("timed out")
		case <-ticker.C:
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
