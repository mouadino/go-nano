package amqp

import (
	"github.com/mouadino/go-nano/transport"
	"github.com/streadway/amqp"
)

type AMQPResponseWriter struct {
	trans    *AMQPTransport
	delivery amqp.Delivery
}

func NewAMQPResponseWriter(trans *AMQPTransport, delivery amqp.Delivery) transport.ResponseWriter {
	return &AMQPResponseWriter{
		trans:    trans,
		delivery: delivery,
	}
}

func (rw *AMQPResponseWriter) Write(data interface{}) error {
	err := rw.trans.sendReply(rw.delivery.ReplyTo, rw.delivery.CorrelationId, data.([]byte))
	if err != nil {
		return err
	}
	rw.delivery.Ack(false)
	return nil
}
