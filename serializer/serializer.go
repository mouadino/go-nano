package serializer

type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

type Decoder interface {
	Decode([]byte, interface{}) error
}

type Serializer interface {
	Encoder
	Decoder
}
