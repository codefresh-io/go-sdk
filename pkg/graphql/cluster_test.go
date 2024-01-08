package graphql

import (
	"context"
	"reflect"
	"testing"

	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_cluster_List(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name    string
		runtime string
		want    []platmodel.Cluster
		wantErr string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cluster{
				client: cfClient,
			}
			got, err := c.List(context.Background(), tt.runtime)
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cluster.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
