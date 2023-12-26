package codefresh

// import (
// 	"codefresh"
// )

type (
	ArgoAPI interface {
		// CreateIntegration(integration IntegrationPayloadData) error
		// UpdateIntegration(name string, integration IntegrationPayloadData) error
		// GetIntegrations() ([]*IntegrationPayload, error)
		// GetIntegrationByName(name string) (*IntegrationPayload, error)
		// DeleteIntegrationByName(name string) error
		// HeartBeat(error string, version string, integration string) error
		// SendResources(kind string, items interface{}, amount int, integration string) error
	}

	argo struct {
		codefresh cfImpl
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

func newArgoAPI(c any) ArgoAPI {
	return &argo{}
}

// func (a *argo) CreateIntegration(integration IntegrationPayloadData) error {
// 	_, err := a.codefresh.requestAPI(&RequestOptions{
// 		path:   "/api/argo",
// 		method: "POST",
// 		body: &IntegrationPayload{
// 			Type: "argo-cd",
// 			Data: integration,
// 		},
// 	})

// 	return err
// }

// func (a *argo) UpdateIntegration(name string, integration IntegrationPayloadData) error {
// 	_, err := a.codefresh.requestAPI(&RequestOptions{
// 		method: "PUT",
// 		path:   fmt.Sprintf("/api/argo/%s", name),
// 		body: &IntegrationPayload{
// 			Type: "argo-cd",
// 			Data: integration,
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (a *argo) GetIntegrations() ([]*IntegrationPayload, error) {
// 	var result []*IntegrationPayload

// 	resp, err := a.codefresh.requestAPI(&RequestOptions{
// 		method: "GET",
// 		path:   "/api/argo",
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	err = a.codefresh.decodeResponseInto(resp, &result)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func (a *argo) GetIntegrationByName(name string) (*IntegrationPayload, error) {
// 	var result IntegrationPayload

// 	resp, err := a.codefresh.requestAPI(&RequestOptions{
// 		method: "GET",
// 		path:   fmt.Sprintf("/api/argo/%s", name),
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	err = a.codefresh.decodeResponseInto(resp, &result)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &result, nil
// }

// func (a *argo) DeleteIntegrationByName(name string) error {
// 	_, err := a.codefresh.requestAPI(&RequestOptions{
// 		method: "DELETE",
// 		path:   fmt.Sprintf("/api/argo/%s", name),
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (a *argo) HeartBeat(error string, version string, integration string) error {
// 	var body = Heartbeat{}

// 	if error != "" {
// 		body.Error = error
// 	}

// 	if version != "" {
// 		body.AgentVersion = version
// 	}

// 	_, err := a.codefresh.requestAPI(&RequestOptions{
// 		method: "POST",
// 		path:   fmt.Sprintf("/api/argo-agent/%s/heartbeat", integration),
// 		body:   body,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (a *argo) SendResources(kind string, items interface{}, amount int, integration string) error {
// 	if items == nil {
// 		return nil
// 	}

// 	_, err := a.codefresh.requestAPI(&RequestOptions{
// 		method: "POST",
// 		path:   fmt.Sprintf("/api/argo-agent/%s", integration),
// 		body:   &AgentState{Kind: kind, Items: items},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
