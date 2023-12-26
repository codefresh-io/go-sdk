package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	V2ClusterAPI interface {
		List(ctx context.Context, runtime string) ([]model.Cluster, error)
	}

	v2Cluster struct {
		codefresh *codefresh
	}

	graphqlClusterListResponse struct {
		Data struct {
			Clusters model.ClusterSlice
		}
		Errors []graphqlError
	}
)

func (c *v2Cluster) List(ctx context.Context, runtime string) ([]model.Cluster, error) {
	after := ""
	clusters := make([]model.Cluster, 0)
	for {
		clusterSlice, err := c.getClusterSlice(ctx, runtime, after)
		if err != nil {
			return nil, err
		}

		for i := range clusterSlice.Edges {
			clusters = append(clusters, *clusterSlice.Edges[i].Node)
		}

		if clusterSlice.PageInfo.HasNextPage {
			after = *clusterSlice.PageInfo.EndCursor
		} else {
			break
		}
	}

	return clusters, nil
}

func (c *v2Cluster) getClusterSlice(ctx context.Context, runtime string, after string) (*model.ClusterSlice, error) {
	jsonData := map[string]interface{}{
		"query": `
			query clusters($runtime: String, $pagination: SlicePaginationArgs) {
				clusters(runtime: $runtime, pagination: $pagination) {
					edges {
						node {
							metadata {
								name
								runtime
							}
							server
							info {
								connectionState {
									status
									message
								}
								serverVersion
								cacheInfo {
									resourcesCount
									apisCount
								}
							}
						}
					}
					pageInfo {
						endCursor
						hasNextPage
					}
				}
			}
		`,
		"variables": map[string]interface{}{
			"runtime": runtime,
			"pagination": map[string]interface{}{
				"after": after,
			},
		},
	}

	res := &graphqlClusterListResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed getting cluster list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return &res.Data.Clusters, nil
}
