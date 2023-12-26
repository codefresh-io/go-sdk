package codefresh

import (
	"context"
	"fmt"
)

type (
	V1UserAPI interface {
		GetCurrent(ctx context.Context) (*v1User, error)
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

	users struct {
		*codefresh
	}
)

func (u *users) GetCurrent(ctx context.Context) (*v1User, error) {
	result := &v1User{}
	resp, err := u.codefresh.requestAPIWithContext(ctx, &requestOptions{
		method: "GET",
		path:   "/api/user",
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(resp.Status)
	}

	if err := u.codefresh.decodeResponseInto(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (u *v1User) GetActiveAccount() *Account {
	for i := 0; i < len(u.Accounts); i++ {
		if u.Accounts[i].Name == u.ActiveAccountName {
			return &u.Accounts[i]
		}
	}
	return nil
}
