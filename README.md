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
                                        +------------ +     +---------+     +---------+
                                        | Middlewares | <-> | Handler | <-> | Service |
                                        +------------ +     +---------+     +---------+
    +-----------+     +----------+
    | Transport | <-> | Protocol | <->
    +-----------+     +----------+
                                        +------------ +     +---------+
                                        | Extensions  | <-> | Client  |
                                        +------------ +     +---------+

## Features

Transports:

- [ ] AMPQ Transport
- [ ] ZeroMQ Transport
- [X] HTTP Transport

Protocols:

- [X] JSON-RPC Protocol (Minimal)
- [ ] MSGPACK-RPC Protocol
- [ ] Lymph Protocol

Client:

- [X] Client
- [X] Command line
- [ ] Circuit Breaker
- [X] Timeout
- [ ] Retry
- [ ] Remote errors
- [ ] Async

Misc:

- [ ] Configuration
- [X] Middlewares (Tracing, Rate limit ...)
- [ ] Logging
- [ ] Context
- [ ] Metrics
- [X] Discovery
- [ ] Testing tools
- [ ] Notification (a.k.a One way message)

## Examples

Please check [examples] folder

[examples]: https://github.com/mouadino/go-nano/tree/master/examples
