package serializer

// TODO: How can a developer use go-nano w/o having to install unwanted deps ?
import (
	"github.com/golang/protobuf/proto"
)

type ProtobufSerializer struct{}

func (ProtobufSerializer) Encode(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (ProtobufSerializer) Decode(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
