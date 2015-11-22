package loadbalancer

import (
	"errors"

	"github.com/mouadino/go-nano/client"
	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/protocol"
)

var NoEndpointError = errors.New("No Endpoint")

type LoadBalancer interface {
	Endpoint([]discovery.Instance) (string, error)
}

type loadBalanderExtension struct {
	client   client.Client
	lb       LoadBalancer
	resolver discovery.Resolver
}

func NewLoadBalancerExtension(resolver discovery.Resolver, lb LoadBalancer) client.ClientExtension {
	return func(c client.Client) client.Client {
		return &loadBalanderExtension{
			client:   c,
			lb:       lb,
			resolver: resolver,
		}
	}
}

func (lb *loadBalanderExtension) CallEndpoint(endpoint string, req *protocol.Request) (interface{}, error) {
	svc, err := lb.resolver.Resolve(endpoint)
	if err != nil {
		return nil, err
	}
	endpoint, err = lb.getEndpoint(svc.Instances)
	if err != nil {
		return nil, err
	}
	return lb.client.CallEndpoint(endpoint, req)
}

func (lb *loadBalanderExtension) getEndpoint(instances []discovery.Instance) (string, error) {
	endpoint, err := lb.lb.Endpoint(instances)
	return endpoint, err
}
