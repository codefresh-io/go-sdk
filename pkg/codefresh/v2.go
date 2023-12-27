package codefresh

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
)

type (
	V2API interface {
		Account() V2AccountAPI
		AppProxy(ctx context.Context, runtime string, insecure bool) (AppProxyAPI, error)
		CliRelease() V2CliReleaseAPI
		Cluster() V2ClusterAPI
		Component() V2ComponentAPI
		GitSource() V2GitSourceAPI
		Pipeline() V2PipelineAPI
		Runtime() V2RuntimeAPI
		User() V2UserAPI
		Workflow() V2WorkflowAPI
	}

	v2Impl struct {
		client *client.CfClient
	}
)

func newV2Client(c *client.CfClient) V2API {
	return &v2Impl{client: c}
}

func (v2 *v2Impl) Account() V2AccountAPI {
	return &v2Account{client: v2.client}
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
	return newAppProxyClient(c), nil
}

func (v2 *v2Impl) CliRelease() V2CliReleaseAPI {
	return &v2CliRelease{client: v2.client}
}

func (v2 *v2Impl) Cluster() V2ClusterAPI {
	return &v2Cluster{client: v2.client}
}

func (v2 *v2Impl) Component() V2ComponentAPI {
	return &v2Component{client: v2.client}
}

func (v2 *v2Impl) GitSource() V2GitSourceAPI {
	return &v2GitSource{client: v2.client}
}

func (v2 *v2Impl) Pipeline() V2PipelineAPI {
	return &v2Pipeline{client: v2.client}
}

func (v2 *v2Impl) Runtime() V2RuntimeAPI {
	return &v2Runtime{client: v2.client}
}

func (v2 *v2Impl) User() V2UserAPI {
	return &v2User{client: v2.client}
}

func (v2 *v2Impl) Workflow() V2WorkflowAPI {
	return &v2Workflow{client: v2.client}
}
