package graphql

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
)

type (
	ComponentAPI interface {
		List(ctx context.Context, runtimeName string) ([]platmodel.Component, error)
	}

	component struct {
		client *client.CfClient
	}
)

func (c *component) List(ctx context.Context, runtimeName string) ([]platmodel.Component, error) {
	query := `
query Components($runtime: String!) {
	components(runtime: $runtime) {
		edges {
			node {
				metadata {
					name
					runtime
				}
				version
				self {
					status {
						syncStatus
						healthStatus
					}
					errors {
						...on SyncError{
							title
							message
							suggestion
							level
						}
					}
				}
			}
		}
	}
}`
	variables := map[string]any{
		"runtime": runtimeName,
	}
	res, err := client.GraphqlAPI[platmodel.ComponentSlice](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting component list: %w", err)
	}

	components := make([]platmodel.Component, len(res.Edges))
	for i := range res.Edges {
		components[i] = *res.Edges[i].Node
	}

	return components, nil
}
