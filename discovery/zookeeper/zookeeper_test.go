// +build integration

package zookeeper

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mouadino/go-nano/discovery"
	"github.com/samuel/go-zookeeper/zk"
)

func waitFor(timeout time.Duration, predicate func() bool) bool {
	ch := make(chan struct{}, 1)
	go func() {
		for {
			if predicate() {
				ch <- struct{}{}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-time.After(timeout):
		return false
	case <-ch:
		return true
	}
}

func TestZookeeperAnnouncing(t *testing.T) {
	ts, err := zk.StartTestCluster(3, nil, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Stop()

	hosts := make([]string, len(ts.Servers))
	for i, srv := range ts.Servers {
		hosts[i] = fmt.Sprintf("127.0.0.1:%d", srv.Port)
	}

	an := New(hosts)

	inst := discovery.NewInstance("<dummy>", discovery.InstanceMeta{})
	an.Announce("nano_test", inst)

	srv, err := an.Resolve("nano_test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if srv.Name != "nano_test" {
		t.Errorf("service name want 'nano_test', got %s", srv.Name)
	}

	if len(srv.Instances) != 1 {
		t.Errorf("len(instances) want 1 instance, got %d", len(srv.Instances))
	}

	if srv.Instances[0].Endpoint != "<dummy>" {
		t.Errorf("first instance endpoint want '<dummy>', got %s", srv.Instances[0].Endpoint)
	}
}

func TestZookeeperResolving(t *testing.T) {
	ts, err := zk.StartTestCluster(3, nil, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Stop()

	hosts := make([]string, len(ts.Servers))
	for i, srv := range ts.Servers {
		hosts[i] = fmt.Sprintf("127.0.0.1:%d", srv.Port)
	}

	an := New(hosts)

	an.Announce("nano_test", discovery.NewInstance("<dummy>", discovery.InstanceMeta{}))

	srv, err := an.Resolve("nano_test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(srv.Instances) != 1 {
		t.Errorf("len(instances) want 1 instance, got %d", len(srv.Instances))
	}

	an.Announce("nano_test", discovery.NewInstance("<dummy2>", discovery.InstanceMeta{}))

	waitFor(3*time.Second, func() bool { return len(srv.Instances) == 2 })

	if len(srv.Instances) != 2 {
		t.Errorf("len(instances) want 2 instance, got %d", len(srv.Instances))
	}
}
