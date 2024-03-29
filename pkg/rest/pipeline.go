package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/codefresh-io/go-sdk/pkg/client"
)

type (
	PipelineAPI interface {
		List(query map[string]string) ([]Pipeline, error)
		Run(string, *RunOptions) (string, error)
	}

	pipeline struct {
		client *client.CfClient
	}

	PipelineMetadata struct {
		Name     string `json:"name"`
		IsPublic bool   `json:"isPublic"`
		Labels   struct {
			Tags []string `json:"tags"`
		} `json:"labels"`
		Deprecate struct {
			ApplicationPort string `json:"applicationPort"`
			RepoPipeline    bool   `json:"repoPipeline"`
		} `json:"deprecate"`
		OriginalYamlString string    `json:"originalYamlString"`
		AccountID          string    `json:"accountId"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		Project            string    `json:"project"`
		ID                 string    `json:"id"`
	}

	PipelineSpec struct {
		Triggers []struct {
			Type     string   `json:"type"`
			Repo     string   `json:"repo"`
			Events   []string `json:"events"`
			Provider string   `json:"provider"`
			Context  string   `json:"context"`
		} `json:"triggers"`
		Contexts  []any `json:"contexts"`
		Variables []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"variables"`
		Steps  map[string]any `json:"steps"`
		Stages []any          `json:"stages"`
		Mode   string         `json:"mode"`
	}

	Pipeline struct {
		Metadata PipelineMetadata `json:"metadata"`
		Spec     PipelineSpec     `json:"spec"`
	}

	getPipelineResponse struct {
		Docs  []Pipeline `json:"docs"`
		Count int        `json:"count"`
	}

	RunOptions struct {
		Branch    string
		Variables map[string]string
	}
)

// Get - returns pipelines from API
func (p *pipeline) List(query map[string]string) ([]Pipeline, error) {
	anyQuery := map[string]any{}
	for k, v := range query {
		anyQuery[k] = v
	}

	res, err := p.client.RestAPI(context.TODO(), &client.RequestOptions{
		Method: "GET",
		Path:   "/api/pipelines",
		Query:  anyQuery,
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting pipeline list: %w", err)
	}

	result := &getPipelineResponse{}
	return result.Docs, json.Unmarshal(res, result)
}

func (p *pipeline) Run(name string, options *RunOptions) (string, error) {
	if options == nil {
		options = &RunOptions{}
	}

	res, err := p.client.RestAPI(context.TODO(), &client.RequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/api/pipelines/run/%s", url.PathEscape(name)),
		Body: map[string]any{
			"branch":    options.Branch,
			"variables": options.Variables,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed running pipeline: %w", err)
	}

	return strings.Replace(string(res), "\"", "", -1), nil
}
