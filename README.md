[![Build Status](https://travis-ci.org/mouadino/go-nano.svg?branch=master)](https://travis-ci.org/mouadino/go-nano)
[![Coverage Status](https://coveralls.io/repos/mouadino/go-nano/badge.svg?branch=master&service=github)](https://coveralls.io/github/mouadino/go-nano?branch=master)

# go-nano

A framework for writting web services in Go.

## Design

### Rational

Developers should focus on writing the business logic.

### Highlight

- Transport agnostic (TCP, HTTP, AMQP ...)
- Protocol agnostic (JSON-RPC, ProtocolBuffer, ...)
- Plain Old Go Struct for business logic
- Handle boilerplate for setting up a service
- Convention over configuration
- RPC.

### Components

                                        +------------ +     +---------+
                                        | Middlewares | <-> | Handler |
                                        +------------ +     +---------+
    +-----------+     +----------+
    | Transport | <-> | Protocol | <->
    +-----------+     +----------+
                                        +------------ +     +---------+
                                        | Extensions  | <-> | Client  |
                                        +------------ +     +---------+


## Features

Transports:

- [X] AMQP Transport (Minimal)
- [ ] HTTP2 Transport
- [X] HTTP Transport

Protocols:

- [X] JSON-RPC Protocol (Minimal)
- [ ] ProtocolBuffer
- [X] Lymph Protocol (https://github.com/mouadino/go-lymph)

Client:

- [X] Client
- [ ] Command line
- [X] Circuit Breaker
- [X] Timeout
- [X] Retry
- [ ] Remote errors
- [X] Async

Server:

- [X] Multiple services namespaces (a la net/rpc)
- [X] Tracing middleware
- [ ] Rate limit middleware
- [X] Registration

Misc:

- [ ] Configuration
- [ ] Logging
- [ ] Context
- [ ] Metrics
- [X] Discovery
- [ ] Testing tools
- [ ] Protocol/Transport negotiation
- [ ] Notification (a.k.a One way message)
- [X] Add to travis/circleCI ...
- [ ] GoDoc
- [ ] More complete example (e.g. user service)

## Examples

Please check [examples] folder

[examples]: https://github.com/mouadino/go-nano/tree/master/examples
