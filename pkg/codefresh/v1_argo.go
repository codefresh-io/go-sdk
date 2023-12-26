package codefresh

import "fmt"

type (
	V1ArgoAPI interface {
		CreateIntegration(integration IntegrationPayloadData) error
		DeleteIntegrationByName(name string) error
		GetIntegrationByName(name string) (*IntegrationPayload, error)
		GetIntegrations() ([]*IntegrationPayload, error)
		HeartBeat(error string, version string, integration string) error
		SendResources(kind string, items interface{}, amount int, integration string) error
		UpdateIntegration(name string, integration IntegrationPayloadData) error
	}

	v1Argo struct {
		codefresh *codefresh
	}

	IntegrationItem struct {
		Amount int `json:"amount"`
	}

	IntegrationPayloadData struct {
		Name          string          `json:"name"`
		Url           string          `json:"url"`
		Clusters      IntegrationItem `json:"clusters"`
		Applications  IntegrationItem `json:"applications"`
		Repositories  IntegrationItem `json:"repositories"`
		Username      *string         `json:"username"`
		Password      *string         `json:"password"`
		Token         *string         `json:"token"`
		ClusterName   *string         `json:"clusterName"`
		ServerVersion *string         `json:"serverVersion"`
		Provider      *string         `json:"provider"`
	}

	IntegrationPayload struct {
		Type string                 `json:"type"`
		Data IntegrationPayloadData `json:"data"`
	}

	Heartbeat struct {
		Error        string `json:"error"`
		AgentVersion string `json:"agentVersion"`
	}
	AgentState struct {
		Kind  string      `json:"type"`
		Items interface{} `json:"items"`
	}
)

func (a *v1Argo) CreateIntegration(integration IntegrationPayloadData) error {
	_, err := a.codefresh.requestAPI(&requestOptions{
		path:   "/api/argo",
		method: "POST",
		body: &IntegrationPayload{
			Type: "argo-cd",
			Data: integration,
		},
	})

	return err
}

func (a *v1Argo) DeleteIntegrationByName(name string) error {
	_, err := a.codefresh.requestAPI(&requestOptions{
		method: "DELETE",
		path:   fmt.Sprintf("/api/argo/%s", name),
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *v1Argo) GetIntegrationByName(name string) (*IntegrationPayload, error) {
	var result IntegrationPayload

	resp, err := a.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   fmt.Sprintf("/api/argo/%s", name),
	})

	if err != nil {
		return nil, err
	}

	err = a.codefresh.decodeResponseInto(resp, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a *v1Argo) GetIntegrations() ([]*IntegrationPayload, error) {
	var result []*IntegrationPayload

	resp, err := a.codefresh.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/argo",
	})

	if err != nil {
		return nil, err
	}

	err = a.codefresh.decodeResponseInto(resp, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *v1Argo) HeartBeat(error string, version string, integration string) error {
	var body = Heartbeat{}

	if error != "" {
		body.Error = error
	}

	if version != "" {
		body.AgentVersion = version
	}

	_, err := a.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   fmt.Sprintf("/api/argo-agent/%s/heartbeat", integration),
		body:   body,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *v1Argo) SendResources(kind string, items interface{}, amount int, integration string) error {
	if items == nil {
		return nil
	}

	_, err := a.codefresh.requestAPI(&requestOptions{
		method: "POST",
		path:   fmt.Sprintf("/api/argo-agent/%s", integration),
		body:   &AgentState{Kind: kind, Items: items},
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *v1Argo) UpdateIntegration(name string, integration IntegrationPayloadData) error {
	_, err := a.codefresh.requestAPI(&requestOptions{
		method: "PUT",
		path:   fmt.Sprintf("/api/argo/%s", name),
		body: &IntegrationPayload{
			Type: "argo-cd",
			Data: integration,
		},
	})
	if err != nil {
		return err
	}

	return nil
}