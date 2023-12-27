package ap

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	apmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model/app-proxy"
)

type (
	APGitSourceAPI interface {
		Create(ctx context.Context, opts *apmodel.CreateGitSourceInput) error
		Delete(ctx context.Context, appName string) error
		Edit(ctx context.Context, opts *apmodel.EditGitSourceInput) error
	}

	apGitSource struct {
		client *client.CfClient
	}
)

func (c *apGitSource) Create(ctx context.Context, opts *apmodel.CreateGitSourceInput) error {
	query := `
mutation CreateGitSource($args: CreateGitSourceInput!) { 
	createGitSource(args: $args)
}`
	args := map[string]any{
		"appName":       opts.AppName,
		"appSpecifier":  opts.AppSpecifier,
		"destServer":    opts.DestServer,
		"destNamespace": opts.DestNamespace,
		"isInternal":    opts.IsInternal,
		"include":       opts.Include,
		"exclude":       opts.Exclude,
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, args)
	if err != nil {
		return fmt.Errorf("failed creating a git-source: %w", err)
	}

	return nil
}

func (c *apGitSource) Delete(ctx context.Context, appName string) error {
	query := `
mutation DeleteApplication($args: DeleteApplicationInput!) { 
	deleteApplication(args: $args)
}`
	args := map[string]any{
		"appName": appName,
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, args)
	if err != nil {
		return fmt.Errorf("failed deleting a git-source: %w", err)
	}

	return nil
}

func (c *apGitSource) Edit(ctx context.Context, opts *apmodel.EditGitSourceInput) error {
	query := `
mutation EditGitSource($args: EditGitSourceInput!) { 
	editGitSource(args: $args)
}`
	args := map[string]any{
		"appName":      opts.AppName,
		"appSpecifier": opts.AppSpecifier,
		"include":      opts.Include,
		"exclude":      opts.Exclude,
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, args)
	if err != nil {
		return fmt.Errorf("failed editing a git-source: %w", err)
	}

	return nil
}
