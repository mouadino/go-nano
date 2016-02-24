package serializer

import (
	"encoding/json"
	"io"
)

type JSONSerializer struct{}

func (JSONSerializer) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (JSONSerializer) Decode(r io.Reader, result interface{}) error {
	dec := json.NewDecoder(r)
	return dec.Decode(result)
}
