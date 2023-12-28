package graphql

import (
	"context"
	"reflect"
	"testing"

	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_component_List(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name        string
		runtimeName string
		want        []platmodel.Component
		wantErr     string
	}{
		{
			name:        "should return list of components",
			runtimeName: "atgardner-noam",
			want:        []platmodel.Component{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
