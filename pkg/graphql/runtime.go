package graphql

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
)

type (
	RuntimeAPI interface {
		Create(ctx context.Context, opts *platmodel.RuntimeInstallationArgs) (*platmodel.RuntimeCreationResponse, error)
		Delete(ctx context.Context, runtimeName string) (int, error)
		DeleteManaged(ctx context.Context, runtimeName string) (int, error)
		Get(ctx context.Context, name string) (*platmodel.Runtime, error)
		List(ctx context.Context) ([]platmodel.Runtime, error)
		MigrateRuntime(ctx context.Context, runtimeName string) error
		ReportErrors(ctx context.Context, opts *platmodel.ReportRuntimeErrorsArgs) (int, error)
		SetSharedConfigRepo(ctx context.Context, suggestedSharedConfigRepo string) (string, error)
	}

	runtime struct {
		client *client.CfClient
	}
)

func (c *runtime) Create(ctx context.Context, opts *platmodel.RuntimeInstallationArgs) (*platmodel.RuntimeCreationResponse, error) {
	query := `
mutation CreateRuntime($installationArgs: RuntimeInstallationArgs!) {
	createRuntime(installationArgs: $installationArgs) {
		name
		newAccessToken
	}
}`
	variables := map[string]any{
		"installationArgs": opts,
	}
	resp, err := client.GraphqlAPI[platmodel.RuntimeCreationResponse](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed creating a runtime: %w", err)
	}

	return &resp, nil
}

func (c *runtime) Delete(ctx context.Context, runtimeName string) (int, error) {
	query := `
mutation DeleteRuntime($name: String!) {
	deleteRuntime(name: $name)
}`
	variables := map[string]any{
		"name": runtimeName,
	}
	resp, err := client.GraphqlAPI[int](ctx, c.client, query, variables)
	if err != nil {
		return 0, fmt.Errorf("failed deleting a runtime: %w", err)
	}

	return resp, nil
}

func (c *runtime) DeleteManaged(ctx context.Context, runtimeName string) (int, error) {
	query := `
mutation DeleteManagedRuntime(
	$name: String!
) {
	deleteManagedRuntime(name: $name)
}`
	variables := map[string]any{
		"name": runtimeName,
	}
	resp, err := client.GraphqlAPI[int](ctx, c.client, query, variables)
	if err != nil {
		return 0, fmt.Errorf("failed deleting a hosted runtime: %w", err)
	}

	return resp, nil
}

func (c *runtime) Get(ctx context.Context, name string) (*platmodel.Runtime, error) {
	query := `
query GetRuntime($name: String!) {
	runtime(name: $name) {
		metadata {
			name
			namespace
		}
		self {
			syncStatus
			healthMessage
			healthStatus
		}
		syncStatus
		healthStatus
		healthMessage
		cluster
		managed
		isRemoteClusterConnected
		ingressHost
		internalIngressHost
		ingressClass
		ingressController
		runtimeVersion
		installationStatus
		installationType
		repo
		managedClustersNum
		gitProvider
		accessMode
	}
}`
	variables := map[string]any{
		"name": name,
	}
	resp, err := client.GraphqlAPI[platmodel.Runtime](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting a runtime: %w", err)
	}

	if resp.Metadata.Name == "" {
		return nil, fmt.Errorf("runtime '%s' does not exist", name)
	}

	return &resp, nil
}

func (c *runtime) List(ctx context.Context) ([]platmodel.Runtime, error) {
	query := `
query Runtimes {
	runtimes {
		edges {
			node {
				metadata {
					name
					namespace
				}
				self {
					syncStatus
					healthMessage
					healthStatus
				}
				syncStatus
				healthMessage
				healthStatus
				managed
				cluster
				ingressHost
				runtimeVersion
				installationStatus
				installationType
			}
		}
	}
}`
	variables := map[string]any{}
	resp, err := client.GraphqlAPI[platmodel.RuntimeSlice](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting runtime list: %w", err)
	}

	runtimes := make([]platmodel.Runtime, len(resp.Edges))
	for i := range resp.Edges {
		runtimes[i] = *resp.Edges[i].Node
	}

	return runtimes, nil
}

func (c *runtime) MigrateRuntime(ctx context.Context, runtimeName string) error {
	query := `
mutation migrateRuntime($runtimeName: String!) {
	migrateRuntime(runtimeName: $runtimeName)
}`
	variables := map[string]any{
		"runtimeName": runtimeName,
	}
	_, err := client.GraphqlAPI[client.GraphqlBaseResponse](ctx, c.client, query, variables)
	if err != nil {
		return fmt.Errorf("failed migrating a runtime: %w", err)
	}

	return nil
}

func (c *runtime) ReportErrors(ctx context.Context, opts *platmodel.ReportRuntimeErrorsArgs) (int, error) {
	query := `
mutation ReportRuntimeErrors($reportErrorsArgs: ReportRuntimeErrorsArgs!) {
	reportRuntimeErrors(reportErrorsArgs: $reportErrorsArgs)
}`
	variables := map[string]any{
		"reportErrorsArgs": opts,
	}
	resp, err := client.GraphqlAPI[int](ctx, c.client, query, variables)
	if err != nil {
		return 0, fmt.Errorf("failed reporting errors: %w", err)
	}

	return resp, nil
}

func (c *runtime) SetSharedConfigRepo(ctx context.Context, suggestedSharedConfigRepo string) (string, error) {
	query := `
mutation SuggestIscRepo($suggestedSharedConfigRepo: String!) {
	suggestIscRepo(suggestedSharedConfigRepo: $suggestedSharedConfigRepo)
}`
	variables := map[string]any{
		"suggestedSharedConfigRepo": suggestedSharedConfigRepo,
	}
	resp, err := client.GraphqlAPI[string](ctx, c.client, query, variables)
	if err != nil {
		return "", fmt.Errorf("failed suggesting ISC repo: %w", err)
	}

	return resp, nil
}
