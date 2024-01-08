package graphql

import "github.com/codefresh-io/go-sdk/pkg/client"

type (
	GraphQLAPI interface {
		Account() AccountAPI
		CliRelease() CliReleaseAPI
		Cluster() ClusterAPI
		Component() ComponentAPI
		GitSource() GitSourceAPI
		Pipeline() PipelineAPI
		Runtime() RuntimeAPI
		User() UserAPI
		Workflow() WorkflowAPI
	}

	gqlImpl struct {
		client *client.CfClient
	}
)

func NewGraphQLClient(c *client.CfClient) GraphQLAPI {
	return &gqlImpl{client: c}
}

func (v2 *gqlImpl) Account() AccountAPI {
	return &account{client: v2.client}
}

func (v2 *gqlImpl) CliRelease() CliReleaseAPI {
	return &cliRelease{client: v2.client}
}

func (v2 *gqlImpl) Cluster() ClusterAPI {
	return &cluster{client: v2.client}
}

func (v2 *gqlImpl) Component() ComponentAPI {
	return &component{client: v2.client}
}

func (v2 *gqlImpl) GitSource() GitSourceAPI {
	return &gitSource{client: v2.client}
}

func (v2 *gqlImpl) Pipeline() PipelineAPI {
	return &pipeline{client: v2.client}
}

func (v2 *gqlImpl) Runtime() RuntimeAPI {
	return &runtime{client: v2.client}
}

func (v2 *gqlImpl) User() UserAPI {
	return &user{client: v2.client}
}

func (v2 *gqlImpl) Workflow() WorkflowAPI {
	return &workflow{client: v2.client}
}
