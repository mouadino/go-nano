package loadbalancer

import (
	"testing"

	"github.com/mouadino/go-nano/discovery"
)

func TestRoundRobinLoadBalancer(t *testing.T) {
	lb := NewRoundRobin()
	instances := []discovery.Instance{
		discovery.Instance{Endpoint: "1"},
		discovery.Instance{Endpoint: "2"},
		discovery.Instance{Endpoint: "3"},
	}

	iterations := 100
	want := []string{"1", "2", "3"}
	for i := 0; i < iterations; i++ {
		endpoint, _ := lb.Endpoint(instances)

		if endpoint != want[i%3] {
			t.Errorf("want %s got %s (iteration %d)", want[i%3], endpoint, i)
		}
	}
}

func TestNegativeRoundRobinLoadBalancer(t *testing.T) {
	lb := NewRoundRobin()
	instances := []discovery.Instance{}

	_, err := lb.Endpoint(instances)
	if err != NoEndpointError {
		t.Errorf("expected to fail with NoEndpointError, got %s ", err)
	}
}
