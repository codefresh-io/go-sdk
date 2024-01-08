package rest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
)

type (
	ClusterAPI interface {
		GetAccountClusters() ([]ClusterMinified, error)
		GetClusterCredentialsByAccountId(selector string) (*Cluster, error)
	}

	cluster struct {
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

func (p *cluster) GetAccountClusters() ([]ClusterMinified, error) {
	res, err := p.client.RestAPI(context.TODO(), &client.RequestOptions{
		Method: "GET",
		Path:   "/api/clusters",
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting account cluster list: %w", err)
	}

	result := make([]ClusterMinified, 0)
	return result, json.Unmarshal(res, &result)
}

func (p *cluster) GetClusterCredentialsByAccountId(selector string) (*Cluster, error) {
	res, err := p.client.RestAPI(context.TODO(), &client.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/clusters/%s/credentials", selector),
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting account cluster credentials: %w", err)
	}

	result := &Cluster{}
	return result, json.Unmarshal(res, result)
}
