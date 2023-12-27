package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	V2GitSourceAPI interface {
		List(ctc context.Context, runtimeName string) ([]platmodel.GitSource, error)
	}

	v2GitSource struct {
		client *client.CfClient
	}
)

func (c *v2GitSource) List(ctx context.Context, runtimeName string) ([]platmodel.GitSource, error) {
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
	args := map[string]interface{}{
		"runtime": runtimeName,
	}
	resp, err := client.GraphqlAPI[platmodel.GitSourceSlice](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting git-source list: %w", err)
	}

	gitSources := make([]platmodel.GitSource, len(resp.Edges))
	for i := range resp.Edges {
		gitSources[i] = *resp.Edges[i].Node
	}

	return gitSources, nil
}
