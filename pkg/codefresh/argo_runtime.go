package codefresh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IArgoRuntimeAPI interface {
		List() ([]model.Runtime, error)
		Create(runtimeName, cluster, runtimeVersion string) (*model.RuntimeCreationResponse, error)
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

	graphQlRuntimeCreationResponse struct {
		Data struct {
			Runtime model.RuntimeCreationResponse
		}
		Errors []graphqlError
	}
)

var qlEndPoint = "/2.0/api/graphql"

func newArgoRuntimeAPI(codefresh *codefresh) IArgoRuntimeAPI {
	return &argoRuntime{codefresh: codefresh}
}

func (r *argoRuntime) Create(runtimeName, cluster, runtimeVersion string) (*model.RuntimeCreationResponse, error) {
	jsonData := map[string]interface{}{
		"query": `mutation CreateRuntime($name: String!, $cluster: String!, $runtimeVersion: String!) {
		runtime(name: $name, cluster: $cluster, runtimeVersion: $runtimeVersion) {
		  name
		  newAccessToken
		}
	  }`,
		"variables": map[string]interface{}{
			"name":           runtimeName,
			"cluster":        cluster,
			"runtimeVersion": runtimeVersion,
		},
	}

	response, err := r.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   qlEndPoint,
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

	res := graphQlRuntimeCreationResponse{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return &res.Data.Runtime, nil
}

func (r *argoRuntime) List() ([]model.Runtime, error) {
	jsonData := map[string]interface{}{
		"query": ` 
		{
			runtimes
			(
				pagination: {}
				project: ""
			) {
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
		path:   qlEndPoint,
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
