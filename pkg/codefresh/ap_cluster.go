package codefresh

import (
	"context"
	"fmt"
)

type (
	APClusterAPI interface {
		CreateArgoRollouts(ctx context.Context, server string, namespace string) error
		Delete(ctx context.Context, server string, runtime string) error
	}

	apCluster struct {
		codefresh *codefresh
	}
)

func (c *apCluster) CreateArgoRollouts(ctx context.Context, server string, namespace string) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation createArgoRollouts($args: CreateArgoRolloutsInput!) {
				createArgoRollouts(args: $args)
			}
		`,
		"variables": map[string]interface{}{
			"args": map[string]interface{}{
				"destServer":    server,
				"destNamespace": namespace,
			},
		},
	}

	res := &graphqlVoidResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to add rollouts: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}

func (c *apCluster) Delete(ctx context.Context, server string, runtime string) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation RemoveCluster($server: String!, $runtime: String!) {
				removeCluster(server: $server, runtime: $runtime)
			}
		`,
		"variables": map[string]interface{}{
			"server":  server,
			"runtime": runtime,
		},
	}

	res := &graphqlVoidResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to remove cluster: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}
