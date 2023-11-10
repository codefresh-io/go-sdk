package codefresh

import (
	"context"
	"fmt"

	model "github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IGitSourceAPI interface {
		List(ctx context.Context, runtimeName string) ([]model.GitSource, error)
		Delete(ctx context.Context, runtimeName string, name string) error
	}

	gitSource struct {
		codefresh *codefresh
	}

	graphQlGitSourcesListResponse struct {
		Data struct {
			GitSources model.GitSourceSlice
		}
		Errors []graphqlError
	}

	graphQlDeleteGitSourceResponse struct {
		Errors []graphqlError
	}
)

func newGitSourceAPI(codefresh *codefresh) IGitSourceAPI {
	return &gitSource{codefresh: codefresh}
}

func (g *gitSource) List(ctx context.Context, runtimeName string) ([]model.GitSource, error) {
	jsonData := map[string]interface{}{
		"query": `
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
			}`,
		"variables": map[string]interface{}{
			"runtime": runtimeName,
		},
	}

	res := &graphQlGitSourcesListResponse{}
	err := g.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting git-source list: %w", err)
	}

	gitSources := make([]model.GitSource, len(res.Data.GitSources.Edges))
	for i := range res.Data.GitSources.Edges {
		gitSources[i] = *res.Data.GitSources.Edges[i].Node
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return gitSources, nil
}

func (g *gitSource) Delete(ctx context.Context, runtimeName string, name string) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation RemoveGitSource($runtime: String, $name: String) {
				removeGitSource(runtime: $runtime, name: $name) 
			}
		`,
		"variables": map[string]interface{}{
			"runtime": runtimeName,
			"name":    name,
		},
	}

	res := graphQlDeleteGitSourceResponse{}
	err := g.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return fmt.Errorf("failed making a graphql API call to removeGitSource: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}
