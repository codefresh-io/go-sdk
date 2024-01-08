package graphql

import (
	"context"
	"testing"

	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
	"github.com/codefresh-io/go-sdk/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func Test_account_UpdateCsdpSettings(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	type args struct {
		gitProvider      platmodel.GitProviders
		gitApiUrl        string
		sharedConfigRepo string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &account{
				client: cfClient,
			}
			if err := c.UpdateCsdpSettings(context.Background(), tt.args.gitProvider, tt.args.gitApiUrl, tt.args.sharedConfigRepo); err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
