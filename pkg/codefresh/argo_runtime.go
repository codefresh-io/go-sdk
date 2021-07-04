package codefresh

import (
	"bytes"
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
		Data graphqlRuntimesInnerResponse`json:"data"`
	//	Errors []byte
	}
	graphqlRuntimesInnerResponse struct {
		Runtimes model.RuntimePage `json:"runtimes"`
	}
)

func newArgoRuntimeAPI(codefresh *codefresh) IArgoRuntimeAPI {
	return &argoRuntime{codefresh: codefresh}
}
func (r *argoRuntime)  List() ([]model.Runtime, error) {

	jsonData := map[string]interface{}{
        "query":` 
		{
			runtimes{
			  edges{
				node {
				  id
				  namespace
				  objectMeta {
					name
					description
				  }
				}
			  }
			}
		  }
        `,
    }


	response, err := r.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path: "/argo/api/graphql",
		body:  jsonData,
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
	runtimes := []model.Runtime{}
	for _, v := range res.Data.Runtimes.Edges {
		runtimes = append(runtimes, *v.Node)
	}
	
	return runtimes, nil
	
	
}