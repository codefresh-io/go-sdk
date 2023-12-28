package rest

import (
	"reflect"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_cluster_GetAccountClusters(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name    string
		want    []ClusterMinified
		wantErr string
	}{
		{
			name: "should return clusters",
			want: []ClusterMinified{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &cluster{
				client: cfClient,
			}
			got, err := p.GetAccountClusters()
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cluster.GetAccountClusters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cluster_GetClusterCredentialsByAccountId(t *testing.T) {
	cfClient := utils.NewClientFromCurrentContext()
	tests := []struct {
		name     string
		selector string
		want     *Cluster
		wantErr  string
	}{
		{
			name:     "should return cluster credentials",
			selector: "noam",
			want:     &Cluster{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &cluster{
				client: cfClient,
			}
			got, err := p.GetClusterCredentialsByAccountId(tt.selector)
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cluster.GetClusterCredentialsByAccountId() = %v, want %v", got, tt.want)
			}
		})
	}
}
