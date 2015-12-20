package loadbalancer

import (
	"testing"

	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/dummy"
)

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

	sender := lb(dummy.New())
	req := &protocol.Request{
		Method: "foo",
	}
	_, err := sender.Send("nanoTest", req)
	if err != nil {
		t.Fatalf("Unexcpected failure %s", err)
	}
}
