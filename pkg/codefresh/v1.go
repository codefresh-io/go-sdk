package codefresh

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
		codefresh *codefresh
	}
)

func newV1Client(c *codefresh) V1API {
	return &v1Impl{codefresh: c}
}

func (v1 *v1Impl) Argo() V1ArgoAPI {
	return &v1Argo{codefresh: v1.codefresh}
}

func (v1 *v1Impl) Cluster() V1ClusterAPI {
	return &v1Cluster{codefresh: v1.codefresh}
}

func (v1 *v1Impl) Context() V1ContextAPI {
	return &v1Context{codefresh: v1.codefresh}
}

func (v1 *v1Impl) Gitops() V1GitopsAPI {
	return &v1Gitops{codefresh: v1.codefresh}
}

func (v1 *v1Impl) Pipeline() V1PipelineAPI {
	return &v1Pipeline{codefresh: v1.codefresh}
}

func (v1 *v1Impl) Progress() V1ProgressAPI {
	return &v1Progress{codefresh: v1.codefresh}
}

func (v1 *v1Impl) RuntimeEnvironment() V1RuntimeEnvironmentAPI {
	return &runtimeEnvironment{codefresh: v1.codefresh}
}

func (v1 *v1Impl) Token() V1TokenAPI {
	return &token{codefresh: v1.codefresh}
}

func (v1 *v1Impl) User() V1UserAPI {
	return &users{codefresh: v1.codefresh}
}

func (v1 *v1Impl) Workflow() V1WorkflowAPI {
	return &v1Workflow{codefresh: v1.codefresh}
}
