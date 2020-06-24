package codefresh

import (
	"fmt"
)

type (
	IProgressAPI interface {
		Get(string) (*Progress, error)
	}

	progress struct {
		codefresh Codefresh
	}

	Progress struct {
		ID       string   `json:"id"`
		Status   string   `json:"status"`
		Location Location `json:"location"`
	}

	Location struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	}
)

func newProgressAPI(codefresh Codefresh) IProgressAPI {
	return &progress{codefresh}
}

func (p *progress) Get(id string) (*Progress, error) {
	result := &Progress{}
	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/progress/%s", id),
		method: "GET",
	})
	// failed in api call
	if err != nil {
		return nil, err
	}
	err = p.codefresh.decodeResponseInto(resp, result)
	// failed to decode
	if err != nil {
		return nil, err
	}
	return result, nil
}
