package discovery

import "testing"

func TestService(t *testing.T) {

	srv := Service{
		Name: "foobar",
		Instances: []Instance{
			NewInstance("127.0.0.1:8080", InstanceMeta{}),
			NewInstance("127.0.0.1:8081", InstanceMeta{}),
		},
	}

	if srv.String() != "foobar [2]" {
		t.Errorf("service name want 'foobar [2]', got %s", srv.String())
	}

	endpoint := srv.Instances[0].Endpoint
	if endpoint != "127.0.0.1:8080" {
		t.Errorf("service endpoint want '127.0.0.1:8080', got %s", endpoint)
	}
}
