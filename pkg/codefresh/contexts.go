package codefresh

type (
	IContextAPI interface {
		GetGitContexts() (error, *[]ContextPayload)
		GetGitContextByName(name string) (error, *ContextPayload)
	}

	context struct {
		codefresh Codefresh
	}

	ContextPayload struct {
		Metadata struct {
			Name string `json:"name"`
		}
		Spec struct {
			Type string `json:"type"`
			Data struct {
				Auth struct {
					Password      string `json:"password"`
					ApiHost       string `json:"apiHost"`
					ApiPathPrefix string `json:"apiPathPrefix"`
				} `json:"auth"`
			} `json:"data"`
		} `json:"spec"`
	}
)

func (c context) GetGitContexts() (error, *[]ContextPayload) {
	var result []ContextPayload
	var qs = map[string]string{
		"type":    "git.github",
		"decrypt": "true",
	}

	resp, err := c.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/contexts",
		qs:     qs,
	})
	if err != nil {
		return err, nil
	}

	err = c.codefresh.decodeResponseInto(resp, &result)

	return nil, &result
}

func (c context) GetGitContextByName(name string) (error, *ContextPayload) {
	var result ContextPayload
	var qs = map[string]string{
		"decrypt": "true",
	}

	resp, err := c.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/contexts/" + name,
		qs:     qs,
	})
	if err != nil {
		return err, nil
	}

	err = c.codefresh.decodeResponseInto(resp, &result)

	return nil, &result
}
