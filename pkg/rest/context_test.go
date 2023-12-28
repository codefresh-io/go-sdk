package rest

import (
	"reflect"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_v1Context_GetDefaultGitContext(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name    string
		want    *ContextPayload
		wantErr string
	}{
		{
			name: "should return default git context",
			want: &ContextPayload{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	cfClient := utils.NewClientFromCurrentContext()
	type args struct {
	}
	tests := []struct {
		name        string
		contextName string
		want        *ContextPayload
		wantErr     string
	}{
		{
			name:        "should return git context by name",
			contextName: "github",
			want:        &ContextPayload{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name    string
		want    []ContextPayload
		wantErr string
	}{
		{
			name: "should return git contexts",
			want: []ContextPayload{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
