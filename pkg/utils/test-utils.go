package utils

import (
	"net/http"
	"testing"

	"github.com/codefresh-io/go-sdk/pkg/client"
	"github.com/codefresh-io/go-sdk/pkg/mocks"
)

func NewMockClient(t *testing.T) (*client.CfClient, *mocks.MockRoundTripper) {
	mockRT := mocks.NewMockRoundTripper(t)
	cfClient := client.NewCfClient("https://some.host", "some-token", "grpahql-path", &http.Client{
		Transport: mockRT,
	})
	return cfClient, mockRT
}
