package graphql

import (
	"context"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_cliRelease_GetLatest(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name    string
		want    string
		wantErr string
	}{
		{
			name: "should return latest cli release",
			want: "0.1.55",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
