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
	sender   protocol.Sender
	maxTries int
	backoff  backoff.BackOff
}

// NewRetryExt returns an extension that wraps a client to retry
// failed requests using given backoff.
// This extension assuming nothing about whether service is idempotent
// or not, only request that are not sent are retried.
func NewRetryExt(maxTries int, backoff backoff.BackOff) Extension {
	return func(s protocol.Sender) protocol.Sender {
		return &retryExt{
			sender:   s,
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
		resp, err = e.sender.Send(endpoint, req)

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
