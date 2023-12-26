package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	V2RuntimeAPI interface {
		Create(ctx context.Context, opts *model.RuntimeInstallationArgs) (*model.RuntimeCreationResponse, error)
		Delete(ctx context.Context, runtimeName string) (int, error)
		DeleteManaged(ctx context.Context, runtimeName string) (int, error)
		Get(ctx context.Context, name string) (*model.Runtime, error)
		List(ctx context.Context) ([]model.Runtime, error)
		MigrateRuntime(ctx context.Context, runtimeName string) error
		ReportErrors(ctx context.Context, opts *model.ReportRuntimeErrorsArgs) (int, error)
		SetSharedConfigRepo(ctx context.Context, suggestedSharedConfigRepo string) (string, error)
	}

	v2Runtime struct {
		codefresh *codefresh
	}

	graphqlRuntimesResponse struct {
		Data struct {
			Runtimes model.RuntimeSlice
		}
		Errors []graphqlError
	}

	graphqlRuntimeResponse struct {
		Data struct {
			Runtime *model.Runtime
		}
		Errors []graphqlError
	}

	graphQlRuntimeCreationResponse struct {
		Data struct {
			CreateRuntime model.RuntimeCreationResponse
		}
		Errors []graphqlError
	}

	graphQlReportRuntimeErrorsResponse struct {
		Data struct {
			ReportRuntimeErrors int
		}
		Errors []graphqlError
	}

	graphQlDeleteRuntimeResponse struct {
		Data struct {
			DeleteRuntime int
		}
		Errors []graphqlError
	}

	graphQlDeleteManagedRuntimeResponse struct {
		Data struct {
			DeleteManagedRuntime int
		}
		Errors []graphqlError
	}

	graphQlSuggestIscRepoResponse struct {
		Data struct {
			SuggestIscRepo string
		}
		Errors []graphqlError
	}
)

func (r *v2Runtime) Create(ctx context.Context, opts *model.RuntimeInstallationArgs) (*model.RuntimeCreationResponse, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation CreateRuntime($installationArgs: RuntimeInstallationArgs!) {
				createRuntime(installationArgs: $installationArgs) {
					name
					newAccessToken
				}
			}
		`,
		"variables": map[string]interface{}{
			"installationArgs": opts,
		},
	}

	res := &graphQlRuntimeCreationResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call while creating runtime: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return &res.Data.CreateRuntime, nil
}

func (r *v2Runtime) Delete(ctx context.Context, runtimeName string) (int, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation DeleteRuntime(
				$name: String!
			) {
				deleteRuntime(name: $name)
			}
		`,
		"variables": map[string]interface{}{
			"name": runtimeName,
		},
	}

	res := graphQlDeleteRuntimeResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, &res)
	if err != nil {
		return 0, fmt.Errorf("failed making a graphql API call to deleteRuntime: %w", err)
	}

	if len(res.Errors) > 0 {
		return 0, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.DeleteRuntime, nil
}

func (r *v2Runtime) DeleteManaged(ctx context.Context, runtimeName string) (int, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation DeleteManagedRuntime(
				$name: String!
			) {
				deleteManagedRuntime(name: $name)
			}
		`,
		"variables": map[string]interface{}{
			"name": runtimeName,
		},
	}

	res := graphQlDeleteManagedRuntimeResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, &res)
	if err != nil {
		return 0, fmt.Errorf("failed making a graphql API call to deleteManagedRuntime: %w", err)
	}

	if len(res.Errors) > 0 {
		return 0, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.DeleteManagedRuntime, nil
}

func (r *v2Runtime) Get(ctx context.Context, name string) (*model.Runtime, error) {
	jsonData := map[string]interface{}{
		"query": `
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
			}
		`,
		"variables": map[string]interface{}{
			"name": name,
		},
	}

	res := graphqlRuntimeResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, &res)
	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call to runtime: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	if res.Data.Runtime == nil {
		return nil, fmt.Errorf("runtime '%s' does not exist", name)
	}

	return res.Data.Runtime, nil
}

func (r *v2Runtime) List(ctx context.Context) ([]model.Runtime, error) {
	jsonData := map[string]interface{}{
		"query": `{
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
		}`,
	}

	res := &graphqlRuntimesResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting runtime list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	runtimes := make([]model.Runtime, len(res.Data.Runtimes.Edges))
	for i := range res.Data.Runtimes.Edges {
		runtimes[i] = *res.Data.Runtimes.Edges[i].Node
	}

	return runtimes, nil
}

func (r *v2Runtime) MigrateRuntime(ctx context.Context, runtimeName string) error {
	jsonData := map[string]interface{}{
		"query": `
		  mutation migrateRuntime($runtimeName: String!) {
			migrateRuntime(runtimeName: $runtimeName)
		  }
		`,
		"variables": map[string]interface{}{
			"runtimeName": runtimeName,
		},
	}
	res := &graphqlVoidResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return fmt.Errorf("failed making a graphql API call to migrate runtime: %w", err)
	}

	if len(res.Errors) > 0 {
		return graphqlErrorResponse{errors: res.Errors}
	}

	return nil
}

func (r *v2Runtime) ReportErrors(ctx context.Context, opts *model.ReportRuntimeErrorsArgs) (int, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation ReportRuntimeErrors(
				$reportErrorsArgs: ReportRuntimeErrorsArgs!
			) {
				reportRuntimeErrors(reportErrorsArgs: $reportErrorsArgs)
			}
		`,
		"variables": map[string]interface{}{
			"reportErrorsArgs": opts,
		},
	}

	res := graphQlReportRuntimeErrorsResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, &res)
	if err != nil {
		return 0, fmt.Errorf("failed making a graphql API call to runtimeErrorReport: %w", err)
	}

	if len(res.Errors) > 0 {
		return 0, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.ReportRuntimeErrors, nil
}

func (r *v2Runtime) SetSharedConfigRepo(ctx context.Context, suggestedSharedConfigRepo string) (string, error) {
	jsonData := map[string]interface{}{
		"query": `
			mutation suggestIscRepo($suggestedSharedConfigRepo: String!) {
				suggestIscRepo(suggestedSharedConfigRepo: $suggestedSharedConfigRepo)
			}
		`,
		"variables": map[string]interface{}{
			"suggestedSharedConfigRepo": suggestedSharedConfigRepo,
		},
	}

	res := &graphQlSuggestIscRepoResponse{}
	err := r.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return "", fmt.Errorf("failed making a graphql API call while setting shared config repo: %w", err)
	}

	if len(res.Errors) > 0 {
		return "nil", graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.SuggestIscRepo, nil
}