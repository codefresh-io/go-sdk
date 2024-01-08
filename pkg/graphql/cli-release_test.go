package graphql

import (
	"context"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/mocks"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_cliRelease_GetLatest(t *testing.T) {
	tests := []struct {
		name     string
		want     string
		wantErr  string
		beforeFn func(rt *mocks.MockRoundTripper)
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := utils.NewMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			c := &cliRelease{
				client: cfClient,
			}
			got, err := c.GetLatest(context.Background())
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("cliRelease.GetLatest() = %v, want %v", got, tt.want)
			}
		})
	}
}
