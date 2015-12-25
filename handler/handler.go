/*
package handler includes the logic that define a service as an object
that implements the Handler interface.

For pragmatic reasons we can construct a Handler from any struct using
reflection, where methods that satisfy these criteria will be made
available for remote access:

	- the method is exported.
	- the method returns two result, the second is of type error.

To be more concret for a method to be exposed, this later should look like:

	func (t *T) MethodName(arg1 T1, arg2 T2, ...) (replyType, error)

Here is a simple example of a valid service definition using a normal struct:

	type S struct {}

	func (S) Add(a, b int) (int, error) {
	return a + b, nil
	}

The same service implemented using a plain old handler will look:

	func AddHandler(rw protocol.ResponseWriter, req *protocol.Request) {
		a = req.Params["_0"].(int)
		b = req.Params["_1"].(int)
		rw.Set(a + b)
	}

*/
package handler

import "github.com/mouadino/go-nano/protocol"

// Handler interface for handling RPC requests.
type Handler interface {
	Handle(protocol.ResponseWriter, *protocol.Request)
}

// HandlerFunc represents a function that implement Handler interface.
type HandlerFunc func(protocol.ResponseWriter, *protocol.Request)

func (h HandlerFunc) Handle(rw protocol.ResponseWriter, req *protocol.Request) {
	h(rw, req)
}

// Middleware represents a function that wraps around a Handler and return
// another Handler.
type Middleware func(Handler) Handler

// New returns a handler from interface{} object.
// If object implements already Handler interface return as it's else apply
// reflection rules as described above.
func New(svc interface{}) Handler {
	if hdlr, ok := svc.(Handler); ok {
		return hdlr
	}
	return newStructHandler(svc)
}
