package utils

import (
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
)

func CastToCodefreshOrDie(candidate interface{}) codefresh.Codefresh {
	client, ok := candidate.(codefresh.Codefresh)
	if !ok {
		panic("Failed to cast candidate to Codefresh client")
	}
	return client
}
