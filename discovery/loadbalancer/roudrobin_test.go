package loadbalancer

import (
	"reflect"
	"testing"

	"github.com/mouadino/go-nano/discovery"
)

func TestRoundRobinLoadBalancer(t *testing.T) {
	lb := RoundRobinLoadBalancer()
	instances := []discovery.Instance{
		discovery.Instance{Meta: discovery.ServiceMetadata{"endpoint": "1"}},
		discovery.Instance{Meta: discovery.ServiceMetadata{"endpoint": "2"}},
		discovery.Instance{Meta: discovery.ServiceMetadata{"endpoint": "3"}},
	}

	iterations := 100
	want := []string{"1", "2", "3"}
	for i := 0; i < iterations; i++ {
		endpoint, _ := lb.Endpoint(instances)

		if endpoint != want[i%3] {
			t.Errorf("round robin fail at %s: excpected %s got %s", i, want[i%3], endpoint)
		}
	}
}

func TestConcurrentRoundRobinLoadBalancer(t *testing.T) {
	lb := RoundRobinLoadBalancer()
	instances := []discovery.Instance{
		discovery.Instance{Meta: discovery.ServiceMetadata{"endpoint": "1"}},
		discovery.Instance{Meta: discovery.ServiceMetadata{"endpoint": "2"}},
		discovery.Instance{Meta: discovery.ServiceMetadata{"endpoint": "3"}},
	}

	iterations := 9
	endpoints := make(chan string)
	goroutine := func() {
		for i := 0; i < iterations; i++ {
			e, err := lb.Endpoint(instances)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			endpoints <- e
		}
	}

	go goroutine()
	go goroutine()
	go goroutine()

	cnts := make(map[string]int)
	for j := 0; j < iterations*3; j++ {
		cnts[<-endpoints]++
	}

	want := map[string]int{
		"1": iterations,
		"2": iterations,
		"3": iterations,
	}
	if !reflect.DeepEqual(cnts, want) {
		t.Errorf("round robin is not uniform %v ", cnts)
	}
}

func TestNegativeRoundRobinLoadBalancer(t *testing.T) {
	lb := RoundRobinLoadBalancer()
	instances := []discovery.Instance{}

	_, err := lb.Endpoint(instances)
	if err != NoEndpointError {
		t.Errorf("expected to fail with NoEndpointError, got %s ", err)
	}
}
