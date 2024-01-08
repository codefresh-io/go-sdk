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

func Test_component_List(t *testing.T) {
	tests := []struct {
		name        string
		runtimeName string
		want        []platmodel.Component
		wantErr     string
		beforeFn    func(rt *mocks.MockRoundTripper)
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := utils.NewMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			c := &component{
				client: cfClient,
			}
			got, err := c.List(context.Background(), tt.runtimeName)
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("component.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
