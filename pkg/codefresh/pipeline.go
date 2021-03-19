package codefresh

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type (
	// IPipelineAPI declers Codefresh pipeline API
	IPipelineAPI interface {
		List(qs map[string]string) ([]*Pipeline, error)
		Run(string, *RunOptions) (string, error)
		Create(name string, spec PipelineSpec) (string, error)
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
		Contexts  []interface{} `json:"contexts"`
		Variables []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"variables"`
		Steps  map[string]interface{} `json:"steps"`
		Stages []interface{}          `json:"stages"`
		Mode   string                 `json:"mode"`
	}

	Pipeline struct {
		Metadata PipelineMetadata `json:"metadata"`
		Spec     PipelineSpec     `json:"spec"`
	}

	getPipelineResponse struct {
		Docs  []*Pipeline `json:"docs"`
		Count int         `json:"count"`
	}

	pipeline struct {
		codefresh Codefresh
	}

	RunOptions struct {
		Branch    string
		Variables map[string]string
	}

	CreatePipelineResponse struct {
		Metadata struct {
			Name      string `json:"name"`
			Project   string `json:"project"`
			ProjectID string `json:"projectId"`
			Revision  int    `json:"revision"`
			AccountID string `json:"accountId"`
			Labels    struct {
				Tags []interface{} `json:"tags"`
			} `json:"labels"`
			OriginalYamlString string    `json:"originalYamlString"`
			CreatedAt          time.Time `json:"created_at"`
			UpdatedAt          time.Time `json:"updated_at"`
			ID                 string    `json:"id"`
		} `json:"metadata"`
		Version string `json:"version"`
		Kind    string `json:"kind"`
		Spec    struct {
			Triggers          []interface{} `json:"triggers"`
			Stages            []string      `json:"stages"`
			Variables         []interface{} `json:"variables"`
			Contexts          []interface{} `json:"contexts"`
			TerminationPolicy []interface{} `json:"terminationPolicy"`
			ExternalResources []interface{} `json:"externalResources"`
			Steps             struct {
				Freestyle struct {
					Title            string `json:"title"`
					Type             string `json:"type"`
					WorkingDirectory string `json:"working_directory"`
					Arguments        struct {
						Image    string   `json:"image"`
						Commands []string `json:"commands"`
					} `json:"arguments"`
				} `json:"Freestyle"`
			} `json:"steps"`
		} `json:"spec"`
	}
)

func newPipelineAPI(codefresh Codefresh) IPipelineAPI {
	return &pipeline{codefresh}
}

// Get - returns pipelines from API
func (p *pipeline) List(qs map[string]string) ([]*Pipeline, error) {
	r := &getPipelineResponse{}
	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   "/api/pipelines",
		method: "GET",
		qs:     qs,
	})
	err = p.codefresh.decodeResponseInto(resp, r)
	return r.Docs, err
}

func (p *pipeline) Run(name string, options *RunOptions) (string, error) {
	if options == nil {
		options = &RunOptions{}
	}
	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/pipelines/run/%s", url.PathEscape(name)),
		method: "POST",
		body: map[string]interface{}{
			"branch":    options.Branch,
			"variables": options.Variables,
		},
	})
	if err != nil {
		return "", err
	}
	res, err := p.codefresh.getBodyAsString(resp)
	return strings.Replace(res, "\"", "", -1), err
}

func (p *pipeline) Create(name string, spec PipelineSpec) (string, error) {
	r := &CreatePipelineResponse{}
	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   "/api/pipelines",
		method: "POST",
		body: map[string]interface{}{
			"metadata": name,
			"spec":     spec,
		},
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("non 200 response from create: %v (%v)", resp.Status, resp.StatusCode)
	}
	defer resp.Body.Close()
	err = p.codefresh.decodeResponseInto(resp, r)
	if err != nil {
		return "", err
	}
	return r.Metadata.ID, nil
}
