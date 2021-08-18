package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IArgoPipelineAPI interface {
		List(ctx context.Context, runtimeName string) ([]model.Pipeline, error)
	}

	argoPipeline struct {
		codefresh *codefresh
	}

	graphqlPipelinesResponse struct {
		Data struct {
			Pipelines model.PipelinePage
		}
		Errors []graphqlError
	}
)

func newArgoPipelineAPI(codefresh *codefresh) IArgoPipelineAPI {
	return &argoPipeline{codefresh: codefresh}
}

func (r *argoPipeline) List(ctx context.Context, runtimeName string) ([]model.Pipeline, error) {
	jsonData := map[string]interface{}{
		"query": `
			query Pipelines($runtime: String!) {
				pipelines(runtime: $runtime) {
					edges {
						node {
							metadata {
								name
							}
							version
							self {
								status {
									syncStatus
									healthStatus
								}
							}
						}
					}
				}
			}`,
		"variables": map[string]interface{}{
			"runtime": runtimeName,
		},
	}

	res := &graphqlPipelinesResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting pipelines list: %w", err)
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
