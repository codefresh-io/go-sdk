package v2

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	V2AccountAPI interface {
		UpdateCsdpSettings(ctx context.Context, gitProvider platmodel.GitProviders, gitApiUrl, sharedConfigRepo string) error
	}

	v2Account struct {
		client *client.CfClient
	}
)

func (c *v2Account) UpdateCsdpSettings(ctx context.Context, gitProvider platmodel.GitProviders, gitApiUrl, sharedConfigRepo string) error {
	query := `
mutation updateCsdpSettings($gitProvider: GitProviders!, $gitApiUrl: String!, $sharedConfigRepo: String!) {
	updateCsdpSettings(gitProvider: $gitProvider, gitApiUrl: $gitApiUrl, sharedConfigRepo: $sharedConfigRepo)
}`
	args := map[string]interface{}{
		"gitProvider":      gitProvider,
		"gitApiUrl":        gitApiUrl,
		"sharedConfigRepo": sharedConfigRepo,
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, args)
	if err != nil {
		return fmt.Errorf("failed updating csdp settings: %w", err)
	}

	return nil
}
