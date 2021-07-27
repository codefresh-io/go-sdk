package codefresh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/codefresh-io/argo-platform/libs/ql/graph/model"
)

type (
	IArgoRuntimeAPI interface {
		List() ([]model.Runtime, error)
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
)

func newArgoRuntimeAPI(codefresh *codefresh) IArgoRuntimeAPI {
	return &argoRuntime{codefresh: codefresh}
}
func (r *argoRuntime) List() ([]model.Runtime, error) {

	jsonData := map[string]interface{}{
		"query": ` 
		{
			runtimes(
				pagination: {}
				project: "") {
			  edges {
				node {
				  metadata {
					name
					namespace
				  }
				  healthMessage
				  runtimeVersion
				  cluster
				}
			  }
			}
		  }
        `,
	}

	response, err := r.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   "/argo/api/graphql",
		body:   jsonData,
	})
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("failed to read from response body")
		return nil, err
	}
	res := graphqlRuntimesResponse{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	runtimes := make([]model.Runtime, len(res.Data.Runtimes.Edges))
	for i := range res.Data.Runtimes.Edges {
		runtimes[i] = *res.Data.Runtimes.Edges[i].Node
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return runtimes, nil

}
