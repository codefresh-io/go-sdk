package rest

import (
	"encoding/json"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/client"
)

type (
	ArgoAPI interface {
		CreateIntegration(integration IntegrationPayloadData) error
		DeleteIntegrationByName(name string) error
		GetIntegrationByName(name string) (*IntegrationPayload, error)
		GetIntegrations() ([]IntegrationPayload, error)
		HeartBeat(error string, version string, integration string) error
		SendResources(kind string, items any, amount int, integration string) error
		UpdateIntegration(name string, integration IntegrationPayloadData) error
	}

	argo struct {
		client *client.CfClient
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
		Kind  string `json:"type"`
		Items any    `json:"items"`
	}
)

func (a *argo) CreateIntegration(integration IntegrationPayloadData) error {
	_, err := a.client.RestAPI(nil, &client.RequestOptions{
		Path:   "/api/argo",
		Method: "POST",
		Body: &IntegrationPayload{
			Type: "argo-cd",
			Data: integration,
		},
	})
	if err != nil {
		return fmt.Errorf("failed creating an argo integration: %w", err)
	}

	return nil
}

func (a *argo) DeleteIntegrationByName(name string) error {
	_, err := a.client.RestAPI(nil, &client.RequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("/api/argo/%s", name),
	})
	if err != nil {
		return fmt.Errorf("failed deleting an argo integration: %w", err)
	}

	return nil
}

func (a *argo) GetIntegrationByName(name string) (*IntegrationPayload, error) {
	res, err := a.client.RestAPI(nil, &client.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/argo/%s", name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting an argo integration: %w", err)
	}

	result := &IntegrationPayload{}
	return result, json.Unmarshal(res, result)
}

func (a *argo) GetIntegrations() ([]IntegrationPayload, error) {
	res, err := a.client.RestAPI(nil, &client.RequestOptions{
		Method: "GET",
		Path:   "/api/argo",
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting argo integration list: %w", err)
	}

	result := make([]IntegrationPayload, 0)
	return result, json.Unmarshal(res, &result)
}

func (a *argo) HeartBeat(error string, version string, integration string) error {
	var body = Heartbeat{}
	if error != "" {
		body.Error = error
	}

	if version != "" {
		body.AgentVersion = version
	}

	_, err := a.client.RestAPI(nil, &client.RequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/api/argo-agent/%s/heartbeat", integration),
		Body:   body,
	})
	if err != nil {
		return fmt.Errorf("failed sending argo heartbeat: %w", err)
	}

	return nil
}

func (a *argo) SendResources(kind string, items any, amount int, integration string) error {
	if items == nil {
		return nil
	}

	_, err := a.client.RestAPI(nil, &client.RequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/api/argo-agent/%s", integration),
		Body:   &AgentState{Kind: kind, Items: items},
	})
	if err != nil {
		return fmt.Errorf("failed sending argo resources: %w", err)
	}

	return nil
}

func (a *argo) UpdateIntegration(name string, integration IntegrationPayloadData) error {
	_, err := a.client.RestAPI(nil, &client.RequestOptions{
		Method: "PUT",
		Path:   fmt.Sprintf("/api/argo/%s", name),
		Body: &IntegrationPayload{
			Type: "argo-cd",
			Data: integration,
		},
	})
	if err != nil {
		return fmt.Errorf("failed updating an argo integration: %w", err)
	}

	return nil
}
