package v2

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model/platform"
)

type (
	PipelineAPI interface {
		Get(ctx context.Context, name, namespace, runtime string) (*platmodel.Pipeline, error)
		List(ctx context.Context, filterArgs platmodel.PipelinesFilterArgs) ([]platmodel.Pipeline, error)
	}

	pipeline struct {
		client *client.CfClient
	}
)

func (c *pipeline) Get(ctx context.Context, name, namespace, runtime string) (*platmodel.Pipeline, error) {
	query := `
query Pipeline(
	$runtime: String!
	$name: String!
	$namespace: String
) {
	pipeline(name: $name, namespace: $namespace, runtime: $runtime) {
		metadata {
			name
			namespace
			runtime
		}
		self {
			healthStatus
			syncStatus
			version
		}
		projects
		spec {
			trigger
		}
	}
}`
	args := map[string]interface{}{
		"runtime":   runtime,
		"name":      name,
		"namespace": namespace,
	}
	resp, err := client.GraphqlAPI[platmodel.Pipeline](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting a pipeline: %w", err)
	}

	return &resp, nil
}

func (c *pipeline) List(ctx context.Context, filterArgs platmodel.PipelinesFilterArgs) ([]platmodel.Pipeline, error) {
	query := `
query Pipelines($filters: PipelinesFilterArgs) {
	pipelines(filters: $filters) {
		edges {
			node {
				metadata {
					name
					namespace
					runtime
				}
				self {
					healthStatus
					syncStatus
					version
				}
				projects
				spec {
					trigger
				}
			}
		}
	}
}`
	args := map[string]interface{}{
		"filters": filterArgs,
	}
	resp, err := client.GraphqlAPI[platmodel.PipelineSlice](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting pipeline list: %w", err)
	}

	pipelines := make([]platmodel.Pipeline, len(resp.Edges))
	for i := range resp.Edges {
		pipelines[i] = *resp.Edges[i].Node
	}

	return pipelines, nil
}
