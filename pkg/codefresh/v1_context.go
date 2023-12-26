package codefresh

type (
	V1ContextAPI interface {
		GetDefaultGitContext() (*ContextPayload, error)
		GetGitContextByName(name string) (*ContextPayload, error)
		GetGitContexts() ([]ContextPayload, error)
	}

	v1Context struct {
		codefresh *codefresh
	}

	ContextPayload struct {
		Metadata struct {
			Name string `json:"name"`
		}
		Spec struct {
			Type string `json:"type"`
			Data struct {
				Auth struct {
					Type     string `json:"type"`
					Username string `json:"username"`
					Password string `json:"password"`
					ApiHost  string `json:"apiHost"`
					// for gitlab
					ApiURL         string `json:"apiURL"`
					ApiPathPrefix  string `json:"apiPathPrefix"`
					SshPrivateKey  string `json:"sshPrivateKey"`
					AppId          string `json:"appId"`
					InstallationId string `json:"installationId"`
					PrivateKey     string `json:"privateKey"`
				} `json:"auth"`
			} `json:"data"`
		} `json:"spec"`
	}

	GitContextsQs struct {
		Type    []string `url:"type"`
		Decrypt string   `url:"decrypt"`
	}
)

func (c v1Context) GetDefaultGitContext() (*ContextPayload, error) {
	resp, err := c.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/contexts/git/default",
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	result := &ContextPayload{}
	err = c.codefresh.decodeResponseInto(resp, result)
	return result, err
}

func (c v1Context) GetGitContextByName(name string) (*ContextPayload, error) {
	var qs = map[string]string{
		"decrypt": "true",
	}

	resp, err := c.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/contexts/" + name,
		qs:     qs,
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	result := &ContextPayload{}
	err = c.codefresh.decodeResponseInto(resp, result)
	return result, err
}

func (c v1Context) GetGitContexts() ([]ContextPayload, error) {
	qs := GitContextsQs{
		Type:    []string{"git.github", "git.gitlab", "git.github-app"},
		Decrypt: "true",
	}

	resp, err := c.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/contexts",
		qs:     qs,
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	result := make([]ContextPayload, 0)
	err = c.codefresh.decodeResponseInto(resp, result)
	return result, err
}
