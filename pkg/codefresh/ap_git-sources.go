package codefresh

import (
	"context"
	"fmt"

	platformModel "github.com/codefresh-io/go-sdk/pkg/codefresh/model"
	appProxyModel "github.com/codefresh-io/go-sdk/pkg/codefresh/model/app-proxy"
)

type (
	IAppProxyGitSourcesAPI interface {
		Create(ctx context.Context, opts *appProxyModel.CreateGitSourceInput) error
		Delete(ctx context.Context, appName string) error
		Edit(ctx context.Context, opts *appProxyModel.EditGitSourceInput) error
	}

	appProxyGitSources struct {
		codefresh *codefresh
	}

	graphqlGitSourceListResponse struct {
		Data struct {
			GitSources platformModel.GitSourceSlice
		}
		Errors []graphqlError
	}
)

func newAppProxyGitSourcesAPI(c *codefresh) IAppProxyGitSourcesAPI {
	return &appProxyGitSources{codefresh: c}
}

func (c *appProxyGitSources) Create(ctx context.Context, opts *appProxyModel.CreateGitSourceInput) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation CreateGitSource($args: CreateGitSourceInput!) { 
				createGitSource(args: $args)
			}
		`,
		"variables": map[string]interface{}{
			"args": appProxyModel.CreateGitSourceInput{
				AppName:       opts.AppName,
				AppSpecifier:  opts.AppSpecifier,
				DestServer:    opts.DestServer,
				DestNamespace: opts.DestNamespace,
				IsInternal:    opts.IsInternal,
				Include:       opts.Include,
				Exclude:       opts.Exclude,
			},
		},
	}

	res := &graphqlVoidResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to create git source: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}

func (c *appProxyGitSources) Delete(ctx context.Context, appName string) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation DeleteApplication($args: DeleteApplicationInput!) { 
				deleteApplication(args: $args)
			}
		`,
		"variables": map[string]interface{}{
			"args": appProxyModel.DeleteApplicationInput{
				AppName: appName,
			},
		},
	}

	res := &graphqlVoidResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to delete git source: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}

func (c *appProxyGitSources) Edit(ctx context.Context, opts *appProxyModel.EditGitSourceInput) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation EditGitSource($args: EditGitSourceInput!) { 
				editGitSource(args: $args)
			}
		`,
		"variables": map[string]interface{}{
			"args": appProxyModel.EditGitSourceInput{
				AppName:      opts.AppName,
				AppSpecifier: opts.AppSpecifier,
				Include:      opts.Include,
				Exclude:      opts.Exclude,
			},
		},
	}

	res := &graphqlVoidResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to edit git source: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}
