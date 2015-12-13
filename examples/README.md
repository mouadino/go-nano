# Examples

This folder contain examples on how to use go-nano:

## upper

A service that expose a function that can transform a string to upper
cases.

## demo

It's a client that talk with upper service.

# Running

You can use [docker-compose](https://docs.docker.com/compose/) to run
the services above, the docket-compose.yml already contain dependencies
needed like zookeeper, rabbitmq.

Start by executing:

    $ docker-compose up

