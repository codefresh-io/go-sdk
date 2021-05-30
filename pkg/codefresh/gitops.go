package codefresh

import "fmt"

type (
	GitopsAPI interface {
		CreateEnvironment(name string, project string, application string, integration string) error
		SendEnvironment(environment Environment) (map[string]interface{}, error)
		DeleteEnvironment(name string) error
		GetEnvironments() ([]CFEnvironment, error)
		SendEvent(name string, props map[string]string) error
		SendApplicationResources(resources *ApplicationResources) error
	}

	gitops struct {
		codefresh *codefresh
	}
	CodefreshEvent struct {
		Event string            `json:"event"`
		Props map[string]string `json:"props"`
	}

	MongoCFEnvWrapper struct {
		Docs []CFEnvironment `json:"docs"`
	}

	CFEnvironment struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Spec struct {
			Type        string `json:"type"`
			Application string `json:"application"`
		} `json:"spec"`
	}

	EnvironmentMetadata struct {
		Name string `json:"name"`
	}

	EnvironmentSpec struct {
		Type        string `json:"type"`
		Context     string `json:"context"`
		Project     string `json:"project"`
		Application string `json:"application"`
	}

	EnvironmentPayload struct {
		Version  string              `json:"version"`
		Metadata EnvironmentMetadata `json:"metadata"`
		Spec     EnvironmentSpec     `json:"spec"`
	}

	SyncPolicy struct {
		AutoSync bool `json:"autoSync"`
	}

	Commit struct {
		Message *string `json:"message"`
		Avatar  *string `json:"avatar"`
	}

	EnvironmentActivityRS struct {
		From ReplicaState `json:"from"`
		To   ReplicaState `json:"to"`
	}

	ReplicaState struct {
		Current int64 `json:"current"`
		Desired int64 `json:"desired"`
	}

	Annotation struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	Gitops struct {
		Comitters []User       `json:"comitters"`
		Prs       []Annotation `json:"prs"`
		Issues    []Annotation `json:"issues"`
	}

	Environment struct {
		Gitops       Gitops                `json:"gitops"`
		FinishedAt   string                `json:"finishedAt"`
		HealthStatus string                `json:"healthStatus"`
		SyncStatus   string                `json:"status"`
		HistoryId    int64                 `json:"historyId"`
		SyncRevision string                `json:"revision"`
		Name         string                `json:"name"`
		Activities   []EnvironmentActivity `json:"activities"`
		Resources    interface{}           `json:"resources"`
		RepoUrl      string                `json:"repoUrl"`
		Commit       Commit                `json:"commit"`
		SyncPolicy   SyncPolicy            `json:"syncPolicy"`
		Date         string                `json:"date"`
		ParentApp    string                `json:"parentApp"`
	}

	EnvironmentActivity struct {
		Name         string                `json:"name"`
		TargetImages []string              `json:"targetImages"`
		Status       string                `json:"status"`
		LiveImages   []string              `json:"liveImages"`
		ReplicaSet   EnvironmentActivityRS `json:"replicaSet"`
	}

	ApplicationResources struct {
		Name      string      `json:"name, omitempty"`
		HistoryId int64       `json:"historyId"`
		Revision  string      `json:"revision, omitempty"`
		Resources interface{} `json:"resources"`
	}
)

func newGitopsAPI(codefresh *codefresh) GitopsAPI {
	return &gitops{codefresh}
}

func (a *gitops) CreateEnvironment(name string, project string, application string, integration string) error {
	_, err := a.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   "/api/environments-v2",
		body: &EnvironmentPayload{
			Version: "1.0",
			Metadata: EnvironmentMetadata{
				Name: name,
			},
			Spec: EnvironmentSpec{
				Type:        "argo",
				Context:     integration,
				Project:     project,
				Application: application,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *gitops) SendEnvironment(environment Environment) (map[string]interface{}, error) {
	var result map[string]interface{}
	resp, err := a.codefresh.requestAPI(&requestOptions{method: "POST", path: "/api/environments-v2/argo/events", body: environment})
	if err != nil {
		return nil, err
	}

	err = a.codefresh.decodeResponseInto(resp, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *gitops) DeleteEnvironment(name string) error {
	_, err := a.codefresh.requestAPI(&requestOptions{
		method: "DELETE",
		path:   fmt.Sprintf("/api/environments-v2/%s", name),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *gitops) GetEnvironments() ([]CFEnvironment, error) {
	var result MongoCFEnvWrapper
	resp, err := a.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/environments-v2?plain=true&isEnvironment=false",
	})
	if err != nil {
		return nil, err
	}

	err = a.codefresh.decodeResponseInto(resp, &result)

	if err != nil {
		return nil, err
	}

	return result.Docs, nil
}

func (a *gitops) SendEvent(name string, props map[string]string) error {
	event := CodefreshEvent{Event: name, Props: props}

	_, err := a.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   "/api/gitops/system/events",
		body:   event,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *gitops) SendApplicationResources(resources *ApplicationResources) error {
	_, err := a.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   fmt.Sprintf("/api/gitops/resources"),
		body:   &resources,
	})
	if err != nil {
		return err
	}
	return nil
}
