package utils

import (
	"reflect"
	"testing"

	"github.com/mouadino/go-nano/protocol"
)

type Slice []interface{}

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
