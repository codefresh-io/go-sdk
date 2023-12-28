package graphql

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
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
	variables := map[string]any{
		"runtime":   runtime,
		"name":      name,
		"namespace": namespace,
	}
	res, err := client.GraphqlAPI[platmodel.Pipeline](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting a pipeline: %w", err)
	}

	return &res, nil
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
	variables := map[string]any{
		"filters": filterArgs,
	}
	res, err := client.GraphqlAPI[platmodel.PipelineSlice](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting pipeline list: %w", err)
	}

	pipelines := make([]platmodel.Pipeline, len(res.Edges))
	for i := range res.Edges {
		pipelines[i] = *res.Edges[i].Node
	}

	return pipelines, nil
}
