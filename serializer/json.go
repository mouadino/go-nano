package serializer

import "encoding/json"

type JSONSerializer struct{}

func (JSONSerializer) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (JSONSerializer) Decode(data []byte, result interface{}) error {
	return json.Unmarshal(data, result)
}
