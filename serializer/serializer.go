package serializer

import "io"

type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

type Decoder interface {
	Decode(io.Reader, interface{}) error
}

type Serializer interface {
	Encoder
	Decoder
}
