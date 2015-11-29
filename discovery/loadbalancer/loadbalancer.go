package loadbalancer

import (
	"errors"

	"github.com/mouadino/go-nano/client/extension"
	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/protocol"
)

var NoEndpointError = errors.New("No Endpoint")

type LoadBalancer interface {
	Endpoint([]discovery.Instance) (string, error)
}

type loadBalanderExtension struct {
	sender   protocol.Sender
	lb       LoadBalancer
	resolver discovery.Resolver
}

func NewLoadBalancerExtension(resolver discovery.Resolver, lb LoadBalancer) extension.Extension {
	return func(s protocol.Sender) protocol.Sender {
		return &loadBalanderExtension{
			sender:   s,
			lb:       lb,
			resolver: resolver,
		}
	}
}

func (lb *loadBalanderExtension) Send(endpoint string, req *protocol.Request) (*protocol.Response, error) {
	svc, err := lb.resolver.Resolve(endpoint)
	if err != nil {
		return nil, err
	}
	endpoint, err = lb.getEndpoint(svc.Instances)
	if err != nil {
		return nil, err
	}
	return lb.sender.Send(endpoint, req)
}

func (lb *loadBalanderExtension) getEndpoint(instances []discovery.Instance) (string, error) {
	endpoint, err := lb.lb.Endpoint(instances)
	return endpoint, err
}
