package codefresh

import "github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"

type (
	Codefresh interface {
		V1() V1API
		V2() V2API
	}

	codefresh struct {
		client *client.CfClient
	}
)

func New(opt *client.ClientOptions) Codefresh {
	client := client.NewCfClient(opt)
	return &codefresh{ client: client }
}

func (c *codefresh) V1() V1API {
	return newV1Client(c.client)
}

func (c *codefresh) V2() V2API {
	return newV2Client(c.client)
}
