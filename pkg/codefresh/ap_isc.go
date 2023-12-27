package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	APIscAPI interface {
		RemoveRuntimeFromIscRepo(ctx context.Context) (int, error)
	}

	apIsc struct {
		client *client.CfClient
	}
)

func (c *apIsc) RemoveRuntimeFromIscRepo(ctx context.Context) (int, error) {
	query := `
mutation RemoveRuntimeFromIscRepo {
	removeRuntimeFromIscRepo
}`
	args := map[string]any{}
	res, err := client.GraphqlAPI[int](ctx, c.client, query, args)
	if err != nil {
		return 0, fmt.Errorf("failed removing runtime from ISC repo: %w", err)
	}

	return res, nil
}
