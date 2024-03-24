package model

import (
	"encoding/json"
	"fmt"
	"io"
)

type StringMap map[string]string

func (m *StringMap) UnmarshalGQL(v interface{}) error {
	anyMap, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("StringMap must be a map")
	}

	*m = make(map[string]string, len(anyMap))
	for k, v := range anyMap {
		(*m)[k], ok = v.(string)
		if !ok {
			return fmt.Errorf("StringMap value %q must be strings", k)
		}
	}

	return nil
}

func (m StringMap) MarshalGQL(w io.Writer) {
	_ = json.NewEncoder(w).Encode(m)
}
