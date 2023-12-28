package appproxy

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
	apmodel "github.com/codefresh-io/go-sdk/pkg/model/app-proxy"
)

type (
	GitSourceAPI interface {
		Create(ctx context.Context, opts *apmodel.CreateGitSourceInput) error
		Delete(ctx context.Context, appName string) error
		Edit(ctx context.Context, opts *apmodel.EditGitSourceInput) error
	}

	gitSource struct {
		client *client.CfClient
	}
)

func (c *gitSource) Create(ctx context.Context, args *apmodel.CreateGitSourceInput) error {
	query := `
mutation CreateGitSource($args: CreateGitSourceInput!) { 
	createGitSource(args: $args)
}`
	variables := map[string]any{
		"args": args,
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, variables)
	if err != nil {
		return fmt.Errorf("failed creating a git-source: %w", err)
	}

	return nil
}

func (c *gitSource) Delete(ctx context.Context, appName string) error {
	query := `
mutation DeleteApplication($args: DeleteApplicationInput!) { 
	deleteApplication(args: $args)
}`
	variables := map[string]any{
		"args": map[string]string{
			"appName": appName,
		},
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, variables)
	if err != nil {
		return fmt.Errorf("failed deleting a git-source: %w", err)
	}

	return nil
}

func (c *gitSource) Edit(ctx context.Context, args *apmodel.EditGitSourceInput) error {
	query := `
mutation EditGitSource($args: EditGitSourceInput!) { 
	editGitSource(args: $args)
}`
	variables := map[string]any{
		"args": args,
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, variables)
	if err != nil {
		return fmt.Errorf("failed editing a git-source: %w", err)
	}

	return nil
}
