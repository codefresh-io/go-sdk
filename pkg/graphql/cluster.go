package graphql

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
)

type (
	ClusterAPI interface {
		List(ctx context.Context, runtime string) ([]platmodel.Cluster, error)
	}

	cluster struct {
		client *client.CfClient
	}
)

func (c *cluster) List(ctx context.Context, runtime string) ([]platmodel.Cluster, error) {
	after := ""
	clusters := make([]platmodel.Cluster, 0)
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

func (c *cluster) getClusterSlice(ctx context.Context, runtime string, after string) (*platmodel.ClusterSlice, error) {
	query := `
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
}`
	variables := map[string]any{
		"runtime": runtime,
		"pagination": map[string]any{
			"after": after,
		},
	}
	res, err := client.GraphqlAPI[platmodel.ClusterSlice](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting cluster list: %w", err)
	}

	return &res, nil
}
