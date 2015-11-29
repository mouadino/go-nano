package loadbalancer

import (
	"math"
	"testing"

	"github.com/mouadino/go-nano/discovery"
)

func TestRandomLoadBalancer(t *testing.T) {
	lb := RandomLoadBalancer()
	instances := []discovery.Instance{
		discovery.Instance{Meta: discovery.ServiceMetadata{"endpoint": "1"}},
		discovery.Instance{Meta: discovery.ServiceMetadata{"endpoint": "2"}},
		discovery.Instance{Meta: discovery.ServiceMetadata{"endpoint": "3"}},
	}

	iterations := 100000
	want := iterations / 3
	tolerance := float64(want / 100)
	cnts := make(map[string]int)

	for i := 0; i < iterations; i++ {
		e, err := lb.Endpoint(instances)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		cnts[e]++
	}

	for _, v := range cnts {
		if math.Abs(float64(v-want)) > tolerance {
			t.Errorf("expected %d[±%f], got %d", want, tolerance, v)
		}
	}
}

func TestNegativeRandomLoadBalancer(t *testing.T) {
	lb := RandomLoadBalancer()
	instances := []discovery.Instance{}

	_, err := lb.Endpoint(instances)
	if err != NoEndpointError {
		t.Errorf("expected to fail with NoEndpointError, got %s ", err)
	}
}
