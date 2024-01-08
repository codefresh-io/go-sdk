package rest

import "github.com/codefresh-io/go-sdk/pkg/client"

type (
	RestAPI interface {
		Argo() ArgoAPI
		Cluster() ClusterAPI
		Context() ContextAPI
		Gitops() GitopsAPI
		Pipeline() PipelineAPI
		Progress() ProgressAPI
		RuntimeEnvironment() RuntimeEnvironmentAPI
		Token() TokenAPI
		User() UserAPI
		Workflow() WorkflowAPI
	}

	restImpl struct {
		client *client.CfClient
	}
)

func NewRestClient(c *client.CfClient) RestAPI {
	return &restImpl{client: c}
}

func (v1 *restImpl) Argo() ArgoAPI {
	return &argo{client: v1.client}
}

func (v1 *restImpl) Cluster() ClusterAPI {
	return &cluster{client: v1.client}
}

func (v1 *restImpl) Context() ContextAPI {
	return &v1Context{client: v1.client}
}

func (v1 *restImpl) Gitops() GitopsAPI {
	return &gitops{client: v1.client}
}

func (v1 *restImpl) Pipeline() PipelineAPI {
	return &pipeline{client: v1.client}
}

func (v1 *restImpl) Progress() ProgressAPI {
	return &progress{client: v1.client}
}

func (v1 *restImpl) RuntimeEnvironment() RuntimeEnvironmentAPI {
	return &runtimeEnvironment{client: v1.client}
}

func (v1 *restImpl) Token() TokenAPI {
	return &token{client: v1.client}
}

func (v1 *restImpl) User() UserAPI {
	return &user{client: v1.client}
}

func (v1 *restImpl) Workflow() WorkflowAPI {
	return &workflow{codefresh: v1.client}
}
