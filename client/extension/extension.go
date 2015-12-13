package extension

import "github.com/mouadino/go-nano/protocol"

// Extension define a decorator to apply on client for customization purposes.
type Extension func(protocol.Sender) protocol.Sender

// Decorate wraps a sender using the given client extensions.
func Decorate(rs protocol.Sender, exts ...Extension) protocol.Sender {
	for _, ext := range exts {
		rs = ext(rs)
	}
	return rs
}
