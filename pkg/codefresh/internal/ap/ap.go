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
	return &apCluster{client: ap.client}
}

func (ap *apImpl) GitIntegration() GitIntegrationAPI {
	return &apGitIntegration{client: ap.client}
}

func (ap *apImpl) VersionInfo() VersionInfoAPI {
	return &apVersionInfo{client: ap.client}
}

func (ap *apImpl) GitSource() GitSourceAPI {
	return &apGitSource{client: ap.client}
}

func (ap *apImpl) ISC() IscAPI {
	return &apIsc{client: ap.client}
}
