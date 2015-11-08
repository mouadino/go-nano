# go-nano

A framework for writting web services in Go.

## Design

### Rational

Developers should focus on writing the business logic.

### Highlight

- Transport agnostic (TCP, HTTP, AMQP ...)
- Protocol agnostic (JSON-RPC, MSGPACK-RPC, ...)
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

- [X] JSON-RPC Protocol
- [ ] MSGPACK-RPC Protocol
- [ ] Lymph Protocol

Client:

- [X] Client
- [X] Command line
- [ ] Circuit Breaker

Misc:

- [ ] Remote errors
- [ ] Configuration handling
- [X] Middlewares (Tracing, Rate limit ...)
- [ ] Logging
- [ ] Metrics
- [ ] Discovery
- [ ] Testing tools
- [ ] Notification (a.k.a One way message)

## Examples

Please check [examples] folder

[examples]: https://github.com/mouadino/go-nano/tree/master/examples
