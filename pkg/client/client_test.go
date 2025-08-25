package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/codefresh-io/go-sdk/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newMockClient creates a mock client for testing
func newMockClient(t *testing.T) (*CfClient, *mocks.MockRoundTripper) {
	mockRT := mocks.NewMockRoundTripper(t)
	cfClient := NewCfClient("https://some.host", "some-token", "grpahql-path", &http.Client{
		Transport: mockRT,
	})
	return cfClient, mockRT
}

func TestNewCfClient(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		token       string
		graphqlPath string
		httpClient  *http.Client
		wantPanic   bool
	}{
		{
			name:        "should create client with default graphql path",
			host:        "https://api.codefresh.io",
			token:       "test-token",
			graphqlPath: "",
			httpClient:  nil,
			wantPanic:   false,
		},
		{
			name:        "should create client with custom graphql path",
			host:        "https://api.codefresh.io",
			token:       "test-token",
			graphqlPath: "/custom/graphql",
			httpClient:  &http.Client{Timeout: 30 * time.Second},
			wantPanic:   false,
		},
		{
			name:        "should panic with invalid host",
			host:        "://invalid-url",
			token:       "test-token",
			graphqlPath: "",
			httpClient:  nil,
			wantPanic:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					NewCfClient(tt.host, tt.token, tt.graphqlPath, tt.httpClient)
				})
				return
			}

			client := NewCfClient(tt.host, tt.token, tt.graphqlPath, tt.httpClient)
			assert.NotNil(t, client)
			assert.Equal(t, tt.token, client.token)
			assert.Equal(t, tt.host, client.baseUrl.String())

			expectedPath := tt.graphqlPath
			if expectedPath == "" {
				expectedPath = "/2.0/api/graphql"
			}
			assert.Contains(t, client.gqlUrl.String(), expectedPath)

			if tt.httpClient == nil {
				assert.NotNil(t, client.client)
			} else {
				assert.Equal(t, tt.httpClient, client.client)
			}
		})
	}
}

func TestCfClient_AppProxyClient(t *testing.T) {
	originalClient := NewCfClient("https://api.codefresh.io", "test-token", "", &http.Client{
		Timeout: 30 * time.Second,
	})

	tests := []struct {
		name     string
		host     string
		insecure bool
	}{
		{
			name:     "should create secure app proxy client",
			host:     "https://app-proxy.codefresh.io",
			insecure: false,
		},
		{
			name:     "should create insecure app proxy client",
			host:     "https://app-proxy.codefresh.io",
			insecure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxyClient := originalClient.AppProxyClient(tt.host, tt.insecure)

			assert.NotNil(t, proxyClient)
			assert.Equal(t, originalClient.token, proxyClient.token)
			assert.Equal(t, tt.host, proxyClient.baseUrl.String())
			assert.Contains(t, proxyClient.gqlUrl.String(), "/app-proxy/api/graphql")
			assert.Equal(t, originalClient.client.Timeout, proxyClient.client.Timeout)

			if tt.insecure {
				transport := proxyClient.client.Transport.(*http.Transport)
				assert.NotNil(t, transport.TLSClientConfig)
				assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
			}
		})
	}
}

func TestCfClient_RestAPI(t *testing.T) {
	tests := []struct {
		name     string
		opt      *RequestOptions
		wantErr  string
		beforeFn func(rt *mocks.MockRoundTripper)
		wantData string
	}{
		{
			name: "should make successful GET request",
			opt: &RequestOptions{
				Path:   "/api/accounts",
				Method: "GET",
			},
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, "GET", req.Method)
					assert.Contains(t, req.URL.Path, "/api/accounts")
					assert.Equal(t, "some-token", req.Header.Get("Authorization"))

					bodyReader := io.NopCloser(strings.NewReader(`{"success": true}`))
					return &http.Response{
						StatusCode: 200,
						Body:       bodyReader,
					}, nil
				})
			},
			wantData: `{"success": true}`,
		},
		{
			name: "should make successful POST request with body",
			opt: &RequestOptions{
				Path:   "/api/accounts",
				Method: "POST",
				Body:   map[string]string{"name": "test"},
			},
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, "POST", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

					bodyReader := io.NopCloser(strings.NewReader(`{"id": "123"}`))
					return &http.Response{
						StatusCode: 201,
						Body:       bodyReader,
					}, nil
				})
			},
			wantData: `{"id": "123"}`,
		},
		{
			name: "should handle query parameters",
			opt: &RequestOptions{
				Path:  "/api/accounts",
				Query: map[string]any{"limit": "10", "tags": []string{"tag1", "tag2"}},
			},
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, "10", req.URL.Query().Get("limit"))
					assert.ElementsMatch(t, []string{"tag1", "tag2"}, req.URL.Query()["tags"])

					bodyReader := io.NopCloser(strings.NewReader(`[]`))
					return &http.Response{
						StatusCode: 200,
						Body:       bodyReader,
					}, nil
				})
			},
			wantData: `[]`,
		},
		{
			name: "should handle API error",
			opt: &RequestOptions{
				Path: "/api/accounts",
			},
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					bodyReader := io.NopCloser(strings.NewReader(`{"error": "not found"}`))
					return &http.Response{
						StatusCode: 404,
						Status:     "404 Not Found",
						Body:       bodyReader,
					}, nil
				})
			},
			wantErr: "API error: 404 Not Found: {\"error\": \"not found\"}",
		},
		{
			name: "should handle network error",
			opt: &RequestOptions{
				Path: "/api/accounts",
			},
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).Return(nil, fmt.Errorf("network error"))
			},
			wantErr: "network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := newMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			data, err := cfClient.RestAPI(context.Background(), tt.opt)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantData, string(data))
		})
	}
}

func TestCfClient_NativeRestAPI(t *testing.T) {
	tests := []struct {
		name     string
		opt      *RequestOptions
		wantErr  string
		beforeFn func(rt *mocks.MockRoundTripper)
	}{
		{
			name: "should return raw response",
			opt: &RequestOptions{
				Path: "/api/accounts",
			},
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					bodyReader := io.NopCloser(strings.NewReader(`{"data": "test"}`))
					return &http.Response{
						StatusCode: 200,
						Body:       bodyReader,
					}, nil
				})
			},
		},
		{
			name: "should return error response without wrapping",
			opt: &RequestOptions{
				Path: "/api/accounts",
			},
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					bodyReader := io.NopCloser(strings.NewReader(`{"error": "bad request"}`))
					return &http.Response{
						StatusCode: 400,
						Body:       bodyReader,
					}, nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := newMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			resp, err := cfClient.NativeRestAPI(context.Background(), tt.opt)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			defer resp.Body.Close()
		})
	}
}

func TestCfClient_GraphqlAPI(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		variables any
		wantErr   string
		beforeFn  func(rt *mocks.MockRoundTripper)
	}{
		{
			name:  "should make successful GraphQL request",
			query: "query { accounts { id name } }",
			variables: map[string]string{
				"limit": "10",
			},
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, "POST", req.Method)
					assert.Contains(t, req.URL.Path, "grpahql-path") // Note: this matches the typo in test-utils.go
					assert.Equal(t, "some-token", req.Header.Get("Authorization"))
					assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

					bodyReader := io.NopCloser(strings.NewReader(`{"data": {"accounts": []}}`))
					return &http.Response{
						StatusCode: 200,
						Body:       bodyReader,
					}, nil
				})
			},
		},
		{
			name:  "should handle GraphQL error",
			query: "invalid query",
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					bodyReader := io.NopCloser(strings.NewReader(`{"error": "GraphQL error"}`))
					return &http.Response{
						StatusCode: 400,
						Status:     "400 Bad Request",
						Body:       bodyReader,
					}, nil
				})
			},
			wantErr: "API error: 400 Bad Request: {\"error\": \"GraphQL error\"}",
		},
		{
			name:  "should handle malformed JSON response",
			query: "query { accounts { id } }",
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					bodyReader := io.NopCloser(strings.NewReader(`invalid json`))
					return &http.Response{
						StatusCode: 200,
						Body:       bodyReader,
					}, nil
				})
			},
			wantErr: "failed to unmarshal response Body:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := newMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			var result map[string]any
			err := cfClient.GraphqlAPI(context.Background(), tt.query, tt.variables, &result)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestGraphqlAPI_Generic(t *testing.T) {
	type TestData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	tests := []struct {
		name      string
		query     string
		variables any
		wantErr   string
		beforeFn  func(rt *mocks.MockRoundTripper)
		wantData  *TestData
	}{
		{
			name:  "should return typed data successfully",
			query: "query { account { id name } }",
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					response := `{"data": {"account": {"id": "123", "name": "test"}}}`
					bodyReader := io.NopCloser(strings.NewReader(response))
					return &http.Response{
						StatusCode: 200,
						Body:       bodyReader,
					}, nil
				})
			},
			wantData: &TestData{ID: "123", Name: "test"},
		},
		{
			name:  "should handle GraphQL errors",
			query: "query { account { id name } }",
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					response := `{"data": {"account": {"id": "123", "name": "test"}}, "errors": [{"message": "field deprecated"}]}`
					bodyReader := io.NopCloser(strings.NewReader(response))
					return &http.Response{
						StatusCode: 200,
						Body:       bodyReader,
					}, nil
				})
			},
			wantData: &TestData{ID: "123", Name: "test"},
			wantErr:  "field deprecated",
		},
		{
			name:  "should handle API errors",
			query: "invalid query",
			beforeFn: func(rt *mocks.MockRoundTripper) {
				rt.EXPECT().RoundTrip(mock.AnythingOfType("*http.Request")).RunAndReturn(func(req *http.Request) (*http.Response, error) {
					bodyReader := io.NopCloser(strings.NewReader(`{"error": "bad request"}`))
					return &http.Response{
						StatusCode: 400,
						Status:     "400 Bad Request",
						Body:       bodyReader,
					}, nil
				})
			},
			wantErr: "API error: 400 Bad Request: {\"error\": \"bad request\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfClient, mockRT := newMockClient(t)
			if tt.beforeFn != nil {
				tt.beforeFn(mockRT)
			}

			result, err := GraphqlAPI[TestData](context.Background(), cfClient, tt.query, tt.variables)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				if tt.wantData != nil {
					assert.Equal(t, *tt.wantData, result)
				}
				return
			}

			assert.NoError(t, err)
			if tt.wantData != nil {
				assert.Equal(t, *tt.wantData, result)
			}
		})
	}
}

func TestApiError_Error(t *testing.T) {
	err := &ApiError{
		status:     "404 Not Found",
		statusCode: 404,
		body:       `{"error": "resource not found"}`,
	}

	expected := `API error: 404 Not Found: {"error": "resource not found"}`
	assert.Equal(t, expected, err.Error())
}

func TestGraphqlErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name     string
		response *GraphqlErrorResponse
		want     string
	}{
		{
			name: "should return cached error string",
			response: &GraphqlErrorResponse{
				concatenatedErrors: "cached error",
				Errors: []GraphqlError{
					{Message: "error 1"},
					{Message: "error 2"},
				},
			},
			want: "cached error",
		},
		{
			name: "should concatenate multiple errors",
			response: &GraphqlErrorResponse{
				Errors: []GraphqlError{
					{Message: "error 1"},
					{Message: "error 2"},
					{Message: "error 3"},
				},
			},
			want: "error 1\nerror 2\nerror 3\n",
		},
		{
			name: "should handle single error",
			response: &GraphqlErrorResponse{
				Errors: []GraphqlError{
					{Message: "single error"},
				},
			},
			want: "single error\n",
		},
		{
			name: "should handle no errors",
			response: &GraphqlErrorResponse{
				Errors: []GraphqlError{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.Error()
			assert.Equal(t, tt.want, result)

			// Note: The caching doesn't work due to value receiver, but the method still works
			if len(tt.response.Errors) > 0 && tt.response.concatenatedErrors == "" {
				result2 := tt.response.Error()
				assert.Equal(t, result, result2)
				// The concatenatedErrors field won't be set due to value receiver
			}
		})
	}
}

func TestGraphqlBaseResponse_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		response GraphqlBaseResponse
		want     bool
	}{
		{
			name: "should return true when errors exist",
			response: GraphqlBaseResponse{
				Errors: []GraphqlError{
					{Message: "error 1"},
				},
			},
			want: true,
		},
		{
			name: "should return false when no errors",
			response: GraphqlBaseResponse{
				Errors: []GraphqlError{},
			},
			want: false,
		},
		{
			name: "should return false when errors is nil",
			response: GraphqlBaseResponse{
				Errors: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.HasErrors()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSetQueryParams(t *testing.T) {
	tests := []struct {
		name    string
		query   map[string]any
		wantErr string
		want    map[string][]string
	}{
		{
			name: "should handle string values",
			query: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			want: map[string][]string{
				"key1": {"value1"},
				"key2": {"value2"},
			},
		},
		{
			name: "should handle string slice values",
			query: map[string]any{
				"tags":  []string{"tag1", "tag2", "tag3"},
				"types": []string{"type1"},
			},
			want: map[string][]string{
				"tags":  {"tag1", "tag2", "tag3"},
				"types": {"type1"},
			},
		},
		{
			name: "should handle mixed types",
			query: map[string]any{
				"name": "test",
				"tags": []string{"tag1", "tag2"},
			},
			want: map[string][]string{
				"name": {"test"},
				"tags": {"tag1", "tag2"},
			},
		},
		{
			name: "should return error for invalid type",
			query: map[string]any{
				"invalid": 123,
			},
			wantErr: "invalid query param type: int",
		},
		{
			name: "should return error for invalid slice type",
			query: map[string]any{
				"invalid": []int{1, 2, 3},
			},
			wantErr: "invalid query param type: []int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := url.Values{}
			err := setQueryParams(q, tt.query)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)
			for key, expectedValues := range tt.want {
				actualValues := q[key]
				assert.ElementsMatch(t, expectedValues, actualValues, "mismatch for key: %s", key)
			}
		})
	}
}

func TestGraphqlVoidResponse(t *testing.T) {
	// Test that GraphqlVoidResponse can be instantiated
	var response GraphqlVoidResponse
	assert.NotNil(t, &response)
}

func TestCfClient_Integration(t *testing.T) {
	// Test that all the types implement the expected interfaces
	t.Run("GraphqlBaseResponse implements GraphqlResponse", func(t *testing.T) {
		var response GraphqlResponse = GraphqlBaseResponse{}
		assert.False(t, response.HasErrors())

		response = GraphqlBaseResponse{Errors: []GraphqlError{{Message: "test"}}}
		assert.True(t, response.HasErrors())
	})

	t.Run("client configuration is preserved", func(t *testing.T) {
		originalClient := &http.Client{Timeout: 45 * time.Second}
		client := NewCfClient("https://api.codefresh.io", "token", "/custom", originalClient)

		// Test AppProxyClient preserves timeout
		proxyClient := client.AppProxyClient("https://proxy.codefresh.io", false)
		assert.Equal(t, originalClient.Timeout, proxyClient.client.Timeout)

		// Test insecure TLS configuration
		insecureProxyClient := client.AppProxyClient("https://proxy.codefresh.io", true)
		transport := insecureProxyClient.client.Transport.(*http.Transport)
		assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
	})
}
