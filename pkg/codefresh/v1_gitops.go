package codefresh

import (
	"fmt"
	"time"
)

type (
	V1GitopsAPI interface {
		CreateEnvironment(name string, project string, application string, integration string) error
		DeleteEnvironment(name string) error
		GetEnvironments() ([]CFEnvironment, error)
		SendApplicationResources(resources *ApplicationResources) error
		SendEnvironment(environment Environment) (map[string]interface{}, error)
		SendEvent(name string, props map[string]string) error
	}

	v1Gitops struct {
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
		Time    *time.Time `json:"time,omitempty"`
		Message *string    `json:"message"`
		Avatar  *string    `json:"avatar"`
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

	GitopsUser struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}

	Gitops struct {
		Comitters []GitopsUser `json:"comitters"`
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
		Namespace    string                `json:"namespace"`
		Server       string                `json:"server"`
		Context      *string               `json:"context"`
	}

	EnvironmentActivity struct {
		Name         string                `json:"name"`
		TargetImages []string              `json:"targetImages"`
		Status       string                `json:"status"`
		LiveImages   []string              `json:"liveImages"`
		ReplicaSet   EnvironmentActivityRS `json:"replicaSet"`
	}

	ApplicationResources struct {
		Name      string      `json:"name,omitempty"`
		HistoryId int64       `json:"historyId"`
		Revision  string      `json:"revision,omitempty"`
		Resources interface{} `json:"resources"`
		Context   *string     `json:"context"`
	}
)

func (a *v1Gitops) CreateEnvironment(name string, project string, application string, integration string) error {
	resp, err := a.codefresh.requestAPI(&requestOptions{
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

	defer resp.Body.Close()
	return nil
}

func (a *v1Gitops) DeleteEnvironment(name string) error {
	resp, err := a.codefresh.requestAPI(&requestOptions{
		method: "DELETE",
		path:   fmt.Sprintf("/api/environments-v2/%s", name),
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (a *v1Gitops) GetEnvironments() ([]CFEnvironment, error) {
	resp, err := a.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/environments-v2?plain=true&isEnvironment=false",
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	result := &MongoCFEnvWrapper{}
	err = a.codefresh.decodeResponseInto(resp, result)
	if err != nil {
		return nil, err
	}

	return result.Docs, nil
}

func (a *v1Gitops) SendApplicationResources(resources *ApplicationResources) error {
	resp, err := a.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   fmt.Sprintf("/api/gitops/resources"),
		body:   &resources,
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (a *v1Gitops) SendEnvironment(environment Environment) (map[string]interface{}, error) {
	resp, err := a.codefresh.requestAPI(&requestOptions{method: "POST", path: "/api/environments-v2/argo/events", body: environment})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	result := make(map[string]interface{})
	err = a.codefresh.decodeResponseInto(resp, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *v1Gitops) SendEvent(name string, props map[string]string) error {
	event := CodefreshEvent{Event: name, Props: props}

	resp, err := a.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   "/api/gitops/system/events",
		body:   event,
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}
