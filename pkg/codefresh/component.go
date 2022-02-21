package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IComponentAPI interface {
		List(ctx context.Context, runtimeName string) ([]model.Component, error)
	}

	component struct {
		codefresh *codefresh
	}

	graphqlComponentsResponse struct {
		Data struct {
			Components model.ComponentSlice
		}
		Errors []graphqlError
	}
)

func newComponentAPI(codefresh *codefresh) IComponentAPI {
	return &component{codefresh: codefresh}
}

func (r *component) List(ctx context.Context, runtimeName string) ([]model.Component, error) {
	jsonData := map[string]interface{}{
		"query": `
			query Components($runtime: String!) {
				components(runtime: $runtime) {
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
								errors {
									title
									message
									suggestion
									level
									lastSeen
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

	res := &graphqlComponentsResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting components list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	components := make([]model.Component, len(res.Data.Components.Edges))
	for i := range res.Data.Components.Edges {
		components[i] = *res.Data.Components.Edges[i].Node
	}

	return components, nil
}
