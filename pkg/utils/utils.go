package utils

import (
	"encoding/json"
	"fmt"

	"github.com/codefresh-io/go-sdk/internal"
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
)

func CastToCodefreshOrDie(candidate interface{}) codefresh.Codefresh {
	client, ok := candidate.(codefresh.Codefresh)
	if !ok {
		internal.DieOnError(fmt.Errorf("Failed to cast candidate to Codefresh client"))
	}
	return client
}

func Convert(from interface{}, to interface{}) {
	rs, _ := json.Marshal(from)
	_ = json.Unmarshal(rs, to)
}
