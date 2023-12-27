package codefresh

import (
	"github.com/codefresh-io/go-sdk/pkg/codefresh/internal/client"
	v1 "github.com/codefresh-io/go-sdk/pkg/codefresh/internal/v1"
	v2 "github.com/codefresh-io/go-sdk/pkg/codefresh/internal/v2"
)

type (
	Codefresh interface {
		V1() v1.V1API
		V2() v2.V2API
	}

	codefresh struct {
		client *client.CfClient
	}
)

func New(opt *client.ClientOptions) Codefresh {
	client := client.NewCfClient(opt)
	return &codefresh{client: client}
}

func (c *codefresh) V1() v1.V1API {
	return v1.NewV1Client(c.client)
}

func (c *codefresh) V2() v2.V2API {
	return v2.NewV2Client(c.client)
}
