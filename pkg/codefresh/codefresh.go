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

	v1 "codefresh/v1"
	"github.com/google/go-querystring/query"
)

//go:generate mockery -name Codefresh -filename codefresh.go

//go:generate mockery -name UsersAPI -filename users.go

type (
	Codefresh interface {
		V1() V1API
		// V2() V2API
	}

	V1API interface {
		// Argo() ArgoAPI
		// Clusters() IClusterAPI
		// Contexts() IContextAPI
		// Gitops() GitopsAPI
		// Pipelines() IPipelineAPI
		// Progresses() IProgressAPI
		// RuntimeEnvironments() IRuntimeEnvironmentAPI
		// Tokens() ITokenAPI
		// Users() UsersAPI
		// Workflows() IWorkflowAPI
	}

	V2API interface {
		// AppProxy(ctx context.Context, runtime string, insecure bool) (AppProxyAPI, error)
		// AccountV2() IAccountV2API
		// CliReleases() ICliReleasesAPI
		// Cluster() IClusterV2API
		// Component() IComponentAPI
		// GitSource() IGitSourceAPI
		// Pipeline() IPipelineV2API
		// Runtime() IRuntimeAPI
		// UsersV2() IUsersV2API
		// Workflow() IWorkflowV2API
	}

	AppProxyAPI interface {
		// AppProxyClusters() IAppProxyClustersAPI
		// AppProxyGitSources() IAppProxyGitSourcesAPI
		// AppProxyIsc() IAppProxyIscAPI
		// GitIntegrations() IAppProxyGitIntegrationsAPI
		// VersionInfo() IAppProxyVersionInfoAPI
	}

	graphqlVoidResponse struct {
		Errors []graphqlError
	}

	// Options
	ClientOptions struct {
		Token       string
		Debug       bool
		Host        string
		Client      *http.Client
		graphqlPath string
	}

	cfImpl struct {
		token       string
		host        string
		graphqlPath string
		client      *http.Client
	}

	v1Impl struct {
		codefresh *cfImpl
	}

	v2Impl struct {
		codefresh *cfImpl
	}

	apImpl struct {
		codefresh *cfImpl
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
)

func New(opt *ClientOptions) Codefresh {
	return newClient(opt)
}

func (c *cfImpl) V1() V1API {
	return &v1Impl{
		codefresh: c,
	}
}

func (c *cfImpl) V2() V2API {
	return &v2Impl{
		codefresh: c,
	}
}

// func (c *codefresh) Pipelines() IPipelineAPI {
// 	return newPipelineAPI(c)
// }

// func (c *codefresh) Users() UsersAPI {
// 	return newUsersAPI(c)
// }

// func (c *codefresh) UsersV2() IUsersV2API {
// 	return newUsersV2API(c)
// }

// func (c *codefresh) Tokens() ITokenAPI {
// 	return newTokenAPI(c)
// }

// func (c *codefresh) RuntimeEnvironments() IRuntimeEnvironmentAPI {
// 	return newRuntimeEnvironmentAPI(c)
// }

// func (c *codefresh) Workflows() IWorkflowAPI {
// 	return newWorkflowAPI(c)
// }

// func (c *codefresh) Progresses() IProgressAPI {
// 	return newProgressAPI(c)
// }

// func (c *codefresh) Clusters() IClusterAPI {
// 	return newClusterAPI(c)
// }

// func (c *codefresh) Contexts() IContextAPI {
// 	return newContextAPI(c)
// }

func (c *v1Impl) Argo() v1.ArgoAPI {
	return v1.newArgoAPI(c)
}

// func (c *codefresh) Gitops() GitopsAPI {
// 	return newGitopsAPI(c)
// }

// func (c *codefresh) V2() V2API {
// 	return c
// }

// func (c *codefresh) AccountV2() IAccountV2API {
// 	return newAccountV2API(c)
// }

// func (c *codefresh) Runtime() IRuntimeAPI {
// 	return newArgoRuntimeAPI(c)
// }

// func (c *codefresh) GitSource() IGitSourceAPI {
// 	return newGitSourceAPI(c)
// }

// func (c *codefresh) Component() IComponentAPI {
// 	return newComponentAPI(c)
// }

// func (c *codefresh) Workflow() IWorkflowV2API {
// 	return newWorkflowV2API(c)
// }

// func (c *codefresh) Pipeline() IPipelineV2API {
// 	return newPipelineV2API(c)
// }

// func (c *codefresh) CliReleases() ICliReleasesAPI {
// 	return newCliReleasesAPI(c)
// }

// func (c *codefresh) Cluster() IClusterV2API {
// 	return newClusterV2API(c)
// }

// func (c *codefresh) AppProxy(ctx context.Context, runtime string, insecure bool) (AppProxyAPI, error) {
// 	rt, err := c.V2().Runtime().Get(ctx, runtime)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create app-proxy client for runtime %s: %w", runtime, err)
// 	}

// 	var host string

// 	if rt.InternalIngressHost != nil && *rt.InternalIngressHost != "" {
// 		host = *rt.InternalIngressHost
// 	} else if rt.IngressHost != nil && *rt.IngressHost != "" {
// 		host = *rt.IngressHost
// 	} else {
// 		return nil, fmt.Errorf("failed to create app-proxy client for runtime %s: runtime does not have ingressHost configured", runtime)
// 	}

// 	httpClient := &http.Client{}
// 	httpClient.Timeout = c.client.Timeout
// 	if insecure {
// 		customTransport := http.DefaultTransport.(*http.Transport).Clone()
// 		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 		httpClient.Transport = customTransport
// 	}

// 	return newClient(&ClientOptions{
// 		Host:        host,
// 		Auth:        AuthOptions{Token: c.token},
// 		Client:      httpClient,
// 		graphqlPath: "/app-proxy/api/graphql",
// 	}), nil
// }

// func (c *codefresh) AppProxyClusters() IAppProxyClustersAPI {
// 	return newAppProxyClustersAPI(c)
// }

// func (c *codefresh) GitIntegrations() IAppProxyGitIntegrationsAPI {
// 	return newAppProxyGitIntegrationsAPI(c)
// }

// func (c *codefresh) VersionInfo() IAppProxyVersionInfoAPI {
// 	return newAppProxyVersionInfoAPI(c)
// }

// func (c *codefresh) AppProxyGitSources() IAppProxyGitSourcesAPI {
// 	return newAppProxyGitSourcesAPI(c)
// }

// func (c *codefresh) AppProxyIsc() IAppProxyIscAPI {
// 	return newAppProxyIscAPI(c)
// }

func (c *cfImpl) requestAPI(opt *requestOptions) (*http.Response, error) {
	return c.requestAPIWithContext(context.Background(), opt)
}

func (c *cfImpl) requestAPIWithContext(ctx context.Context, opt *requestOptions) (*http.Response, error) {
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

func (c *cfImpl) graphqlAPI(ctx context.Context, body map[string]interface{}, res interface{}) error {
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

func (c *cfImpl) decodeResponseInto(resp *http.Response, target interface{}) error {
	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *cfImpl) getBodyAsString(resp *http.Response) (string, error) {
	body, err := c.getBodyAsBytes(resp)
	return string(body), err
}

func (c *cfImpl) getBodyAsBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func newClient(opt *ClientOptions) *cfImpl {
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

	return &cfImpl{
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
