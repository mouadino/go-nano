package jsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func equalJSON(first, second []byte) error {
	f, err := toJSON(first)
	if err != nil {
		return err
	}
	s, err := toJSON(second)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(f, s) {
		return fmt.Errorf("Expected %s, got %s", first, second)
	}
	return nil
}

func toJSON(b []byte) (interface{}, error) {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: %s", err)
	}
	return f, err
}
