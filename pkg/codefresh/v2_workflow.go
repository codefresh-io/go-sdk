package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	V2WorkflowAPI interface {
		Get(ctx context.Context, uid string) (*model.Workflow, error)
		List(ctx context.Context, filterArgs model.WorkflowsFilterArgs) ([]model.Workflow, error)
	}

	v2Workflow struct {
		codefresh *codefresh
	}

	graphqlListWorkflowsResponse struct {
		Data struct {
			Workflows model.WorkflowSlice
		}
		Errors []graphqlError
	}

	graphqlGetWorkflowResponse struct {
		Data struct {
			Workflow model.Workflow
		}
		Errors []graphqlError
	}
)

func (w *v2Workflow) Get(ctx context.Context, uid string) (*model.Workflow, error) {
	jsonData := map[string]interface{}{
		"query": `
			query Workflow(
				$uid: String!
			) {
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
			}`,
		"variables": map[string]interface{}{
			"uid": uid,
		},
	}

	res := &graphqlGetWorkflowResponse{}
	err := w.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting workflow: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	if res.Data.Workflow.Metadata == nil {
		return nil, err
	}

	return &res.Data.Workflow, nil
}

func (w *v2Workflow) List(ctx context.Context, filterArgs model.WorkflowsFilterArgs) ([]model.Workflow, error) {
	jsonData := map[string]interface{}{
		"query": `
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
			}`,
		"variables": map[string]interface{}{
			"filters": filterArgs,
		},
	}

	res := &graphqlListWorkflowsResponse{}
	err := w.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting workflow list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	workflows := make([]model.Workflow, len(res.Data.Workflows.Edges))
	for i := range res.Data.Workflows.Edges {
		workflows[i] = *res.Data.Workflows.Edges[i].Node
	}

	return workflows, nil
}
