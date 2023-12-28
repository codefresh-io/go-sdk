package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type (
	GraphqlResponse interface {
		HasErrors() bool
	}

	ClientOptions struct {
		Token       string
		Host        string
		Client      *http.Client
		GraphqlPath string
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
		Query  map[string]any
		Body   any
	}

	ApiError struct {
		response *http.Response
	}

	GraphqlError struct {
		Message    string
		Extensions any
	}

	GraphqlErrorResponse struct {
		errors             []GraphqlError
		concatenatedErrors string
	}

	GraphqlBaseResponse struct {
		Errors []GraphqlError
	}

	GraphqlVoidResponse struct{}
)

func (e *ApiError) Error() string {
	return e.response.Status
}

func NewCfClient(host, token, graphqlPath string, httpClient *http.Client) *CfClient {
	baseUrl, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	if graphqlPath == "" {
		graphqlPath = "/2.0/api/graphql"
	}

	gqlUrl := baseUrl.JoinPath(graphqlPath)
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &CfClient{
		baseUrl: baseUrl,
		token:   token,
		gqlUrl:  gqlUrl,
		client:  httpClient,
	}
}

func (c *CfClient) AppProxyClient(host string, insecure bool) *CfClient {
	httpClient := &http.Client{}
	httpClient.Timeout = c.client.Timeout
	if insecure {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		httpClient.Transport = customTransport
	}

	return NewCfClient(host, c.token, "/app-proxy/api/graphql", httpClient)
}

func (c *CfClient) RestAPI(ctx context.Context, opt *RequestOptions) ([]byte, error) {
	res, err := c.apiCall(ctx, c.baseUrl, opt)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response Body: %w", err)
	}

	return bytes, nil
}

func (c *CfClient) GraphqlAPI(ctx context.Context, query string, variables any, result any) error {
	body := map[string]any{
		"query":     query,
		"variables": variables,
	}
	res, err := c.apiCall(ctx, c.gqlUrl, &RequestOptions{
		Method: "POST",
		Body:   body,
	})
	if err != nil {
		return err
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response Body: %w", err)
	}

	err = json.Unmarshal(bytes, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response Body: %w", err)
	}

	return nil
}

func (c *CfClient) apiCall(ctx context.Context, baseUrl *url.URL, opt *RequestOptions) (*http.Response, error) {
	var body []byte
	finalUrl := baseUrl.JoinPath(opt.Path)
	q := finalUrl.Query()
	setQueryParams(q, opt.Query)
	finalUrl.RawQuery = q.Encode()

	if opt.Body != nil {
		body, _ = json.Marshal(opt.Body)
	}

	method := http.MethodGet
	if opt.Method != "" {
		method = opt.Method
	}

	if ctx == nil {
		ctx = context.Background()
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

	if res.StatusCode >= http.StatusBadRequest {
		return nil, &ApiError{response: res}
	}

	return res, nil
}

func (e GraphqlErrorResponse) Error() string {
	if e.concatenatedErrors != "" {
		return e.concatenatedErrors
	}

	var sb strings.Builder
	for _, err := range e.errors {
		sb.WriteString(fmt.Sprintln(err.Message))
	}

	e.concatenatedErrors = sb.String()
	return e.concatenatedErrors
}

func (e GraphqlBaseResponse) HasErrors() bool {
	return len(e.Errors) > 0
}

func GraphqlAPI[T any](ctx context.Context, client *CfClient, query string, variables any) (T, error) {
	var (
		wrapper struct {
			Data   map[string]T   `json:"data,omitempty"`
			Errors []GraphqlError `json:"errors,omitempty"`
		}
		result T
	)

	err := client.GraphqlAPI(ctx, query, variables, &wrapper)
	if err != nil {
		return result, err
	}

	// we assume there is only a single data key in the result (= a single query in the request)
	for k := range wrapper.Data {
		result = wrapper.Data[k]
		break
	}

	if wrapper.Errors != nil {
		err = &GraphqlErrorResponse{errors: wrapper.Errors}
	}

	return result, err
}

func setQueryParams(q url.Values, query map[string]any) error {
	for k, v := range query {
		if str, ok := v.(string); ok {
			q.Set(k, str)
		} else if arr, ok := v.([]string); ok {
			for _, item := range arr {
				q.Add(k, item)
			}
		} else {
			return fmt.Errorf("invalid query param type: %T", v)
		}
	}

	return nil
}
