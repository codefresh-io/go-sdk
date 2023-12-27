package ap

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	apmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model/app-proxy"
)

type (
	GitIntegrationAPI interface {
		Add(ctx context.Context, args *apmodel.AddGitIntegrationArgs) (*apmodel.GitIntegration, error)
		Deregister(ctx context.Context, name *string) (*apmodel.GitIntegration, error)
		Edit(ctx context.Context, args *apmodel.EditGitIntegrationArgs) (*apmodel.GitIntegration, error)
		Get(ctx context.Context, name *string) (*apmodel.GitIntegration, error)
		List(ctx context.Context) ([]apmodel.GitIntegration, error)
		Register(ctx context.Context, args *apmodel.RegisterToGitIntegrationArgs) (*apmodel.GitIntegration, error)
		Remove(ctx context.Context, name string) error
	}

	gitIntegration struct {
		client *client.CfClient
	}
)

func (c *gitIntegration) Add(ctx context.Context, args *apmodel.AddGitIntegrationArgs) (*apmodel.GitIntegration, error) {
	query := `
mutation AddGitIntegration($args: AddGitIntegrationArgs!) {
	addGitIntegration(args: $args) {
		name
		sharingPolicy
		provider
		apiUrl
		registeredUsers
	}
}`
	res, err := client.GraphqlAPI[apmodel.GitIntegration](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed adding a git integration: %w", err)
	}

	return &res, nil
}

func (c *gitIntegration) Deregister(ctx context.Context, name *string) (*apmodel.GitIntegration, error) {
	query := `
mutation DeregisterFromGitIntegration($name: String) {
	deregisterFromGitIntegration(name: $name) {
		name
		sharingPolicy
		provider
		apiUrl
		registeredUsers
	}
}`
	args := map[string]any{
		"name": name,
	}
	res, err := client.GraphqlAPI[apmodel.GitIntegration](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed deregistering a git integration: %w", err)
	}

	return &res, nil
}

func (c *gitIntegration) Edit(ctx context.Context, args *apmodel.EditGitIntegrationArgs) (*apmodel.GitIntegration, error) {
	query := `
mutation EditGitIntegration($args: EditGitIntegrationArgs!) {
	editGitIntegration(args: $args) {
		name
		sharingPolicy
		provider
		apiUrl
		registeredUsers
	}
}`
	res, err := client.GraphqlAPI[apmodel.GitIntegration](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed editing a git integration: %w", err)
	}

	return &res, nil
}

func (c *gitIntegration) Get(ctx context.Context, name *string) (*apmodel.GitIntegration, error) {
	query := `
query GitIntegration($name: String) {
	gitIntegration(name: $name) {
		name
		sharingPolicy
		provider
		apiUrl
		registeredUsers
	}
}`
	args := map[string]any{
		"name": name,
	}
	res, err := client.GraphqlAPI[apmodel.GitIntegration](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting a git integration: %w", err)
	}

	return &res, nil
}

func (c *gitIntegration) List(ctx context.Context) ([]apmodel.GitIntegration, error) {
	query := `
query GitIntegrations {
	gitIntegrations {
		name
		sharingPolicy
		provider
		apiUrl
		users {
			userId
		}
	}
}`
	args := map[string]any{}
	res, err := client.GraphqlAPI[[]apmodel.GitIntegration](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting git integration list: %w", err)
	}

	return res, nil
}

func (c *gitIntegration) Register(ctx context.Context, args *apmodel.RegisterToGitIntegrationArgs) (*apmodel.GitIntegration, error) {
	query := `
mutation RegisterToGitIntegration($args: RegisterToGitIntegrationArgs!) {
	registerToGitIntegration(args: $args) {
		name
		sharingPolicy
		provider
		apiUrl
		registeredUsers
	}
}`
	res, err := client.GraphqlAPI[apmodel.GitIntegration](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed registering a git integration: %w", err)
	}

	return &res, nil
}

func (c *gitIntegration) Remove(ctx context.Context, name string) error {
	query := `
mutation RemoveGitIntegration($name: String!) {
	removeGitIntegration(name: $name)
}`
	args := map[string]any{
		"name": name,
	}
	_, err := client.GraphqlAPI[client.GraphqlVoidResponse](ctx, c.client, query, args)
	if err != nil {
		return fmt.Errorf("failed removing a git integration: %w", err)
	}

	return nil
}
