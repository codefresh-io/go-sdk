package codefresh

import (
	"fmt"
)

type (
	V1ProgressAPI interface {
		Get(string) (*Progress, error)
	}

	v1Progress struct {
		codefresh *codefresh
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

func (p *v1Progress) Get(id string) (*Progress, error) {
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
