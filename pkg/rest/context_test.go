package rest

import (
	"reflect"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/mocks"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_v1Context_GetDefaultGitContext(t *testing.T) {
	tests := []struct {
		name     string
		want     *ContextPayload
		wantErr  string
		beforeFn func(rt *mocks.MockRoundTripper)
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := utils.NewMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			c := v1Context{
				client: cfClient,
			}
			got, err := c.GetDefaultGitContext()
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("v1Context.GetDefaultGitContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_v1Context_GetGitContextByName(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		want        *ContextPayload
		wantErr     string
		beforeFn    func(rt *mocks.MockRoundTripper)
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := utils.NewMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			c := v1Context{
				client: cfClient,
			}
			got, err := c.GetGitContextByName(tt.contextName)
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("v1Context.GetGitContextByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_v1Context_GetGitContexts(t *testing.T) {
	tests := []struct {
		name     string
		want     []ContextPayload
		wantErr  string
		beforeFn func(rt *mocks.MockRoundTripper)
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := utils.NewMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			c := v1Context{
				client: cfClient,
			}
			got, err := c.GetGitContexts()
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("v1Context.GetGitContexts() = %v, want %v", got, tt.want)
			}
		})
	}
}
