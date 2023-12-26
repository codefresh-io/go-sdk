package codefresh

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/go-querystring/query"
)

//go:generate mockery -name Codefresh -filename codefresh.go

//go:generate mockery -name UsersAPI -filename users.go

type (
	Codefresh interface {
		V1() V1API
		V2() V2API
	}

	ClientOptions struct {
		Token       string
		Debug       bool
		Host        string
		Client      *http.Client
		graphqlPath string
	}

	codefresh struct {
		token       string
		host        string
		graphqlPath string
		client      *http.Client
	}
	requestOptions struct {
		path   string
		method string
		body   interface{}
		qs     interface{}
	}

	graphqlError struct {
		Message    string
		Extensions interface{}
	}

	graphqlErrorResponse struct {
		errors             []graphqlError
		concatenatedErrors string
	}

	graphqlVoidResponse struct {
		Errors []graphqlError
	}
)

func New(opt *ClientOptions) Codefresh {
	return newClient(opt)
}

func (c *codefresh) V1() V1API {
	return newV1Client(c)
}

func (c *codefresh) V2() V2API {
	return newV2Client(c)
}

func (c *codefresh) requestAPI(opt *requestOptions) (*http.Response, error) {
	return c.requestAPIWithContext(context.Background(), opt)
}

func (c *codefresh) requestAPIWithContext(ctx context.Context, opt *requestOptions) (*http.Response, error) {
	var body []byte
	finalURL := fmt.Sprintf("%s%s", c.host, opt.path)
	if opt.qs != nil {
		finalURL += toQS(opt.qs)
	}
	if opt.body != nil {
		body, _ = json.Marshal(opt.body)
	}
	request, err := http.NewRequestWithContext(ctx, opt.method, finalURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", c.token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("origin", c.host)

	response, err := c.client.Do(request)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (c *codefresh) graphqlAPI(ctx context.Context, body map[string]interface{}, res interface{}) error {
	response, err := c.requestAPIWithContext(ctx, &requestOptions{
		method: "POST",
		path:   c.graphqlPath,
		body:   body,
	})
	if err != nil {
		return fmt.Errorf("the HTTP request failed: %w", err)
	}
	defer response.Body.Close()

	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		return errors.New(response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read from response body: %w", err)
	}

	return json.Unmarshal(data, res)
}

func buildQSFromMap(qs map[string]string) string {
	var arr = []string{}
	for k, v := range qs {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v))
	}
	return "?" + strings.Join(arr, "&")
}

func toQS(qs interface{}) string {
	v, _ := query.Values(qs)
	qsStr := v.Encode()
	if qsStr != "" {
		return "?" + qsStr
	}
	var qsMap map[string]string
	rs, _ := json.Marshal(qs)
	err := json.Unmarshal(rs, &qsMap)
	if err != nil {
		return ""
	}
	return buildQSFromMap(qsMap)
}

func (c *codefresh) decodeResponseInto(resp *http.Response, target interface{}) error {
	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *codefresh) getBodyAsString(resp *http.Response) (string, error) {
	body, err := c.getBodyAsBytes(resp)
	return string(body), err
}

func (c *codefresh) getBodyAsBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func newClient(opt *ClientOptions) *codefresh {
	httpClient := &http.Client{}
	if opt.Client != nil {
		httpClient = opt.Client
	}

	graphqlPath := "/2.0/api/graphql"
	if opt.graphqlPath != "" {
		graphqlPath = opt.graphqlPath
	}

	re := regexp.MustCompile("/$")

	if re.FindString(opt.Host) != "" {
		if len(opt.Host) > 1 {
			opt.Host = opt.Host[:len(opt.Host)-1]
		}
	}

	return &codefresh{
		host:        opt.Host,
		token:       opt.Token,
		graphqlPath: graphqlPath,
		client:      httpClient,
	}
}

func (e graphqlErrorResponse) Error() string {
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
