package graphql

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
)

type (
	WorkflowAPI interface {
		Get(ctx context.Context, uid string) (*platmodel.Workflow, error)
		List(ctx context.Context, filterArgs platmodel.WorkflowsFilterArgs) ([]platmodel.Workflow, error)
	}

	workflow struct {
		client *client.CfClient
	}
)

func (c *workflow) Get(ctx context.Context, uid string) (*platmodel.Workflow, error) {
	query := `
query Workflow($uid: String!) {
	workflow(uid: $uid) {
		metadata {
			uid
			name
			namespace
			runtime
		}
		projects
		spec {
			entrypoint
			templates {
				name
			}
			workflowTemplateRef {
				name
				namespace
			}
			}
		status {
			phase
			progress
			nodes {
				type
				name
			}
			}
		pipeline {
			metadata {
				name
				namespace
			}
			}
	}
}`
	variables := map[string]any{
		"uid": uid,
	}
	res, err := client.GraphqlAPI[platmodel.Workflow](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting a workflow: %w", err)
	}

	if res.Metadata == nil {
		return nil, err
	}

	return &res, nil
}

func (c *workflow) List(ctx context.Context, filterArgs platmodel.WorkflowsFilterArgs) ([]platmodel.Workflow, error) {
	query := `
query Workflows($filters: WorkflowsFilterArgs) {
	workflows(filters: $filters) {
		edges {
			node {
				metadata {
					uid
					name
					namespace
					runtime
				}
				projects
				spec {
					entrypoint
					templates {
						name
					}
					workflowTemplateRef {
						name
						namespace
					}
					}
				status {
					phase
					progress
					nodes {
						type
						name
					}
					}
				pipeline {
					metadata {
						name
						namespace
					}
					}
			}
		}
	}
}`
	variables := map[string]any{
		"filters": filterArgs,
	}
	res, err := client.GraphqlAPI[platmodel.WorkflowSlice](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting workflow list: %w", err)
	}

	workflows := make([]platmodel.Workflow, len(res.Edges))
	for i := range res.Edges {
		workflows[i] = *res.Edges[i].Node
	}

	return workflows, nil
}
