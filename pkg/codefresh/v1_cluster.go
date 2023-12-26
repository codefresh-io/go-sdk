package codefresh

import "fmt"

type (
	V1ClusterAPI interface {
		GetAccountClusters() ([]ClusterMinified, error)
		GetClusterCredentialsByAccountId(selector string) (*Cluster, error)
	}

	v1Cluster struct {
		codefresh *codefresh
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
	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/clusters"),
		method: "GET",
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	r := make([]ClusterMinified, 0)
	err = p.codefresh.decodeResponseInto(resp, r)
	return r, err
}

func (p *v1Cluster) GetClusterCredentialsByAccountId(selector string) (*Cluster, error) {
	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/clusters/%s/credentials", selector),
		method: "GET",
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	r := &Cluster{}
	err = p.codefresh.decodeResponseInto(resp, r)
	return r, err
}
