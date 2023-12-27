package codefresh

import (
	"encoding/json"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	V1ProgressAPI interface {
		Get(string) (*Progress, error)
	}

	v1Progress struct {
		client *client.CfClient
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
	resp, err := p.client.RestAPI(nil, &client.RequestOptions{
		Path:   fmt.Sprintf("/api/progress/%s", id),
		Method: "GET",
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting progress: %w", err)
	}

	result := &Progress{}
	return result, json.Unmarshal(resp, result)
}
