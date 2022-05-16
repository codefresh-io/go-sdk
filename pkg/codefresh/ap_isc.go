package codefresh

import (
	"context"
	"fmt"
)

type (
	IAppProxyIscAPI interface {
		Create(ctx context.Context, runtime string, namespace string, repo string, clusterName string, clusterServer string) error
	}

	appProxyIsc struct {
		codefresh *codefresh
	}

	graphqlIscResponse struct {
		Errors []graphqlError
	}
)

func newAppProxyIscAPI(c *codefresh) IAppProxyIscAPI {
	return &appProxyIsc{codefresh: c}
}

func (c *appProxyIsc) Create(ctx context.Context, runtime string, namespace string, repo string, clusterName string, clusterServer string) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation createIsc($runtime: String!, $namespace: String!, $repo: String!, $clusterName: String!, $clusterServer: String!) {
				createIsc(runtime: $runtime, namespace: $namespace, repo: $repo, clusterName: $clusterName, clusterServer: $clusterServer)
			}
		`,
		"variables": map[string]interface{}{
			"runtime":       runtime,
			"namespace":     namespace,
			"repo":          repo,
			"clusterName":   clusterName,
			"clusterServer": clusterServer,
		},
	}

	res := &graphqlIscResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to create isc: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}
