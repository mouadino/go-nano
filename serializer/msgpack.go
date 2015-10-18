package serializer

import (
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

type MsgPackSerializer struct{}

func (MsgPackSerializer) Encode(data interface{}) ([]byte, error) {
	return msgpack.Marshal(data)
}

func (MsgPackSerializer) Decode(data []byte, result interface{}) error {
	return msgpack.Unmarshal(data, result)
}
