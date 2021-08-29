package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IPipelineV2API interface {
		Get(ctx context.Context, name, namespace, runtime string) (*model.Pipeline, error)
		List(ctx context.Context, filterArgs model.PipelinesFilterArgs) ([]model.Pipeline, error)
	}

	pipelineV2 struct {
		codefresh *codefresh
	}

	graphqlListPipelinesResponse struct {
		Data struct {
			Pipelines model.PipelinePage
		}
		Errors []graphqlError
	}

	graphqlGetPipelineResponse struct {
		Data struct {
			Pipeline model.Pipeline
		}
		Errors []graphqlError
	}
)

func newPipelineV2API(codefresh *codefresh) IPipelineV2API {
	return &pipelineV2{codefresh: codefresh}
}

func (p *pipelineV2) Get(ctx context.Context, name, namespace, runtime string) (*model.Pipeline, error) {
	jsonData := map[string]interface{}{
		"query": `
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
			}`,
		"variables": map[string]interface{}{
			"runtime":   runtime,
			"name":      name,
			"namespace": namespace,
		},
	}

	res := &graphqlGetPipelineResponse{}
	err := p.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting pipeline: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	if res.Data.Pipeline.Metadata == nil {
		return nil, err
	}

	return &res.Data.Pipeline, nil
}

func (p *pipelineV2) List(ctx context.Context, filterArgs model.PipelinesFilterArgs) ([]model.Pipeline, error) {
	jsonData := map[string]interface{}{
		"query": `
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
			}`,
		"variables": map[string]interface{}{
			"filters": filterArgs,
		},
	}

	res := &graphqlListPipelinesResponse{}
	err := p.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting pipeline list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	pipelines := make([]model.Pipeline, len(res.Data.Pipelines.Edges))
	for i := range res.Data.Pipelines.Edges {
		pipelines[i] = *res.Data.Pipelines.Edges[i].Node
	}

	return pipelines, nil
}
