package codefresh

import "github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"

type (
	V1API interface {
		Argo() V1ArgoAPI
		Cluster() V1ClusterAPI
		Context() V1ContextAPI
		Gitops() V1GitopsAPI
		Pipeline() V1PipelineAPI
		Progress() V1ProgressAPI
		RuntimeEnvironment() V1RuntimeEnvironmentAPI
		Token() V1TokenAPI
		User() V1UserAPI
		Workflow() V1WorkflowAPI
	}

	v1Impl struct {
		client *client.CfClient
	}
)

func newV1Client(c *client.CfClient) V1API {
	return &v1Impl{client: c}
}

func (v1 *v1Impl) Argo() V1ArgoAPI {
	return &v1Argo{client: v1.client}
}

func (v1 *v1Impl) Cluster() V1ClusterAPI {
	return &v1Cluster{client: v1.client}
}

func (v1 *v1Impl) Context() V1ContextAPI {
	return &v1Context{client: v1.client}
}

func (v1 *v1Impl) Gitops() V1GitopsAPI {
	return &v1Gitops{client: v1.client}
}

func (v1 *v1Impl) Pipeline() V1PipelineAPI {
	return &v1Pipeline{client: v1.client}
}

func (v1 *v1Impl) Progress() V1ProgressAPI {
	return &v1Progress{client: v1.client}
}

func (v1 *v1Impl) RuntimeEnvironment() V1RuntimeEnvironmentAPI {
	return &runtimeEnvironment{client: v1.client}
}

func (v1 *v1Impl) Token() V1TokenAPI {
	return &token{client: v1.client}
}

func (v1 *v1Impl) User() V1UserAPI {
	return &users{client: v1.client}
}

func (v1 *v1Impl) Workflow() V1WorkflowAPI {
	return &v1Workflow{codefresh: v1.client}
}
