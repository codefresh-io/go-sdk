package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type (
	ClientOptions struct {
		Token       string
		Host        string
		Client      *http.Client
		graphqlPath string
	}

	CfClient struct {
		token   string
		baseUrl *url.URL
		gqlUrl  *url.URL
		client  *http.Client
	}

	RequestOptions struct {
		Path   string
		Method string
		Query  map[string]string
		Body   interface{}
	}

	GraphqlOptions struct {
		Query     string
		Variables map[string]any
	}

	ApiError struct {
		Message    string
		StatusCode int
		StatusText string
		Body       string
	}
)

func (e *ApiError) Error() string {
	return fmt.Sprintf("[%d] - %s:\n%s", e.StatusCode, e.Message, e.Body)
}

func NewCfClient(opt *ClientOptions) *CfClient {
	var (
		graphqlPath string
		httpClient  *http.Client
	)

	baseUrl, err := url.Parse(opt.Host)
	if err != nil {
		panic(err)
	}

	if opt.graphqlPath != "" {
		graphqlPath = opt.graphqlPath
	} else {
		graphqlPath = "/2.0/api/graphql"
	}

	gqlUrl := baseUrl.JoinPath(graphqlPath)
	if opt.Client != nil {
		httpClient = opt.Client
	} else {
		httpClient = &http.Client{}
	}

	return &CfClient{
		token:   opt.Token,
		baseUrl: baseUrl,
		gqlUrl:  gqlUrl,
		client:  httpClient,
	}
}

func (c *CfClient) apiCall(ctx context.Context, baseUrl *url.URL, opt *RequestOptions) (*http.Response, error) {
	var body []byte
	finalUrl := baseUrl.JoinPath(opt.Path)
	q := finalUrl.Query()
	for k, v := range opt.Query {
		q.Set(k, v)
	}

	finalUrl.RawQuery = q.Encode()
	if opt.Body != nil {
		body, _ = json.Marshal(opt.Body)
	}

	method := http.MethodGet
	if opt.Method != "" {
		method = opt.Method
	}

	request, err := http.NewRequestWithContext(ctx, method, finalUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Authorization", c.token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("origin", c.baseUrl.Host)

	res, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return res, nil
}

func (c *CfClient) RestAPI(ctx context.Context, opt *RequestOptions) ([]byte, error) {
	res, err := c.apiCall(ctx, c.baseUrl, opt)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode >= 400 {
		return nil, &ApiError{
			Message:    "failed to make a REST API request",
			StatusCode: res.StatusCode,
			StatusText: res.Status,
			Body:       string(bytes),
		}
	}

	return bytes, nil
}

func (c *CfClient) GraphqlAPI(ctx context.Context, opt *GraphqlOptions, result any) error {
	res, err := c.apiCall(ctx, c.gqlUrl, &RequestOptions{
		Method: "POST",
		Body: map[string]any{
			"query":     opt.Query,
			"variables": opt.Variables,
		},
	})
	if err != nil {
		return err
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode >= 400 {
		return &ApiError{
			Message:    "failed to make a GraphQL API request",
			StatusCode: res.StatusCode,
			StatusText: res.Status,
			Body:       string(bytes),
		}
	}

	return json.Unmarshal(bytes, res)
}
