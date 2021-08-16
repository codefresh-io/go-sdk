package codefresh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IGitSourceAPI interface {
		List(runtimeName string) ([]model.GitSource, error)
	}

	gitSource struct {
		codefresh *codefresh
	}

	graphQlGitSourcesListResponse struct {
		Data struct {
			GitSources model.GitSourcePage
		}
		Errors []graphqlError
	}

	// graphQlRuntimeCreationResponse struct {
	// 	Data struct {
	// 		Runtime model.RuntimeCreationResponse
	// 	}
	// 	Errors []graphqlError
	// }
)

func newGitSourceAPI(codefresh *codefresh) IGitSourceAPI {
	return &gitSource{codefresh: codefresh}
}

func (g *gitSource) List(runtimeName string) ([]model.GitSource, error) {
	jsonData := map[string]interface{}{
		"query": `{
			gitSources {
			  edges {
				node {
				  metadata {
					name
				  }
				  self {
					path
					repoURL
					revision
					size
					status {
					  syncStatus
					  syncStartedAt
					  syncFinishedAt
					  healthStatus
					  revision
					}
				  }
				}
			  }
			}
		  }`,
	}

	response, err := g.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   "/argo/api/graphql",
		body:   jsonData,
	})
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("failed to read from response body")
		return nil, err
	}

	res := graphQlGitSourcesListResponse{}
	err = json.Unmarshal(data, &res)

	if err != nil {
		return nil, err
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
