package ap

import "github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"

type (
	AppProxyAPI interface {
		AppProxyClusters() APClusterAPI
		AppProxyGitSources() APGitSourceAPI
		AppProxyIsc() APIscAPI
		GitIntegrations() APGitIntegrationAPI
		VersionInfo() APVersionInfoAPI
	}

	apImpl struct {
		client *client.CfClient
	}
)

func NewAppProxyClient(c *client.CfClient) AppProxyAPI {
	return &apImpl{client: c}
}

func (ap *apImpl) AppProxyClusters() APClusterAPI {
	return &apCluster{client: ap.client}
}

func (ap *apImpl) GitIntegrations() APGitIntegrationAPI {
	return &apGitIntegration{client: ap.client}
}

func (ap *apImpl) VersionInfo() APVersionInfoAPI {
	return &apVersionInfo{client: ap.client}
}

func (ap *apImpl) AppProxyGitSources() APGitSourceAPI {
	return &apGitSource{client: ap.client}
}

func (ap *apImpl) AppProxyIsc() APIscAPI {
	return &apIsc{client: ap.client}
}
