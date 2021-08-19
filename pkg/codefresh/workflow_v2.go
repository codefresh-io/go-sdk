package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IWorkflowV2API interface {
		Get(ctx context.Context, name, namespace, runtime string) (model.Workflow, error)
		List(ctx context.Context, args model.WorkflowsFilterArgs) ([]model.Workflow, error)
	}

	workflowV2 struct {
		codefresh *codefresh
	}

	graphqlListWorkflowsResponse struct {
		Data struct {
			Workflows model.WorkflowPage
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

func newWorkflowV2API(codefresh *codefresh) IWorkflowV2API {
	return &workflowV2{codefresh: codefresh}
}

func (w *workflowV2) Get(ctx context.Context, name, namespace, runtime string) (model.Workflow, error) {
	jsonData := map[string]interface{}{
		"query": `{
			workflow(
				runtime: String!
				name: String!
				namespace: String
			) {
				metadata {
					name
					namespace
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
					createdAt
					startedAt
					finishedAt
					phase
					progress {
					  total
					  done
					}
					nodes {
					  type
					  name
					}
					message
					statuses {
					  since
					  phase
					  message
					}
				  }
				pipeline {
					metadata {
					  name
					  namespace
					}
				  }
				actualManifest
			}
		}`,
		"variables": map[string]interface{}{
			"runtime":   runtime,
			"name":      name,
			"namespace": namespace,
		},
	}

	res := &graphqlGetWorkflowResponse{}
	err := w.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed getting pipeline list: %w", err)
	}

	if len(res.Errors) > 0 {
		return model.Workflow{}, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.Workflow, nil
}

func (w *workflowV2) List(ctx context.Context, args model.WorkflowsFilterArgs) ([]model.Workflow, error) {
	jsonData := map[string]interface{}{
		"query": `{
			workflows {
				edges {
					node {
						metadata {
							name
							namespace
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
							createdAt
							startedAt
							finishedAt
							phase
							progress {
							  total
							  done
							}
							nodes {
							  type
							  name
							}
							message
							statuses {
							  since
							  phase
							  message
							}
						  }
						pipeline {
							metadata {
							  name
							  namespace
							}
						  }
						actualManifest
					}
				}
			}
		}`,
		"variables": map[string]interface{}{
			"project": args.Project,
			"runtime": args.Runtime,
			"pipeline": args.Pipeline,
			"repositories": args.Repositories,
			"branches": args.Branches,
			"eventTypes": args.EventTypes,
			"initiators": args.Initiators,
			"statuses": args.Statuses,
			"startDate": args.StartDate,
		},
	}

	res := &graphqlListWorkflowsResponse{}
	err := w.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting pipeline list: %w", err)
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
