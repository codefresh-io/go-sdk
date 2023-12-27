package v1

import (
	"encoding/json"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	V1ClusterAPI interface {
		GetAccountClusters() ([]ClusterMinified, error)
		GetClusterCredentialsByAccountId(selector string) (*Cluster, error)
	}

	v1Cluster struct {
		client *client.CfClient
	}

	Cluster struct {
		Auth struct {
			Bearer string
		} `json:"auth"`
		Ca  string `json:"ca"`
		Url string `json:"url"`
	}

	ClusterMinified struct {
		Cluster struct {
			Name string `json:"name"`
		} `json:"cluster"`
		BehindFirewall bool   `json:"behindFirewall"`
		Selector       string `json:"selector"`
		Provider       string `json:"provider"`
	}
)

func (p *v1Cluster) GetAccountClusters() ([]ClusterMinified, error) {
	resp, err := p.client.RestAPI(nil, &client.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/clusters"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting account cluster list: %w", err)
	}

	result := make([]ClusterMinified, 0)
	return result, json.Unmarshal(resp, &result)
}

func (p *v1Cluster) GetClusterCredentialsByAccountId(selector string) (*Cluster, error) {
	resp, err := p.client.RestAPI(nil, &client.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/clusters/%s/credentials", selector),
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting account cluster credentials: %w", err)
	}

	result := &Cluster{}
	return result, json.Unmarshal(resp, result)
}
