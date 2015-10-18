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
- RPC (Focus only on request/response pattern).

### Schema

                                                                           +---------+
                                                                           | Service |
                                                                           +---------+
    +-----------+     +----------+     +------------ +      +---------+
    | Transport | <-> | Protocol | <-> | Middlewares | <->  | Handler | <->
    +-----------+     +----------+     +------------ +      +---------+
                                                                           +---------+
                                                                           | Client  |
                                                                           +---------+

## TODO

Transports:

- [ ] AMPQ Transport
- [ ] ZeroMQ Transport
- [ ] HTTP Transport

Protocols:

- [ ] JSON-RPC Protocol
- [ ] MSGPACK-RPC Protocol
- [ ] Lymph Protocol

Client

- [ ] Client
- [ ] Command line
- [ ] Circuit Breaker

Misc:

- [ ] Configuration handling
- [ ] Middlewares (Tracing, Rate limit ...)
- [ ] Logging
- [ ] Metrics
- [ ] Discovery
- [ ] Testing tools
