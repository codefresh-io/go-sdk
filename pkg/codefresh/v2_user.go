package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	V2UserAPI interface {
		GetCurrent(ctx context.Context) (*model.User, error)
	}

	v2User struct {
		*codefresh
	}

	graphQlMeResponse struct {
		Data struct {
			Me model.User
		}
		Errors []graphqlError
	}
)

func (u *v2User) GetCurrent(ctx context.Context) (*model.User, error) {
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
					gitProvider
					gitApiUrl
					sharedConfigRepo
					admins
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
