package codefresh

import "fmt"

type (
	IClusterAPI interface {
		GetClusterCredentialsByAccountId(selector string) (*Cluster, error)
		GetAccountClusters() ([]*ClusterMinified, error)
	}

	cluster struct {
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

func newClusterAPI(codefresh *codefresh) IClusterAPI {
	return &cluster{codefresh}
}

func (p *cluster) GetClusterCredentialsByAccountId(selector string) (*Cluster, error) {
	r := &Cluster{}
	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/clusters/%s/credentials", selector),
		method: "GET",
	})
	err = p.codefresh.decodeResponseInto(resp, &r)

	defer resp.Body.Close()

	return r, err
}

func (p *cluster) GetAccountClusters() ([]*ClusterMinified, error) {
	r := make([]*ClusterMinified, 0)
	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/clusters"),
		method: "GET",
	})
	err = p.codefresh.decodeResponseInto(resp, &r)

	defer resp.Body.Close()

	return r, err
}
