package utils

import (
	"net"
	"reflect"
	"testing"

	"github.com/mouadino/go-nano/protocol"
)

type Slice []interface{}

func checkIP(t *testing.T, addr string) {
	ip := net.ParseIP(addr)
	if ip.IsLoopback() {
		t.Errorf("ip is looback")
	}

	if ip.IsUnspecified() {
		t.Errorf("ip is unspecified")
	}
}

func TestParamsFormat(t *testing.T) {
	paramsTests := []struct {
		in  Slice
		out protocol.Params
	}{
		{Slice{"foobar"}, protocol.Params{"_0": "foobar"}},
		{Slice{"foobar", 42}, protocol.Params{"_0": "foobar", "_1": 42}},
	}

	for _, test := range paramsTests {
		out := ParamsFormat(test.in...)
		if !reflect.DeepEqual(test.out, out) {
			t.Errorf("ParamsFormat(%s) => %q, want %q", test.in, out, test.out)
		}
	}
}

func TestGetExternalIP(t *testing.T) {
	addr, err := GetExternalIP()

	if err != nil {
		t.Fatalf("unexpected failure %s", err)
	}

	checkIP(t, addr)
}

func TestListener(t *testing.T) {
	ln, err := GetListener("10.0.0.1")

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ip := ln.Addr().String()

	if ip != "10.0.0.1" {
		t.Errorf("IP doesn't match want 10.0.0.1, got %q", ip)
	}
}

func TestListenerEmptyAddress(t *testing.T) {
	ln, err := GetListener("")

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ip := ln.Addr().String()

	checkIP(t, ip)
}
