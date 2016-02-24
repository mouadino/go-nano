# Anatomy of an RPC framework

## RPC VS Rest


## Components


### Transport


### Protocol


### Server

#### Announcing


#### Middlewares

### Client

A client

#### Discovery


#### Extensions

Extension are components that implements the `transport.Sender` and can
be stacked by decorating yet another `transport.Sender` and so on, the
end result is yet an extension.

Extensions are used in conjuction with client,


