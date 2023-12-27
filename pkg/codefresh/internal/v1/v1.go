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
	return &v1Argo{client: v1.client}
}

func (v1 *v1Impl) Cluster() ClusterAPI {
	return &v1Cluster{client: v1.client}
}

func (v1 *v1Impl) Context() ContextAPI {
	return &v1Context{client: v1.client}
}

func (v1 *v1Impl) Gitops() GitopsAPI {
	return &v1Gitops{client: v1.client}
}

func (v1 *v1Impl) Pipeline() PipelineAPI {
	return &v1Pipeline{client: v1.client}
}

func (v1 *v1Impl) Progress() ProgressAPI {
	return &v1Progress{client: v1.client}
}

func (v1 *v1Impl) RuntimeEnvironment() RuntimeEnvironmentAPI {
	return &runtimeEnvironment{client: v1.client}
}

func (v1 *v1Impl) Token() TokenAPI {
	return &token{client: v1.client}
}

func (v1 *v1Impl) User() UserAPI {
	return &users{client: v1.client}
}

func (v1 *v1Impl) Workflow() WorkflowAPI {
	return &v1Workflow{codefresh: v1.client}
}
