package codefresh

import (
	"context"
	"fmt"

	// model "github.com/codefresh-io/go-sdk/pkg/codefresh/model/app-proxy"
)

type (
	IAppProxyClustersAPI interface {
		RemoveCluster(ctx context.Context, server string, runtime string) error
	}

	graphqlClusterRemoveResponse struct {
		Errors []graphqlError
	}
)

func (c *gitIntegrations) RemoveCluster(ctx context.Context, server string, runtime string) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation RemoveCluster($server: String!, $runtime: String!) {
				removeCluster(server: $server, runtime: $runtime)
			}
		`,
		"variables": map[string]interface{}{
			"server": server,
			"runtime": runtime,
		},
	}

	res := &graphqlClusterRemoveResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to remove cluster: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}