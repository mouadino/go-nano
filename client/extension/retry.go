package extension

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/cenkalti/backoff"
	"github.com/mouadino/go-nano/protocol"
)

// RetryExceeded represents the error returned when client fail
// to send request after retry.
var RetryExceeded = errors.New("retry attempts exceeded")

type retryExt struct {
	next     protocol.Sender
	maxTries int
	backoff  backoff.BackOff
}

// NewRetryExt returns an extension that wraps a client to retry
// failed requests using given backoff.
// This extension assuming nothing about whether service is idempotent
// or not, only request that are not sent are retried.
// FIXME: Middlewares are client specific this mean that all requests share the same backoff.
func NewRetryExt(maxTries int, backoff backoff.BackOff) Extension {
	return func(next protocol.Sender) protocol.Sender {
		return &retryExt{
			next:     next,
			maxTries: maxTries,
			backoff:  backoff,
		}
	}
}

func (e *retryExt) Send(endpoint string, req *protocol.Request) (*protocol.Response, error) {
	var next time.Duration
	var err error
	var resp *protocol.Response

	// TODO: Shared backoff !?
	e.backoff.Reset()
	for i := 1; i <= e.maxTries; i++ {
		resp, err = e.next.Send(endpoint, req)

		if err == nil {
			break
		}

		next = e.backoff.NextBackOff()
		if next == backoff.Stop {
			break
		}
		time.Sleep(next)

		log.Debug("Sending request fail %s:%s: %s", endpoint, req, err)
	}

	if err != nil {
		return nil, RetryExceeded
	}
	return resp, nil
}
