package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	IUsersV2API interface {
		GetCurrent(ctx context.Context) (*model.User, error)
	}

	usersV2 struct {
		*codefresh
	}

	graphQlMeResponse struct {
		Data struct {
			Me model.User
		}
		Errors []graphqlError
	}
)

func newUsersV2API(codefresh *codefresh) IUsersV2API {
	return &usersV2{codefresh}
}

func (u *usersV2) GetCurrent(ctx context.Context) (*model.User, error) {
	jsonData := map[string]interface{}{
		"query": `{
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
					sharedConfigRepo
				}
			}
		}
		`,
	}

	res := &graphQlMeResponse{}
	err := u.codefresh.graphqlAPI(ctx, jsonData, res)

	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call while getting user info: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	return &res.Data.Me, nil
}
