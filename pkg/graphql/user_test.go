package graphql

import (
	"context"
	"reflect"
	"testing"

	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_user_GetCurrent(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name    string
		want    *platmodel.User
		wantErr string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
