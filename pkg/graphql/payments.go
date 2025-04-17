package graphql

import (
	"context"
	"fmt"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"

	"github.com/codefresh-io/go-sdk/pkg/client"
)

type (
	PaymentsAPI interface {
		GetLimitsStatus(ctx context.Context) (*platmodel.LimitsStatus, error)
	}

	payments struct {
		client *client.CfClient
	}
)

func (c *payments) GetLimitsStatus(ctx context.Context) (*platmodel.LimitsStatus, error) {
	query := `
query LimitsStatus {
  limitsStatus {
    usage {
      clusters
	  applications
    }
    limits {
      clusters
	  applications
    }
    status
  }
}`
	limitsStatus, err := client.GraphqlAPI[platmodel.LimitsStatus](ctx, c.client, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get limits status: %w", err)
	}

	return &limitsStatus, nil
}
