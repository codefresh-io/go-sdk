package v1

import "github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"

type (
	V1API interface {
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

	v1Impl struct {
		client *client.CfClient
	}
)

func NewV1Client(c *client.CfClient) V1API {
	return &v1Impl{client: c}
}

func (v1 *v1Impl) Argo() ArgoAPI {
	return &argo{client: v1.client}
}

func (v1 *v1Impl) Cluster() ClusterAPI {
	return &cluster{client: v1.client}
}

func (v1 *v1Impl) Context() ContextAPI {
	return &v1Context{client: v1.client}
}

func (v1 *v1Impl) Gitops() GitopsAPI {
	return &gitops{client: v1.client}
}

func (v1 *v1Impl) Pipeline() PipelineAPI {
	return &pipeline{client: v1.client}
}

func (v1 *v1Impl) Progress() ProgressAPI {
	return &progress{client: v1.client}
}

func (v1 *v1Impl) RuntimeEnvironment() RuntimeEnvironmentAPI {
	return &runtimeEnvironment{client: v1.client}
}

func (v1 *v1Impl) Token() TokenAPI {
	return &token{client: v1.client}
}

func (v1 *v1Impl) User() UserAPI {
	return &user{client: v1.client}
}

func (v1 *v1Impl) Workflow() WorkflowAPI {
	return &workflow{codefresh: v1.client}
}
