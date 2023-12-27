package ap

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	IscAPI interface {
		RemoveRuntimeFromIscRepo(ctx context.Context) (int, error)
	}

	isc struct {
		client *client.CfClient
	}
)

func (c *isc) RemoveRuntimeFromIscRepo(ctx context.Context) (int, error) {
	query := `
mutation RemoveRuntimeFromIscRepo {
	removeRuntimeFromIscRepo
}`
	variables := map[string]any{}
	res, err := client.GraphqlAPI[int](ctx, c.client, query, variables)
	if err != nil {
		return 0, fmt.Errorf("failed removing runtime from ISC repo: %w", err)
	}

	return res, nil
}
