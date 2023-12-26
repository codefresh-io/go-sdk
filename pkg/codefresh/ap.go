package codefresh

type (
	AppProxyAPI interface {
		AppProxyClusters() IAppProxyClustersAPI
		AppProxyGitSources() IAppProxyGitSourcesAPI
		AppProxyIsc() IAppProxyIscAPI
		GitIntegrations() IAppProxyGitIntegrationsAPI
		VersionInfo() IAppProxyVersionInfoAPI
	}

	apImpl struct {
		codefresh *codefresh
	}
)

func newAppProxyClient(c *codefresh) AppProxyAPI {
	return &apImpl{codefresh: c}
}

func (ap *apImpl) AppProxyClusters() IAppProxyClustersAPI {
	return newAppProxyClustersAPI(ap.codefresh)
}

func (ap *apImpl) GitIntegrations() IAppProxyGitIntegrationsAPI {
	return newAppProxyGitIntegrationsAPI(ap.codefresh)
}

func (ap *apImpl) VersionInfo() IAppProxyVersionInfoAPI {
	return newAppProxyVersionInfoAPI(ap.codefresh)
}

func (ap *apImpl) AppProxyGitSources() IAppProxyGitSourcesAPI {
	return newAppProxyGitSourcesAPI(ap.codefresh)
}

func (ap *apImpl) AppProxyIsc() IAppProxyIscAPI {
	return newAppProxyIscAPI(ap.codefresh)
}
