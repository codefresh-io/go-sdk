package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IRuntimeAPI interface {
		Create(ctx context.Context, runtimeName, cluster, runtimeVersion, ingressHost string, componentNames []string) (*model.RuntimeCreationResponse, error)
		Get(ctx context.Context, name string) (*model.Runtime, error)
		List(ctx context.Context) ([]model.Runtime, error)
		Delete(ctx context.Context, runtimeName string) (int, error)
	}

	argoRuntime struct {
		codefresh *codefresh
	}

	graphqlRuntimesResponse struct {
		Data struct {
			Runtimes model.RuntimePage
		}
		Errors []graphqlError
	}

	graphqlRuntimeResponse struct {
		Data struct {
			Runtime *model.Runtime
		}
		Errors []graphqlError
	}

	graphQlRuntimeCreationResponse struct {
		Data struct {
			Runtime model.RuntimeCreationResponse
		}
		Errors []graphqlError
	}

	graphQlDeleteRuntimeResponse struct {
		Data struct {
			DeleteRuntime int
		}
		Errors []graphqlError
	}
)

func newArgoRuntimeAPI(codefresh *codefresh) IRuntimeAPI {
	return &argoRuntime{codefresh: codefresh}
}

func (r *argoRuntime) Create(ctx context.Context, runtimeName, cluster, runtimeVersion, ingressHost string, componentNames []string) (*model.RuntimeCreationResponse, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation CreateRuntime(
				$runtimeName: String!, $cluster: String!, $runtimeVersion: String!, $ingressHost: String, $componentNames: [String]!
			) {
				runtime(runtimeName: $runtimeName, cluster: $cluster, runtimeVersion: $runtimeVersion, ingressHost: $ingressHost, componentNames: $componentNames) {
					name
					newAccessToken
				}
			}
		`,
		"variables": map[string]interface{}{
			"runtimeName":    runtimeName,
			"cluster":        cluster,
			"runtimeVersion": runtimeVersion,
			"ingressHost":    ingressHost,
			"componentNames": componentNames,
		},
	}

	res := &graphQlRuntimeCreationResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call while creating runtime: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return &res.Data.Runtime, nil
}

func (r *argoRuntime) Get(ctx context.Context, name string) (*model.Runtime, error) {
	jsonData := map[string]interface{}{
		"query": `
			query GetRuntime(
				$name: String!
			) {
				runtime(name: $name) {
					metadata {
						name
						namespace
					}
					self {
						syncStatus
						healthMessage
						healthStatus
					}
					syncStatus
					healthMessage
					healthStatus
					cluster
					ingressHost
					runtimeVersion
					installationStatus
				}
			}
		`,
		"variables": map[string]interface{}{
			"name": name,
		},
	}

	res := graphqlRuntimeResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, &res)
	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call to runtime: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	if res.Data.Runtime == nil {
		return nil, fmt.Errorf("runtime '%s' does not exist", name)
	}

	return res.Data.Runtime, nil
}

func (r *argoRuntime) List(ctx context.Context) ([]model.Runtime, error) {
	jsonData := map[string]interface{}{
		"query": `{
			runtimes {
				edges {
					node {
						metadata {
							name
							namespace
						}
						self {
							syncStatus
							healthMessage
							healthStatus
						}
						syncStatus
						healthMessage
						healthStatus
						cluster
						ingressHost
						runtimeVersion
						installationStatus
					}
				}
			}
		}`,
	}

	res := &graphqlRuntimesResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting runtime list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	runtimes := make([]model.Runtime, len(res.Data.Runtimes.Edges))
	for i := range res.Data.Runtimes.Edges {
		runtimes[i] = *res.Data.Runtimes.Edges[i].Node
	}

	return runtimes, nil
}

func (r *argoRuntime) Delete(ctx context.Context, runtimeName string) (int, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation DeleteRuntime(
				$name: String!
			) {
				deleteRuntime(name: $name)
			}
		`,
		"variables": map[string]interface{}{
			"name": runtimeName,
		},
	}

	res := graphQlDeleteRuntimeResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, &res)
	if err != nil {
		return 0, fmt.Errorf("failed making a graphql API call to deleteRuntime: %w", err)
	}

	if len(res.Errors) > 0 {
		return 0, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.DeleteRuntime, nil
}
