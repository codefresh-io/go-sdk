package graphql

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
)

type (
	GitSourceAPI interface {
		List(ctc context.Context, runtimeName string) ([]platmodel.GitSource, error)
	}

	gitSource struct {
		client *client.CfClient
	}
)

func (c *gitSource) List(ctx context.Context, runtimeName string) ([]platmodel.GitSource, error) {
	query := `
query GitSources($runtime: String) {
	gitSources(runtime: $runtime) {
		edges {
			node {
				metadata {
					name
				}
				self {
					path
					repoURL
					status {
						syncStatus
						healthStatus
					}
				}
			}
		}
	}
}`
	variables := map[string]any{
		"runtime": runtimeName,
	}
	res, err := client.GraphqlAPI[platmodel.GitSourceSlice](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting git-source list: %w", err)
	}

	gitSources := make([]platmodel.GitSource, len(res.Edges))
	for i := range res.Edges {
		gitSources[i] = *res.Edges[i].Node
	}

	return gitSources, nil
}
