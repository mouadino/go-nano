package nano

// TODO: Use me !
type Configurable interface {
	NanoConfigure(interface{}) error
}

type Startable interface {
	NanoStart() error
	NanoStop()
}

// TODO: Should we use io.ReadCloser !?
type Response struct {
	Body []byte
}
