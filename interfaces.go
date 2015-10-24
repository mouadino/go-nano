package nano

// TODO: Use me !
type Configurable interface {
	Configure(interface{}) error
}

// TODO: Use me !
type Startable interface {
	Start() error
	Stop() error
}

// TODO: Should we use io.ReadCloser !?
type Response struct {
	Body []byte
}

// TODO: Use this !?
type Server interface {
	ListenAndServe() error
}
