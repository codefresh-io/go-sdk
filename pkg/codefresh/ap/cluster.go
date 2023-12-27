package ap

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	ClusterAPI interface {
		CreateArgoRollouts(ctx context.Context, server string, namespace string) error
		Delete(ctx context.Context, server string, runtime string) error
	}

	cluster struct {
		client *client.CfClient
	}
)

func (c *cluster) CreateArgoRollouts(ctx context.Context, server string, namespace string) error {
	query := `
mutation CreateArgoRollouts($args: CreateArgoRolloutsInput!) {
	createArgoRollouts(args: $args)
}`
	variables := map[string]any{
		"args": map[string]any{
			"destServer":    server,
			"destNamespace": namespace,
		},
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, variables)
	if err != nil {
		return fmt.Errorf("failed creating argo rollouts: %w", err)
	}

	return nil
}

func (c *cluster) Delete(ctx context.Context, server string, runtime string) error {
	query := `
mutation RemoveCluster($server: String!, $runtime: String!) {
	removeCluster(server: $server, runtime: $runtime)
}`
	variables := map[string]any{
		"server":  server,
		"runtime": runtime,
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, variables)
	if err != nil {
		return fmt.Errorf("failed deleting cluster: %w", err)
	}

	return nil
}
