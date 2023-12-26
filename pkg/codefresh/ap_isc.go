package codefresh

import (
	"context"
	"fmt"
)

type (
	APIscAPI interface {
		RemoveRuntimeFromIscRepo(ctx context.Context) (int, error)
	}

	apIsc struct {
		codefresh *codefresh
	}

	graphqlAppProxyRemoveRuntimeFromIscRepoResponse struct {
		Data struct {
			RemoveRuntimeFromIscRepo int
		}
		Errors []graphqlError
	}
)

func (c *apIsc) RemoveRuntimeFromIscRepo(ctx context.Context) (int, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation RemoveRuntimeFromIscRepo {
				removeRuntimeFromIscRepo
			}
		`,
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
