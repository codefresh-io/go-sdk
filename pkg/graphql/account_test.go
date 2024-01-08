package graphql

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/mocks"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
	"github.com/codefresh-io/go-sdk/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_account_UpdateCsdpSettings(t *testing.T) {
	type args struct {
		gitProvider      platmodel.GitProviders
		gitApiUrl        string
		sharedConfigRepo string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  string
		beforeFn func(rt *mocks.MockRoundTripper)
	}{
		{
			name: "should return nil when csdp settings updated",
			args: args{
				gitProvider:      platmodel.GitProvidersBitbucket,
				gitApiUrl:        "https://bitbucket.org",
				sharedConfigRepo: "codefresh/cf-test",
			},
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(_ *http.Request) (*http.Response, error) {
					bodyReader := io.NopCloser(strings.NewReader("null"))
					res := &http.Response{
						StatusCode: 200,
						Body:       bodyReader,
					}
					return res, nil
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := utils.NewMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			c := &account{
				client: cfClient,
			}
			if err := c.UpdateCsdpSettings(context.Background(), tt.args.gitProvider, tt.args.gitApiUrl, tt.args.sharedConfigRepo); err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
