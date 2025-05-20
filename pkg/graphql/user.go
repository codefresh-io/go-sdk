package graphql

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
)

type (
	UserAPI interface {
		GetCurrent(ctx context.Context) (*platmodel.User, error)
	}

	user struct {
		client *client.CfClient
	}
)

func (c *user) GetCurrent(ctx context.Context) (*platmodel.User, error) {
	query := `
query Me {
	me {
		id
		name
		email
		isAdmin
		accounts {
			id
			name
		}
		activeAccount {
			id
			name
			gitProvider
			gitApiUrl
			sharedConfigRepo
			admins
			features
		}
	}
}`
	variables := map[string]any{}
	res, err := client.GraphqlAPI[platmodel.User](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed getting current user: %w", err)
	}

	return &res, nil
}
