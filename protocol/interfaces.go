package protocol

type Protocol interface {
	SendRequest(string, *Request) (interface{}, error)
	ReceiveRequest() (ResponseWriter, *Request)
}

type Params map[string]interface{}

type Request struct {
	Method string
	Params Params
	// TODO: Headers header.Header
}

type ResponseWriter interface {
	// TODO: Header() header.Header

	Write(interface{}) error
	WriteError(err error) error
}
