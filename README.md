[![Build Status](https://travis-ci.org/mouadino/go-nano.svg?branch=master)](https://travis-ci.org/mouadino/go-nano)
[![Coverage Status](https://coveralls.io/repos/mouadino/go-nano/badge.svg?branch=master&service=github)](https://coveralls.io/github/mouadino/go-nano?branch=master)
[![GoDoc](https://godoc.org/github.com/mouadino/go-nano?status.svg)](https://godoc.org/github.com/mouadino/go-nano)

# go-nano

A framework for writting web services in Go.

## Design

### Rational:

Developers should focus on writing the business logic.

### Highlight:

- Transport agnostic (TCP, HTTP, AMQP ...)
- Protocol agnostic (JSON-RPC, ProtocolBuffer, ...)
- Plain Old Go Struct for business logic
- Handle boilerplate for setting up a service
- Convention over configuration
- RPC.

### Components:

                                        +------------ +     +---------+
                                        | Middlewares | <-> | Handler |
                                        +------------ +     +---------+
    +-----------+     +----------+
    | Transport | <-> | Protocol | <->
    +-----------+     +----------+
                                        +------------ +     +---------+
                                        | Extensions  | <-> | Client  |
                                        +------------ +     +---------+


## Features:

Transports:

- [X] AMQP Transport (Minimal)
- [ ] HTTP2 Transport
- [X] HTTP Transport

Protocols:

- [X] JSON-RPC Protocol (Minimal)
- [ ] ProtocolBuffer
- [X] Lymph Protocol (http://lymph.readthedocs.org/en/latest/protocol.html https://github.com/mouadino/go-lymph)

Client:

- [X] Client
- [X] Circuit Breaker
- [X] Timeout
- [X] Retry
- [X] Remote errors
- [X] Async

Server:

- [X] Multiple services namespaces (a la net/rpc)
- [X] Tracing middleware
- [ ] Rate limit middleware
- [X] Registration

Command lines:

- [ ] gonano request
- [ ] gonano get
- [ ] gonano list

Misc:

- [ ] Logging (https://godoc.org/gopkg.in/inconshreveable/log15.v2)
- [ ] Context (https://blog.golang.org/context)
- [ ] Metrics (https://github.com/rcrowley/go-metrics)
- [X] Discovery
- [ ] Testing tools (mocking, fake services)
- [ ] Protocol/Transport negotiation
- [ ] PubSub
- [ ] Notify (One way messages)
- [X] Add to travis/circleCI ...
- [X] GoDoc
- [ ] More complete example (e.g. user service)

## Examples:

Please check [examples] folder

[examples]: https://github.com/mouadino/go-nano/tree/master/examples
