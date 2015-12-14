package discovery

import (
	"fmt"
	"testing"
)

func TestService(t *testing.T) {
	insts := make([]Instance, 0)
	insts = append(insts, NewInstance("127.0.0.1:8080", InstanceMeta{}))
	insts = append(insts, NewInstance("127.0.0.1:8081", InstanceMeta{}))

	srv := Service{
		Name:      "foobar",
		Instances: insts,
	}

	fmt.Printf("%s - %s", srv.Instances, insts)

	if srv.String() != "foobar [2]" {
		t.Errorf("service name want 'foobar [2]', got %s", srv.String())
	}

	endpoint := srv.Instances[0].Endpoint
	if endpoint != "127.0.0.1:8080" {
		t.Errorf("service endpoint want '127.0.0.1:8080', got %s", endpoint)
	}
}
