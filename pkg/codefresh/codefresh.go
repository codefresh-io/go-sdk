package codefresh

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/go-querystring/query"
)

//go:generate mockery -name Codefresh -filename codefresh.go

//go:generate mockery -name UsersAPI -filename users.go

type (
	Codefresh interface {
		Pipelines() IPipelineAPI
		Tokens() ITokenAPI
		RuntimeEnvironments() IRuntimeEnvironmentAPI
		Workflows() IWorkflowAPI
		Progresses() IProgressAPI
		Clusters() IClusterAPI
		Contexts() IContextAPI
		Users() UsersAPI
		Argo() ArgoAPI
		Gitops() GitopsAPI
		V2() V2API
		AppProxy(ctx context.Context, runtime string, insecure bool) (AppProxyAPI, error)
	}

	V2API interface {
		UsersV2() IUsersV2API
		Runtime() IRuntimeAPI
		Cluster() IClusterV2API
		GitSource() IGitSourceAPI
		Component() IComponentAPI
		Workflow() IWorkflowV2API
		Pipeline() IPipelineV2API
		CliReleases() ICliReleasesAPI
	}

	AppProxyAPI interface {
		AppProxyClusters() IAppProxyClustersAPI
		GitIntegrations() IAppProxyGitIntegrationsAPI
		VersionInfo() IAppProxyVersionInfoAPI
		AppProxyGitSources() IAppProxyGitSourcesAPI
		AppProxyIsc() IAppProxyIscAPI
	}
)

func New(opt *ClientOptions) Codefresh {
	return newClient(opt)
}

func (c *codefresh) Pipelines() IPipelineAPI {
	return newPipelineAPI(c)
}

func (c *codefresh) Users() UsersAPI {
	return newUsersAPI(c)
}

func (c *codefresh) UsersV2() IUsersV2API {
	return newUsersV2API(c)
}

func (c *codefresh) Tokens() ITokenAPI {
	return newTokenAPI(c)
}

func (c *codefresh) RuntimeEnvironments() IRuntimeEnvironmentAPI {
	return newRuntimeEnvironmentAPI(c)
}

func (c *codefresh) Workflows() IWorkflowAPI {
	return newWorkflowAPI(c)
}

func (c *codefresh) Progresses() IProgressAPI {
	return newProgressAPI(c)
}

func (c *codefresh) Clusters() IClusterAPI {
	return newClusterAPI(c)
}

func (c *codefresh) Contexts() IContextAPI {
	return newContextAPI(c)
}

func (c *codefresh) Argo() ArgoAPI {
	return newArgoAPI(c)
}

func (c *codefresh) Gitops() GitopsAPI {
	return newGitopsAPI(c)
}

func (c *codefresh) V2() V2API {
	return c
}

func (c *codefresh) Runtime() IRuntimeAPI {
	return newArgoRuntimeAPI(c)
}

func (c *codefresh) GitSource() IGitSourceAPI {
	return newGitSourceAPI(c)
}

func (c *codefresh) Component() IComponentAPI {
	return newComponentAPI(c)
}

func (c *codefresh) Workflow() IWorkflowV2API {
	return newWorkflowV2API(c)
}

func (c *codefresh) Pipeline() IPipelineV2API {
	return newPipelineV2API(c)
}

func (c *codefresh) CliReleases() ICliReleasesAPI {
	return newCliReleasesAPI(c)
}

func (c *codefresh) Cluster() IClusterV2API {
	return newClusterV2API(c)
}

func (c *codefresh) AppProxy(ctx context.Context, runtime string, insecure bool) (AppProxyAPI, error) {
	rt, err := c.V2().Runtime().Get(ctx, runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to create app-proxy client for runtime %s: %w", runtime, err)
	}

	var host string

	if rt.InternalIngressHost != nil && *rt.InternalIngressHost != "" {
		host = *rt.InternalIngressHost
	} else if rt.IngressHost != nil && *rt.IngressHost != "" {
		host = *rt.IngressHost
	} else {
		return nil, fmt.Errorf("failed to create app-proxy client for runtime %s: runtime does not have ingressHost configured", runtime)
	}

	httpClient := &http.Client{}
	httpClient.Timeout = c.client.Timeout
	if insecure {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		httpClient.Transport = customTransport
	}

	return newClient(&ClientOptions{
		Host:        host,
		Auth:        AuthOptions{Token: c.token},
		Client:      httpClient,
		graphqlPath: "/app-proxy/api/graphql",
	}), nil
}

func (c *codefresh) AppProxyClusters() IAppProxyClustersAPI {
	return newAppProxyClustersAPI(c)
}

func (c *codefresh) GitIntegrations() IAppProxyGitIntegrationsAPI {
	return newAppProxyGitIntegrationsAPI(c)
}

func (c *codefresh) VersionInfo() IAppProxyVersionInfoAPI {
	return newAppProxyVersionInfoAPI(c)
}

func (c *codefresh) AppProxyGitSources() IAppProxyGitSourcesAPI {
	return newAppProxyGitSourcesAPI(c)
}

func (c *codefresh) AppProxyIsc() IAppProxyIscAPI {
	return newAppProxyIscAPI(c)
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

	data, err := ioutil.ReadAll(response.Body)
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
	body, err := ioutil.ReadAll(resp.Body)
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
		token:       opt.Auth.Token,
		graphqlPath: graphqlPath,
		client:      httpClient,
	}
}
