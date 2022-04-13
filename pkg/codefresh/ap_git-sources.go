package codefresh

import (
	"context"
	"fmt"

	platformModel "github.com/codefresh-io/go-sdk/pkg/codefresh/model"
	appProxyModel "github.com/codefresh-io/go-sdk/pkg/codefresh/model/app-proxy"
)

type (
	IAppProxyGitSourcesAPI interface {
		Create(ctx context.Context, appName, appSpecifier, destServer, destNamespace string, isInternal bool) error
		Delete(ctx context.Context, appName string) error
		Edit(ctx context.Context, appName, appSpecifier string) error
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

	graphqlGitSourceCreateResponse struct {
		Errors []graphqlError
	}

	graphqlGitSourceDeleteResponse struct {
		Errors []graphqlError
	}

	graphqlGitSourceEditResponse struct {
		Errors []graphqlError
	}

)

func newAppProxyGitSourcesAPI(c *codefresh) IAppProxyGitSourcesAPI {
	return &appProxyGitSources{codefresh: c}
}

func (c *appProxyGitSources) Create(ctx context.Context, appName, appSpecifier, destServer, destNamespace string, isInternal bool) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation CreateGitSource($args: CreateGitSourceInput!) { 
				createGitSource(args: $args)
			}
		`, 
		"variables": map[string]interface{}{
			"args": appProxyModel.CreateGitSourceInput{
				AppName: appName,
				AppSpecifier: appSpecifier,
				DestServer: destServer,
				DestNamespace: destNamespace,
				IsInternal: &isInternal,
			},
		},
	}

	res := &graphqlGitSourceCreateResponse{}
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

	res := &graphqlGitSourceDeleteResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to delete git source: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}

func (c *appProxyGitSources) Edit(ctx context.Context, appName, appSpecifier string) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation EditGitSource($args: EditGitSourceInput!) { 
				editGitSource(args: $args)
			}
		`, 
		"variables": map[string]interface{}{
			"args": appProxyModel.EditGitSourceInput{
				AppName: appName,
				AppSpecifier: appSpecifier,
			},
		},
	}

	res := &graphqlGitSourceEditResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to edit git source: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}
