package codefresh

import (
	"context"
	"fmt"

	model "github.com/codefresh-io/go-sdk/pkg/codefresh/model/app-proxy"
)

type (
	IAppProxyGitIntegrationsAPI interface {
		List(ctx context.Context) ([]model.GitIntegration, error)
		Get(ctx context.Context, name string) (*model.GitIntegration, error)
		Add(ctx context.Context, args *model.AddGitIntegrationArgs) (*model.GitIntegration, error)
		Edit(ctx context.Context, args *model.EditGitIntegrationArgs) (*model.GitIntegration, error)
		Remove(ctx context.Context, name string) error
		Register(ctx context.Context, args *model.RegisterToGitIntegrationArgs) (*model.GitIntegration, error)
		Deregister(ctx context.Context, name string) (*model.GitIntegration, error)
	}

	gitIntegrations struct {
		codefresh *codefresh
	}

	graphqlGitIntegrationsListResponse struct {
		Data struct {
			GitIntegrations []model.GitIntegration
		}
		Errors []graphqlError
	}

	graphqlGitIntegrationsGetResponse struct {
		Data struct {
			GitIntegration *model.GitIntegration
		}
		Errors []graphqlError
	}

	graphqlGitIntegrationsAddResponse struct {
		Data struct {
			AddGitIntegration *model.GitIntegration
		}
		Errors []graphqlError
	}

	graphqlGitIntegrationsEditResponse struct {
		Data struct {
			EditGitIntegration *model.GitIntegration
		}
		Errors []graphqlError
	}

	graphqlGitIntegrationsRemoveResponse struct {
		Errors []graphqlError
	}

	graphqlGitIntegrationsRegisterResponse struct {
		Data struct {
			RegisterToGitIntegration *model.GitIntegration
		}
		Errors []graphqlError
	}

	graphqlGitIntegrationsDeregisterResponse struct {
		Data struct {
			DeregisterFromGitIntegration *model.GitIntegration
		}
		Errors []graphqlError
	}
)

func newAppProxyGitIntegrationsAPI(c *codefresh) IAppProxyGitIntegrationsAPI {
	return &gitIntegrations{codefresh: c}
}

func (c *gitIntegrations) List(ctx context.Context) ([]model.GitIntegration, error) {
	jsonData := map[string]interface{}{
		"query": `
			{
				gitIntegrations {
					name
					sharingPolicy
					provider
					apiUrl
					registeredUsers
				}
			}`,
	}

	res := &graphqlGitIntegrationsListResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting git-integrations list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.GitIntegrations, nil
}

func (c *gitIntegrations) Get(ctx context.Context, name string) (*model.GitIntegration, error) {
	jsonData := map[string]interface{}{
		"query": `
			query GetGitIntegration($name: String!) {
				gitIntegration(name: $name) {
					name
					sharingPolicy
					provider
					apiUrl
					registeredUsers
				}
			}
		`,
		"variables": map[string]interface{}{
			"name": name,
		},
	}

	res := &graphqlGitIntegrationsGetResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call while getting git integration: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.GitIntegration, nil
}

func (c *gitIntegrations) Add(ctx context.Context, args *model.AddGitIntegrationArgs) (*model.GitIntegration, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation AddGitIntegration($args: AddGitIntegrationArgs!) {
				addGitIntegration(args: $args) {
					name
					sharingPolicy
					provider
					apiUrl
					registeredUsers
				}
			}
		`,
		"variables": map[string]interface{}{
			"args": args,
		},
	}

	res := &graphqlGitIntegrationsAddResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call while adding git integration: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.AddGitIntegration, nil
}

func (c *gitIntegrations) Edit(ctx context.Context, args *model.EditGitIntegrationArgs) (*model.GitIntegration, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation EditGitIntegration($args: EditGitIntegrationArgs!) {
				editGitIntegration(args: $args) {
					name
					sharingPolicy
					provider
					apiUrl
					registeredUsers
				}
			}
		`,
		"variables": map[string]interface{}{
			"args": args,
		},
	}

	res := &graphqlGitIntegrationsEditResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call to edit a git integration: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.EditGitIntegration, nil
}

func (c *gitIntegrations) Remove(ctx context.Context, name string) error {
	jsonData := map[string]interface{}{
		"query": `
			mutation RemoveGitIntegration($name: String!) {
				removeGitIntegration(name: $name)
			}
		`,
		"variables": map[string]interface{}{
			"name": name,
		},
	}

	res := &graphqlGitIntegrationsEditResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return fmt.Errorf("failed making a graphql API call to remove a git integration: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}

func (c *gitIntegrations) Register(ctx context.Context, args *model.RegisterToGitIntegrationArgs) (*model.GitIntegration, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation RegisterToGitIntegration($args: RegisterToGitIntegrationArgs!) {
				registerToGitIntegration(args: $args) {
					name
					sharingPolicy
					provider
					apiUrl
					registeredUsers
				}
			}
		`,
		"variables": map[string]interface{}{
			"args": args,
		},
	}

	res := &graphqlGitIntegrationsRegisterResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call to register to a git integration: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.RegisterToGitIntegration, nil
}

func (c *gitIntegrations) Deregister(ctx context.Context, name string) (*model.GitIntegration, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation DeregisterToGitIntegration($name: String!) {
				deregisterFromGitIntegration(name: $name) {
					name
					sharingPolicy
					provider
					apiUrl
					registeredUsers
				}
			}
		`,
		"variables": map[string]interface{}{
			"name": name,
		},
	}

	res := &graphqlGitIntegrationsDeregisterResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call to deregister from a git integration: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.DeregisterFromGitIntegration, nil
}
