package interfaces

type Transport interface {
	Receive() <-chan Data
	Send(endpoint string, data []byte) (ResponseReader, error)
}

// TODO: Rename me
type Data struct {
	Body []byte
	Resp ResponseWriter
}

type Protocol interface {
	SendRequest(string, *Request) (ResponseReader, error)
	ReceiveRequest() (ResponseWriter, *Request)
	// TODO: SendError !?
}

type Handler interface {
	Handle(ResponseWriter, *Request) error
}

type Middleware func(Handler) Handler

type Header map[string][]string

type ResponseWriter interface {
	Header() Header

	Write(interface{}) error
}

type ResponseReader interface {
	Read() ([]byte, error)
}

type Serializer interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

// TODO: Use me !
type Configurable interface {
	Configure(interface{}) error
}

// TODO: Use me !
type Startable interface {
	Start() error
	Stop() error
}

// TODO: Is this specific to JSON-RPC !?
type Request struct {
	Method string
	Params map[string]interface{}
	// TODO: Headers header
}

// TODO: Should we use io.ReadCloser !?
type Response struct {
	Body []byte
}

// TODO: Use this !?
type Server interface {
	ListenAndServe() error
}
