package loadbalancer

import (
	"testing"

	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/protocol"
)

type dummySender struct{}

func (dummySender) Send(e string, req *protocol.Request) (*protocol.Response, error) {
	resp := &protocol.Response{
		Body: "",
	}
	return resp, nil
}

type firstStrategy struct{}

func (firstStrategy) Endpoint(insts []discovery.Instance) (string, error) {
	return insts[0].Endpoint, nil
}

type staticResolver struct{}

func (staticResolver) Resolve(name string) (*discovery.Service, error) {
	insts := []discovery.Instance{
		discovery.NewInstance("memory:///", nil),
	}
	svc := &discovery.Service{
		Name:      name,
		Instances: insts,
	}

	return svc, nil
}

func TestLoadBalancer(t *testing.T) {
	lb := New(staticResolver{}, firstStrategy{})
	req := &protocol.Request{
		Method: "foo",
	}

	sender := lb(dummySender{})

	_, err := sender.Send(":dummy:", req)

	if err != nil {
		t.Fatalf("Unexcpected failure %s", err)
	}

	// TODO: Add more assertion.
}
