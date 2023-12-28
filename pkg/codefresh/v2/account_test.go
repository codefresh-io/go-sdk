package v2

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	platmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model/platform"
	"github.com/codefresh-io/go-sdk/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func Test_account_UpdateCsdpSettings(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".cfconfig")
	authContext, _ := utils.ReadAuthContext(path, "")
	cfClient := client.NewCfClient(authContext.URL, authContext.Token, "", nil)
	type args struct {
		gitProvider      platmodel.GitProviders
		gitApiUrl        string
		sharedConfigRepo string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "should return error when failed to update csdp settings",
			args: args{
				gitProvider:      platmodel.GitProvidersGithub,
				gitApiUrl:        "https://api.github.com",
				sharedConfigRepo: "https://github.com/noam-codefresh/atgardner-shared-config",
			},
			wantErr: "failed updating csdp settings: failed",
		},
	}
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
