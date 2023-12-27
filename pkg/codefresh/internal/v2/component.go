package v2

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model/platform"
)

type (
	ComponentAPI interface {
		List(ctx context.Context, runtimeName string) ([]platmodel.Component, error)
	}

	v2Component struct {
		client *client.CfClient
	}
)

func (c *v2Component) List(ctx context.Context, runtimeName string) ([]platmodel.Component, error) {
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
	args := map[string]interface{}{
		"runtime": runtimeName,
	}
	resp, err := client.GraphqlAPI[platmodel.ComponentSlice](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting component list: %w", err)
	}

	components := make([]platmodel.Component, len(resp.Edges))
	for i := range resp.Edges {
		components[i] = *resp.Edges[i].Node
	}

	return components, nil
}
