package client

// Future represents the response of an asynchronous client request.
type Future struct {
	resp interface{}
	err  error
	done chan struct{}
}

func newFuture() *Future {
	return &Future{
		done: make(chan struct{}, 1),
	}
}

// Result returns the future result, block until request finish.
func (f *Future) Result() (interface{}, error) {
	<-f.done
	return f.resp, f.err
}

func (f *Future) set(resp interface{}, err error) {
	f.resp = resp
	f.err = err
	f.done <- struct{}{}
}
