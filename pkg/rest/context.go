package rest

import (
	"encoding/json"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
)

type (
	ContextAPI interface {
		GetDefaultGitContext() (*ContextPayload, error)
		GetGitContextByName(name string) (*ContextPayload, error)
		GetGitContexts() ([]ContextPayload, error)
	}

	v1Context struct {
		client *client.CfClient
	}

	ContextPayload struct {
		Metadata struct {
			Name string `json:"name"`
		}
		Spec struct {
			Type string `json:"type"`
			Data struct {
				Auth struct {
					Type     string `json:"type"`
					Username string `json:"username"`
					Password string `json:"password"`
					ApiHost  string `json:"apiHost"`
					// for gitlab
					ApiURL         string `json:"apiURL"`
					ApiPathPrefix  string `json:"apiPathPrefix"`
					SshPrivateKey  string `json:"sshPrivateKey"`
					AppId          string `json:"appId"`
					InstallationId string `json:"installationId"`
					PrivateKey     string `json:"privateKey"`
				} `json:"auth"`
			} `json:"data"`
		} `json:"spec"`
	}

	GitContextsQs struct {
		Type    []string `url:"type"`
		Decrypt string   `url:"decrypt"`
	}
)

func (c v1Context) GetDefaultGitContext() (*ContextPayload, error) {
	res, err := c.client.RestAPI(nil, &client.RequestOptions{
		Method: "GET",
		Path:   "/api/contexts/git/default",
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting default git context: %w", err)
	}

	result := &ContextPayload{}
	return result, json.Unmarshal(res, result)
}

func (c v1Context) GetGitContextByName(name string) (*ContextPayload, error) {
	res, err := c.client.RestAPI(nil, &client.RequestOptions{
		Method: "GET",
		Path:   "/api/contexts/" + name,
		Query: map[string]any{
			"decrypt": "true",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting git context by name: %w", err)
	}

	result := &ContextPayload{}
	return result, json.Unmarshal(res, result)
}

func (c v1Context) GetGitContexts() ([]ContextPayload, error) {
	res, err := c.client.RestAPI(nil, &client.RequestOptions{
		Method: "GET",
		Path:   "/api/contexts",
		Query: map[string]any{
			"type":    []string{"git.github", "git.gitlab", "git.github-app"},
			"decrypt": "true",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting git context list: %w", err)
	}

	result := make([]ContextPayload, 0)
	return result, json.Unmarshal(res, &result)
}
