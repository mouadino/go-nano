package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestHTTPReceive(t *testing.T) {
	trans := New()
	trans.Listen()

	body := `{"hello": "world"}`
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/rpc/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	go trans.handle(w, req)

	select {
	case r := <-trans.Receive():
		b, ok := r.Body.([]byte)
		if !ok {
			t.Fatalf("request body is not []byte")
		}
		if string(b) != body {
			t.Errorf("request body doesn't match want %v, got %v", body, r.Body)
		}

		r.Resp.Write([]byte(body))

		if w.Body.String() != body {
			t.Errorf("response body didn't match, want %s got %s", body, w.Body)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Didn't receive any request after 1 second")
	}
}

func TestHTTPSend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello")
	}))
	defer ts.Close()

	trans := New()

	fmt.Printf(ts.URL)
	resp, err := trans.Send(ts.URL, strings.NewReader("foobar"))

	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if string(resp) != "Hello\n" {
		t.Errorf("unexpected response want %q, got %q", "Hello", resp)
	}
}

func TestHTTPAddr(t *testing.T) {
	trans := New()
	trans.Listen()

	addr, err := url.Parse(trans.Addr())

	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if addr.Scheme != "http" {
		t.Errorf("unexpected scheme want %q, got %q", "http", addr.Scheme)
	}
}

// TODO: Add benchmarks.
