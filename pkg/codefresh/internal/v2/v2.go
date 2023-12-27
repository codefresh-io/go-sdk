package v2

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/ap"
	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	V2API interface {
		Account() AccountAPI
		AppProxy(ctx context.Context, runtime string, insecure bool) (ap.AppProxyAPI, error)
		CliRelease() CliReleaseAPI
		Cluster() ClusterAPI
		Component() ComponentAPI
		GitSource() GitSourceAPI
		Pipeline() PipelineAPI
		Runtime() RuntimeAPI
		User() UserAPI
		Workflow() WorkflowAPI
	}

	v2Impl struct {
		client *client.CfClient
	}
)

func NewV2Client(c *client.CfClient) V2API {
	return &v2Impl{client: c}
}

func (v2 *v2Impl) Account() AccountAPI {
	return &account{client: v2.client}
}

func (v2 *v2Impl) AppProxy(ctx context.Context, runtime string, insecure bool) (ap.AppProxyAPI, error) {
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
	httpClient.Timeout = v2.client.Timeout()
	if insecure {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		httpClient.Transport = customTransport
	}

	c := client.NewCfClient(&client.ClientOptions{
		Token:       v2.client.Token(),
		Host:        host,
		Client:      httpClient,
		GraphqlPath: "/app-proxy/api/graphql",
	})
	return ap.NewAppProxyClient(c), nil
}

func (v2 *v2Impl) CliRelease() CliReleaseAPI {
	return &cliRelease{client: v2.client}
}

func (v2 *v2Impl) Cluster() ClusterAPI {
	return &cluster{client: v2.client}
}

func (v2 *v2Impl) Component() ComponentAPI {
	return &component{client: v2.client}
}

func (v2 *v2Impl) GitSource() GitSourceAPI {
	return &gitSource{client: v2.client}
}

func (v2 *v2Impl) Pipeline() PipelineAPI {
	return &pipeline{client: v2.client}
}

func (v2 *v2Impl) Runtime() RuntimeAPI {
	return &runtime{client: v2.client}
}

func (v2 *v2Impl) User() UserAPI {
	return &user{client: v2.client}
}

func (v2 *v2Impl) Workflow() WorkflowAPI {
	return &workflow{client: v2.client}
}
