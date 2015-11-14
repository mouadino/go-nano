package discovery

import (
	"math/rand"
	"sync/atomic"

	"github.com/mouadino/go-nano/client"
	"github.com/mouadino/go-nano/protocol"
)

type randomLoadBalancer struct {
	rand *rand.Rand
}

func RandomLoadBalancer() *randomLoadBalancer {
	return &randomLoadBalancer{
		rand: rand.New(rand.NewSource(0)),
	}
}

func (lb *randomLoadBalancer) Endpoint(svc *Service) (Endpoint, error) {
	if len(svc.Instances) == 0 {
		return "", NoEndpointError
	}
	instance := svc.Instances[lb.rand.Intn(len(svc.Instances))]
	return instance.Meta.Endpoint(), nil
}

type roundRobinLoadBalancer struct {
	mod uint64
}

func RoundRobinLoadBalancer() *roundRobinLoadBalancer {
	return &roundRobinLoadBalancer{0}
}

func (s *roundRobinLoadBalancer) Endpoint(svc *Service) (Endpoint, error) {
	if len(svc.Instances) == 0 {
		return "", NoEndpointError
	}
	var old uint64
	for {
		old = atomic.LoadUint64(&s.mod)
		if atomic.CompareAndSwapUint64(&s.mod, old, old+1) {
			break
		}
	}
	instance := svc.Instances[s.mod%uint64(len(svc.Instances))]
	return instance.Meta.Endpoint(), nil
}

type loadBalanderExtension struct {
	client   client.Client
	lb       LoadBalancer
	resolver Resolver
}

func NewLoadBalancerExtension(resolver Resolver, lb LoadBalancer) client.ClientExtension {
	return func(c client.Client) client.Client {
		return &loadBalanderExtension{
			client:   c,
			lb:       lb,
			resolver: resolver,
		}
	}
}

func (lbExt *loadBalanderExtension) CallEndpoint(endpoint string, req *protocol.Request) (interface{}, error) {
	svc, err := lbExt.resolver.Resolve(endpoint)
	if err != nil {
		return nil, err
	}
	endpoint, err = lbExt.getEndpoint(svc)
	if err != nil {
		return nil, err
	}
	return lbExt.client.CallEndpoint(endpoint, req)
}

func (lbExt *loadBalanderExtension) getEndpoint(svc *Service) (string, error) {
	endpoint, err := lbExt.lb.Endpoint(svc)
	return string(endpoint), err
}
