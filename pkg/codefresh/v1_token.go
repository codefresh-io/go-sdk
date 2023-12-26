package codefresh

import "time"

type (
	V1TokenAPI interface {
		Create(name string, subject string) (*v1Token, error)
		List() ([]v1Token, error)
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

	token struct {
		codefresh *codefresh
	}
)

const (
	runtimeEnvironmentSubject tokenSubjectType = 0
)

func (s tokenSubjectType) String() string {
	return [...]string{"runtime-environment"}[s]
}

func (t *token) Create(name string, subject string) (*v1Token, error) {
	resp, err := t.codefresh.requestAPI(&requestOptions{
		path:   "/api/auth/key",
		method: "POST",
		body: map[string]interface{}{
			"name": name,
		},
		qs: map[string]string{
			"subjectReference": subject,
			"subjectType":      runtimeEnvironmentSubject.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	value, err := t.codefresh.getBodyAsString(resp)
	if err != nil {
		return nil, err
	}

	return &v1Token{
		Name:  name,
		Value: value,
	}, err
}

func (t *token) List() ([]v1Token, error) {
	emptySlice := make([]v1Token, 0)
	resp, err := t.codefresh.requestAPI(&requestOptions{
		path:   "/api/auth/keys",
		method: "GET",
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = t.codefresh.decodeResponseInto(resp, &emptySlice)
	if err != nil {
		return nil, err
	}

	return emptySlice, err
}
