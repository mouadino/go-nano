package extension

import "github.com/mouadino/go-nano/transport"

// Extension define a decorator to apply on client for customization purposes.
type Extension func(transport.Sender) transport.Sender

// Decorate wraps a sender using the given client extensions.
func Decorate(next transport.Sender, exts ...Extension) transport.Sender {
	for _, ext := range exts {
		next = ext(next)
	}
	return next
}
