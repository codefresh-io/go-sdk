package codefresh

import (
	"context"
	"fmt"
	"net/http"

	ap "github.com/codefresh-io/go-sdk/pkg/appproxy"
	"github.com/codefresh-io/go-sdk/pkg/client"
	gql "github.com/codefresh-io/go-sdk/pkg/graphql"
	"github.com/codefresh-io/go-sdk/pkg/rest"
)

type (
	Codefresh interface {
		AppProxy(ctx context.Context, runtime string, insecure bool) (ap.AppProxyAPI, error)
		GraphQL() gql.GraphQLAPI
		Rest() rest.RestAPI
		InternalClient() *client.CfClient
		HttpClient() HttpClient
	}

	HttpClient interface {
		NativeRestAPI(ctx context.Context, opt *client.RequestOptions) (*http.Response, error)
	}

	ClientOptions struct {
		Token       string
		Host        string
		Client      *http.Client
		GraphqlPath string
	}

	codefresh struct {
		client *client.CfClient
	}
)

func New(opt *ClientOptions) Codefresh {
	client := client.NewCfClient(opt.Host, opt.Token, opt.GraphqlPath, opt.Client)
	return &codefresh{client: client}
}

func (c *codefresh) AppProxy(ctx context.Context, runtime string, insecure bool) (ap.AppProxyAPI, error) {
	rt, err := c.GraphQL().Runtime().Get(ctx, runtime)
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

	apClient := c.client.AppProxyClient(host, insecure)
	return ap.NewAppProxyClient(apClient), nil
}

func (c *codefresh) GraphQL() gql.GraphQLAPI {
	return gql.NewGraphQLClient(c.client)
}

func (c *codefresh) Rest() rest.RestAPI {
	return rest.NewRestClient(c.client)
}

func (c *codefresh) InternalClient() *client.CfClient {
	return c.client
}

func (c *codefresh) HttpClient() HttpClient {
	return c.client
}
