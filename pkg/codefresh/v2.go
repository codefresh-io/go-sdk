package codefresh

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
)

type (
	V2API interface {
		AccountV2() IAccountV2API
		AppProxy(ctx context.Context, runtime string, insecure bool) (AppProxyAPI, error)
		CliReleases() ICliReleasesAPI
		Cluster() IClusterV2API
		Component() IComponentAPI
		GitSource() IGitSourceAPI
		Pipeline() IPipelineV2API
		Runtime() IRuntimeAPI
		UsersV2() IUsersV2API
		Workflow() IWorkflowV2API
	}

	v2Impl struct {
		codefresh *codefresh
	}
)

func newV2Client(c *codefresh) V2API {
	return &v2Impl{codefresh: c}
}

func (v2 *v2Impl) AccountV2() IAccountV2API {
	return newAccountV2API(v2.codefresh)
}

func (v2 *v2Impl) AppProxy(ctx context.Context, runtime string, insecure bool) (AppProxyAPI, error) {
	rt, err := v2.Runtime().Get(ctx, runtime)
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
	httpClient.Timeout = v2.codefresh.client.Timeout
	if insecure {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		httpClient.Transport = customTransport
	}

	c := newClient(&ClientOptions{
		Host:        host,
		Token:       v2.codefresh.token,
		Client:      httpClient,
		graphqlPath: "/app-proxy/api/graphql",
	})
	return newAppProxyClient(c), nil
}

func (v2 *v2Impl) CliReleases() ICliReleasesAPI {
	return newCliReleasesAPI(v2.codefresh)
}

func (v2 *v2Impl) Cluster() IClusterV2API {
	return newClusterV2API(v2.codefresh)
}

func (v2 *v2Impl) Component() IComponentAPI {
	return newComponentAPI(v2.codefresh)
}

func (v2 *v2Impl) GitSource() IGitSourceAPI {
	return newGitSourceAPI(v2.codefresh)
}

func (v2 *v2Impl) Pipeline() IPipelineV2API {
	return newPipelineV2API(v2.codefresh)
}

func (v2 *v2Impl) Runtime() IRuntimeAPI {
	return newArgoRuntimeAPI(v2.codefresh)
}

func (v2 *v2Impl) UsersV2() IUsersV2API {
	return newUsersV2API(v2.codefresh)
}

func (v2 *v2Impl) Workflow() IWorkflowV2API {
	return newWorkflowV2API(v2.codefresh)
}
