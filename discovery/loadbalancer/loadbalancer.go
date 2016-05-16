/*
package loadbalancer defines how to balance and distribute requests between multiple
service instances.

*/
package loadbalancer

import (
	"errors"

	"github.com/mouadino/go-nano/client/extension"
	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
)

// TODO: Move this outside discovery.

// NoEndpointError returned when loadbalancer doesn't find any suitable endpoint.
var NoEndpointError = errors.New("No Endpoint")

// LoadBalancerStrategy is the interface that define how an endpoint get chosen from
// available endpoints.
type LoadBalancerStrategy interface {
	// TODO: Pass request too ?
	Endpoint([]discovery.Instance) (string, error)
}

type loadBalanderExtension struct {
	sender   transport.Sender
	strategy LoadBalancerStrategy
	resolver discovery.Resolver
}

// New returns a client extension that know how to balance requests.
func New(resolver discovery.Resolver, strategy LoadBalancerStrategy) extension.Extension {
	return func(s transport.Sender) transport.Sender {
		return &loadBalanderExtension{
			sender:   s,
			strategy: strategy,
			resolver: resolver,
		}
	}
}

// Send function to implement transport.Sender interface.
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
	endpoint, err := lb.strategy.Endpoint(instances)
	return endpoint, err
}
