package codefresh

import (
	"context"
	"fmt"

	model "github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IGitSourceAPI interface {
		List(ctc context.Context, runtimeName string) ([]model.GitSource, error)
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
