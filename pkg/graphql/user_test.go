package graphql

import (
	"context"
	"reflect"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/mocks"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_user_GetCurrent(t *testing.T) {
	tests := []struct {
		name     string
		want     *platmodel.User
		wantErr  string
		beforeFn func(rt *mocks.MockRoundTripper)
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := utils.NewMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			c := &user{
				client: cfClient,
			}
			got, err := c.GetCurrent(context.Background())
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.GetCurrent() = %v, want %v", got, tt.want)
			}
		})
	}
}
