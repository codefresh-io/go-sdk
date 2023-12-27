package v2

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model/platform"
)

type (
	WorkflowAPI interface {
		Get(ctx context.Context, uid string) (*platmodel.Workflow, error)
		List(ctx context.Context, filterArgs platmodel.WorkflowsFilterArgs) ([]platmodel.Workflow, error)
	}

	v2Workflow struct {
		client *client.CfClient
	}
)

func (c *v2Workflow) Get(ctx context.Context, uid string) (*platmodel.Workflow, error) {
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
	args := map[string]interface{}{
		"uid": uid,
	}
	resp, err := client.GraphqlAPI[platmodel.Workflow](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting a workflow: %w", err)
	}

	if resp.Metadata == nil {
		return nil, err
	}

	return &resp, nil
}

func (c *v2Workflow) List(ctx context.Context, filterArgs platmodel.WorkflowsFilterArgs) ([]platmodel.Workflow, error) {
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
	args := map[string]interface{}{
		"filters": filterArgs,
	}
	resp, err := client.GraphqlAPI[platmodel.WorkflowSlice](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting workflow list: %w", err)
	}

	workflows := make([]platmodel.Workflow, len(resp.Edges))
	for i := range resp.Edges {
		workflows[i] = *resp.Edges[i].Node
	}

	return workflows, nil
}
