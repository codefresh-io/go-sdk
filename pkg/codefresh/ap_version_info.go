package codefresh

import (
	"context"
	"fmt"

	model "github.com/codefresh-io/go-sdk/pkg/codefresh/model/app-proxy"
)

type (
	IAppProxyVersionInfoAPI interface {
		VersionInfo(ctx context.Context) (*model.AppProxyVersionInfo, error)
	}

	appProxyVersionInfo struct {
		codefresh *codefresh
	}

	graphqlAppProxyVersionInfoResponse struct {
		Data struct {
			VersionInfo *model.AppProxyVersionInfo
		}
		Errors []graphqlError
	}
)

func newAppProxyVersionInfoAPI(c *codefresh) IAppProxyVersionInfoAPI {
	return &appProxyVersionInfo{codefresh: c}
}

func (c *appProxyVersionInfo) VersionInfo(ctx context.Context) (*model.AppProxyVersionInfo, error) {
	jsonData := map[string]interface{}{
		"query": `
			{
				versionInfo {
					version
					platformHost
					platformVersion
				}
			}`,
	}

	res := &graphqlAppProxyVersionInfoResponse{}
	err := c.codefresh.graphqlAPI(ctx, jsonData, res)
	if err != nil {
		return nil, fmt.Errorf("failed to get version info: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return res.Data.VersionInfo, nil
}
