package codefresh

type (
	V1API interface {
		Argo() ArgoAPI
		Clusters() IClusterAPI
		Contexts() IContextAPI
		Gitops() GitopsAPI
		Pipelines() IPipelineAPI
		Progresses() IProgressAPI
		RuntimeEnvironments() IRuntimeEnvironmentAPI
		Tokens() ITokenAPI
		Users() UsersAPI
		Workflows() IWorkflowAPI
	}

	v1Impl struct {
		codefresh *codefresh
	}
)

func newV1Client(c *codefresh) V1API {
	return &v1Impl{codefresh: c}
}

func (v1 *v1Impl) Argo() ArgoAPI {
	return newArgoAPI(v1.codefresh)
}

func (v1 *v1Impl) Clusters() IClusterAPI {
	return newClusterAPI(v1.codefresh)
}

func (v1 *v1Impl) Contexts() IContextAPI {
	return newContextAPI(v1.codefresh)
}

func (v1 *v1Impl) Gitops() GitopsAPI {
	return newGitopsAPI(v1.codefresh)
}

func (v1 *v1Impl) Pipelines() IPipelineAPI {
	return newPipelineAPI(v1.codefresh)
}

func (v1 *v1Impl) Progresses() IProgressAPI {
	return newProgressAPI(v1.codefresh)
}

func (v1 *v1Impl) RuntimeEnvironments() IRuntimeEnvironmentAPI {
	return newRuntimeEnvironmentAPI(v1.codefresh)
}

func (v1 *v1Impl) Tokens() ITokenAPI {
	return newTokenAPI(v1.codefresh)
}

func (v1 *v1Impl) Users() UsersAPI {
	return newUsersAPI(v1.codefresh)
}

func (v1 *v1Impl) Workflows() IWorkflowAPI {
	return newWorkflowAPI(v1.codefresh)
}
