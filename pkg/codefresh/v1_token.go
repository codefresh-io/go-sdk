package codefresh

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	V1TokenAPI interface {
		Create(name string, subject string) (*v1Token, error)
		List() ([]v1Token, error)
	}

	token struct {
		client *client.CfClient
	}

	v1Token struct {
		ID          string    `json:"_id"`
		Name        string    `json:"name"`
		TokenPrefix string    `json:"tokenPrefix"`
		Created     time.Time `json:"created"`
		Subject     struct {
			Type string `json:"type"`
			Ref  string `json:"ref"`
		} `json:"subject"`
		Value string
	}

	tokenSubjectType int

	getTokensReponse struct {
		Tokens []v1Token
	}
)

const (
	runtimeEnvironmentSubject tokenSubjectType = 0
)

func (s tokenSubjectType) String() string {
	return [...]string{"runtime-environment"}[s]
}

func (t *token) Create(name string, subject string) (*v1Token, error) {
	resp, err := t.client.RestAPI(nil, &client.RequestOptions{
		Path:   "/api/auth/key",
		Method: "POST",
		Body: map[string]interface{}{
			"name": name,
		},
		Query: map[string]string{
			"subjectReference": subject,
			"subjectType":      runtimeEnvironmentSubject.String(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating token: %w", err)
	}

	return &v1Token{
		Name:  name,
		Value: string(resp),
	}, err
}

func (t *token) List() ([]v1Token, error) {
	resp, err := t.client.RestAPI(nil, &client.RequestOptions{
		Path:   "/api/auth/keys",
		Method: "GET",
	})
	if err != nil {
		return nil, fmt.Errorf("failed listing tokens: %w", err)
	}

	result := make([]v1Token, 0)
	return result, json.Unmarshal(resp, &result)
}
