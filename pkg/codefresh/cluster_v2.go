package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IClusterV2API interface {
		List(ctx context.Context, runtime string) ([]model.Cluster, error)
	}

	clusterV2 struct {
		codefresh *codefresh
	}

	graphqlClusterListResponse struct {
		Data struct {
			Clusters model.ClusterSlice
		}
		Errors []graphqlError
	}
)

func newClusterV2API(codefresh *codefresh) IClusterV2API {
	return &clusterV2{codefresh: codefresh}
}

func (c *clusterV2)List(ctx context.Context, runtime string) ([]model.Cluster, error) {
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


func (c *clusterV2)getClusterSlice(ctx context.Context, runtime string, after string) (*model.ClusterSlice, error) {
	jsonData := map[string]interface{}{
		"query": `
			query clusters($runtime: String!, $pagination: SlicePaginationArgs) {
				clusters(runtime: $runtime, pagination: $pagination) {
					edges {
						node {
							metadata {
								name
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
		return nil, fmt.Errorf("failed getting runtime list: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return &res.Data.Clusters, nil
}