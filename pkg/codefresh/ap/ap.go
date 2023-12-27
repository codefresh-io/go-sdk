package ap

import "github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"

type (
	AppProxyAPI interface {
		Cluster() ClusterAPI
		GitSource() GitSourceAPI
		ISC() IscAPI
		GitIntegration() GitIntegrationAPI
		VersionInfo() VersionInfoAPI
	}

	apImpl struct {
		client *client.CfClient
	}
)

func NewAppProxyClient(c *client.CfClient) AppProxyAPI {
	return &apImpl{client: c}
}

func (ap *apImpl) Cluster() ClusterAPI {
	return &cluster{client: ap.client}
}

func (ap *apImpl) GitIntegration() GitIntegrationAPI {
	return &gitIntegration{client: ap.client}
}

func (ap *apImpl) VersionInfo() VersionInfoAPI {
	return &versionInfo{client: ap.client}
}

func (ap *apImpl) GitSource() GitSourceAPI {
	return &gitSource{client: ap.client}
}

func (ap *apImpl) ISC() IscAPI {
	return &isc{client: ap.client}
}
