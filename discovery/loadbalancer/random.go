package loadbalancer

import (
	"math/rand"

	"github.com/mouadino/go-nano/discovery"
)

type randomLoadBalancer struct {
	rand *rand.Rand
}

func RandomLoadBalancer() *randomLoadBalancer {
	return &randomLoadBalancer{
		rand: rand.New(rand.NewSource(0)),
	}
}

func (lb *randomLoadBalancer) Endpoint(instances []discovery.Instance) (string, error) {
	if len(instances) == 0 {
		return "", NoEndpointError
	}
	instance := instances[lb.rand.Intn(len(instances))]
	return instance.Meta.Endpoint(), nil
}
