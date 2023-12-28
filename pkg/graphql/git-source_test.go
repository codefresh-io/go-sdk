package graphql

import (
	"context"
	"reflect"
	"testing"

	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_gitSource_List(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name        string
		runtimeName string
		want        []platmodel.GitSource
		wantErr     string
	}{
		{
			name:        "should return git sources",
			runtimeName: "atgardner-noam",
			want: []platmodel.GitSource{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &gitSource{
				client: cfClient,
			}
			got, err := c.List(context.Background(), tt.runtimeName)
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gitSource.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
