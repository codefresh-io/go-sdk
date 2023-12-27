package v2

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	UserAPI interface {
		GetCurrent(ctx context.Context) (*platmodel.User, error)
	}

	v2User struct {
		client *client.CfClient
	}
)

func (c *v2User) GetCurrent(ctx context.Context) (*platmodel.User, error) {
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
		}
	}
}`
	args := map[string]interface{}{}
	resp, err := client.GraphqlAPI[platmodel.User](ctx, c.client, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed getting current user: %w", err)
	}

	return &resp, nil
}
