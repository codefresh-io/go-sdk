package codefresh

import (
	"context"
	"fmt"
)

type (
	IAppProxyIscAPI interface {
		RemoveRuntimeFromIscRepo(ctx context.Context, runtimeName string) (int, error)
	}

	appProxyIsc struct {
		codefresh *codefresh
	}

	graphqlAppProxyRemoveRuntimeFromIscRepoResponse struct {
		Data struct {
			RemoveRuntimeFromIscRepo int
		}
		Errors []graphqlError
	}
)

func newAppProxyIscAPI(c *codefresh) IAppProxyIscAPI {
	return &appProxyIsc{codefresh: c}
}

func (c *appProxyIsc) RemoveRuntimeFromIscRepo(ctx context.Context, runtimeName string) (int, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation RemoveRuntimeFromIscRepo(
				$runtime: String!
			) {
				removeRuntimeFromIscRepo(runtime: $runtime)
			}
		`,
		"variables": map[string]interface{}{
			"name": runtimeName,
		},
	}

	res := &graphqlAppProxyRemoveRuntimeFromIscRepoResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return 0, fmt.Errorf("failed to remove runtime from isc repo: %w", err)
	}

	if len(res.Errors) > 0 {
		return 0, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.RemoveRuntimeFromIscRepo, nil
}
