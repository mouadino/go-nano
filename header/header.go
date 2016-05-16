package header

import "net/textproto"

// Header represents the key-value pairs in an RPC request/response header.
type Header map[string][]string

// Get gets the first value associated with the given key.
func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

// Set sets the header entries associated with key to the single element value.
func (h Header) Set(key string, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}
