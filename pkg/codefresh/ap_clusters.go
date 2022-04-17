package codefresh

import (
	"context"
	"fmt"
)

type (
	IAppProxyClustersAPI interface {
		AddRollouts(ctx context.Context, server string, namespace string) error
		RemoveCluster(ctx context.Context, server string, runtime string) error
	}
	appProxyClusters struct {
		codefresh *codefresh
	}

	graphqlClusterResponse struct {
		Errors []graphqlError
	}
)

func newAppProxyClustersAPI(c *codefresh) IAppProxyClustersAPI {
	return &appProxyClusters{codefresh: c}
}

func (c *appProxyClusters)AddRollouts(ctx context.Context, server string, namespace string) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation createRollouts($args: CreateRolloutsInput!) {
				createRollouts(args: $args)
			}
		`,
		"variables": map[string]interface{}{
			"destServer": server,
			"destNamespace": namespace,
		},
	}

	res := &graphqlClusterResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to add rollouts: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}

func (c *appProxyClusters)RemoveCluster(ctx context.Context, server string, runtime string) error {
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

	res := &graphqlClusterResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to remove cluster: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}