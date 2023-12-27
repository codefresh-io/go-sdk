package ap

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	model "github.com/codefresh-io/go-sdk/pkg/codefresh/model/app-proxy"
)

type (
	APVersionInfoAPI interface {
		VersionInfo(ctx context.Context) (*model.AppProxyVersionInfo, error)
	}

	apVersionInfo struct {
		client *client.CfClient
	}
)

func (c *apVersionInfo) VersionInfo(ctx context.Context) (*model.AppProxyVersionInfo, error) {
	query := `
query VersionInfo {
	versionInfo {
		version
		platformHost
		platformVersion
	}
}`
	args := map[string]any{}
	res, err := client.GraphqlAPI[model.AppProxyVersionInfo](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting version info: %w", err)
	}

	return &res, nil
}
