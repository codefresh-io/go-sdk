package rest

import (
	"reflect"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_argo_GetIntegrations(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name    string
		want    []IntegrationPayload
		wantErr string
	}{
		{
			name: "should return integrations",
			want: []IntegrationPayload{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &argo{
				client: cfClient,
			}
			got, err := a.GetIntegrations()
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("argo.GetIntegrations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argo_GetIntegrationByName(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name            string
		integrationName string
		want            *IntegrationPayload
		wantErr         string
	}{
		{
			name:            "should return integration",
			integrationName: "argocd-noam-1",
			want:            &IntegrationPayload{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &argo{
				client: cfClient,
			}
			got, err := a.GetIntegrationByName(tt.integrationName)
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("argo.GetIntegrationByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argo_DeleteIntegrationByName(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name            string
		integrationName string
		wantErr         string
	}{
		{
			name:            "should delete integration",
			integrationName: "argocd-noam-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &argo{
				client: cfClient,
			}
			if err := a.DeleteIntegrationByName(tt.integrationName); err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
