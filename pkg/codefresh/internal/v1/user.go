package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	UserAPI interface {
		GetCurrent(ctx context.Context) (*v1User, error)
	}

	users struct {
		client *client.CfClient
	}

	v1User struct {
		ID                string    `json:"_id"`
		Name              string    `json:"userName"`
		Email             string    `json:"email"`
		Accounts          []Account `json:"account"`
		ActiveAccountName string    `json:"activeAccountName"`
		Roles             []string  `json:"roles"`
		UserData          struct {
			Avatar string `json:"image"`
		} `json:"user_data"`
	}

	Account struct {
		Name string `json:"name"`
		ID   string `json:"_id"`
	}
)

func (u *users) GetCurrent(ctx context.Context) (*v1User, error) {
	resp, err := u.client.RestAPI(ctx, &client.RequestOptions{
		Method: "GET",
		Path:   "/api/user",
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting current user: %w", err)
	}

	result := &v1User{}
	return result, json.Unmarshal(resp, result)
}

func (u *v1User) GetActiveAccount() *Account {
	for i := 0; i < len(u.Accounts); i++ {
		if u.Accounts[i].Name == u.ActiveAccountName {
			return &u.Accounts[i]
		}
	}

	return nil
}
