package header

import "testing"

func TestHeader(t *testing.T) {
	var (
		h = Header{}
		v string
	)

	v = h.Get("key1")
	if v != "" {
		t.Errorf("h.Get('key1') wants '', got %q", v)
	}

	h.Set("key1", "value1")
	v = h.Get("key1")
	if v != "value1" {
		t.Errorf("h.Get('key1') wants 'value1', got %q", v)
	}

	v = h.Get("unknown_key")
	if v != "" {
		t.Errorf("h.Get('unknown_key') wants '', got %q", v)
	}

	h.Set("key1", "value2")
	if v != "value2" {
		t.Errorf("h.Get('key1') wants 'value2', got %q", v)
	}
}
