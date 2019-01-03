package codefresh

import "fmt"

type (
	// IRuntimeEnvironmentAPI declers Codefresh runtime environment API
	IRuntimeEnvironmentAPI interface {
		CreateRuntimeEnvironment(*CreateRuntimeOptions) (*RuntimeEnvironment, error)
	}

	RuntimeEnvironment struct {
		Name string
	}

	CreateRuntimeOptions struct {
		Cluster   string
		Namespace string
		HasAgent  bool
	}

	createRuntimeEnvironmentResponse struct {
		Name string
	}
)

// CreateRuntimeEnvironment - returns pipelines from API
func (c *codefresh) CreateRuntimeEnvironment(opt *CreateRuntimeOptions) (*RuntimeEnvironment, error) {
	// r := &createRuntimeEnvironmentResponse{}
	re := &RuntimeEnvironment{
		Name: fmt.Sprintf("%s/%s", opt.Cluster, opt.Namespace),
	}
	body := map[string]interface{}{
		"clusterName": opt.Cluster,
		"namespace":   opt.Namespace,
	}
	if opt.HasAgent {
		body["agent"] = true
	}
	resp, err := c.requestAPI(&requestOptions{
		path:   "/api/custom_clusters/register",
		method: "POST",
		body:   body,
	})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 400 {
		return re, nil
	} else {
		return nil, fmt.Errorf("Error during runtime environment creation")
	}
}
