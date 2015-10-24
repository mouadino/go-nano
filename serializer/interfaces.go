package serializer

type Serializer interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}
