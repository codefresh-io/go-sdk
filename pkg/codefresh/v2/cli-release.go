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

	cliRelease struct {
		client *client.CfClient
	}
)

func (c *cliRelease) GetLatest(ctx context.Context) (string, error) {
	query := `
query LatestCliRelease {
	latestCliRelease 
}`
	variables := map[string]any{}
	resp, err := client.GraphqlAPI[string](ctx, c.client, query, variables)
	if err != nil {
		return "", fmt.Errorf("failed getting latest cli release: %w", err)
	}

	return resp, nil
}
