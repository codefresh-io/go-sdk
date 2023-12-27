package v2

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	CliReleaseAPI interface {
		GetLatest(ctx context.Context) (string, error)
	}

	v2CliRelease struct {
		client *client.CfClient
	}
)

func (c *v2CliRelease) GetLatest(ctx context.Context) (string, error) {
	query := `
query LatestCliRelease {
	latestCliRelease 
}`
	args := map[string]interface{}{}
	resp, err := client.GraphqlAPI[string](ctx, c.client, query, args)
	if err != nil {
		return "", fmt.Errorf("failed getting latest cli release: %w", err)
	}

	return resp, nil
}
