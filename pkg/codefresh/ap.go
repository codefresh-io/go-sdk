package codefresh

type (
	AppProxyAPI interface {
		AppProxyClusters() APClusterAPI
		AppProxyGitSources() APGitSourceAPI
		AppProxyIsc() APIscAPI
		GitIntegrations() APGitIntegrationAPI
		VersionInfo() APVersionInfoAPI
	}

	apImpl struct {
		codefresh *codefresh
	}
)

func newAppProxyClient(c *codefresh) AppProxyAPI {
	return &apImpl{codefresh: c}
}

func (ap *apImpl) AppProxyClusters() APClusterAPI {
	return &apCluster{codefresh: ap.codefresh}
}

func (ap *apImpl) GitIntegrations() APGitIntegrationAPI {
	return &apGitIntegration{codefresh: ap.codefresh}
}

func (ap *apImpl) VersionInfo() APVersionInfoAPI {
	return &apVersionInfo{codefresh: ap.codefresh}
}

func (ap *apImpl) AppProxyGitSources() APGitSourceAPI {
	return &apGitSource{codefresh: ap.codefresh}
}

func (ap *apImpl) AppProxyIsc() APIscAPI {
	return &apIsc{codefresh: ap.codefresh}
}
