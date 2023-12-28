package appproxy

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
	apmodel "github.com/codefresh-io/go-sdk/pkg/model/app-proxy"
)

type (
	VersionInfoAPI interface {
		VersionInfo(ctx context.Context) (*apmodel.AppProxyVersionInfo, error)
	}

	versionInfo struct {
		client *client.CfClient
	}
)

func (c *versionInfo) VersionInfo(ctx context.Context) (*apmodel.AppProxyVersionInfo, error) {
	query := `
query VersionInfo {
	versionInfo {
		version
		platformHost
		platformVersion
	}
}`
	variables := map[string]any{}
	res, err := client.GraphqlAPI[apmodel.AppProxyVersionInfo](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting version info: %w", err)
	}

	return &res, nil
}
