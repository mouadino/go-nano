package transport

import (
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
)

const defaultExchange = "nano"

type AMQPResponseWriter struct {
	trans    *AMQPTransport
	delivery amqp.Delivery
}

func NewAMQPResponseWriter(trans *AMQPTransport, delivery amqp.Delivery) ResponseWriter {
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

type AMQPTransport struct {
	url             string
	exchange        string
	conn            *amqp.Connection
	reqs            chan Request
	listenQueue     string
	logger          *log.Logger
	pendingRequests map[string]chan []byte
	listening       bool
}

func NewAMQPTransport(url string) Transport {
	return NewCustomAMQPTransport(url, defaultExchange, path.Base(os.Args[0]))
}

func NewCustomAMQPTransport(url, exchange, listenQueue string) Transport {
	return &AMQPTransport{
		url:             url,
		exchange:        exchange, // TODO: Is this used ?
		logger:          log.New(),
		listenQueue:     listenQueue, // FIXME: Should be unique per service.
		reqs:            make(chan Request),
		pendingRequests: make(map[string]chan []byte),
	}
}

func (trans *AMQPTransport) Listen() error {
	err := trans.declareQueue(trans.listenQueue)
	if err != nil {
		return err
	}
	trans.logger.Info("Listening on ", trans.url, trans.listenQueue)
	go trans.consumeMessages()
	trans.listening = true
	return nil
}

func (trans *AMQPTransport) consumeMessages() {
	for {
		ch, err := trans.getChannel()
		if err != nil {
			log.Info("fail to get channel retrying: ", err)
			continue
		}
		defer ch.Close()
		msgs, err := ch.Consume(
			trans.listenQueue,
			trans.exchange,
			false, // TODO: In case of client autoAck is needed.
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Info("fail to consume messages retrying: ", err)
			continue
		}
		for msg := range msgs {
			if msg.ReplyTo == "" {
				trans.logger.Debug("Reply received for ", msg.CorrelationId)
				// TODO: May fail.
				trans.pendingRequests[msg.CorrelationId] <- msg.Body
			} else {
				trans.logger.Debug("New request received")
				trans.reqs <- Request{
					Body: msg.Body,
					Resp: NewAMQPResponseWriter(trans, msg),
				}
			}
		}
	}
}

func (trans *AMQPTransport) Receive() <-chan Request {
	return trans.reqs
}

func (trans *AMQPTransport) Send(endpoint string, message []byte) ([]byte, error) {
	if !trans.listening {
		err := trans.Listen()
		if err != nil {
			return []byte{}, err
		}
	}
	correlationID, err := trans.sendRequest(endpoint, message)
	if err != nil {
		return []byte{}, err
	}
	// TODO: Synchronization.
	trans.logger.Debug("Waiting for ", correlationID)
	body := <-trans.pendingRequests[correlationID]
	trans.logger.Debug("Received  ", correlationID)
	// TODO: close channel and delete key.
	return body, nil
}

func (trans *AMQPTransport) sendRequest(routingKey string, message []byte) (string, error) {
	var correlationID string
	if u, err := uuid.NewV4(); err != nil {
		return "", err
	} else {
		correlationID = u.String()
	}
	// TODO: Expire time ?
	msg := amqp.Publishing{
		ContentType:   "text/plain",
		Body:          message,
		CorrelationId: correlationID,
		ReplyTo:       trans.listenQueue,
	}
	trans.logger.Debug("Send request")
	err := trans.publishMessage(routingKey, msg, false)
	if err != nil {
		return "", err
	}
	// TODO: Should not leak ? remove after some time.
	trans.pendingRequests[correlationID] = make(chan []byte)
	return correlationID, nil
}

func (trans *AMQPTransport) sendReply(routingKey, correlationID string, message []byte) error {
	msg := amqp.Publishing{
		ContentType:   "text/plain",
		Body:          message,
		CorrelationId: correlationID,
	}
	trans.logger.Debug("Send reply to ", routingKey)
	return trans.publishMessage(routingKey, msg, true)
}

func (trans *AMQPTransport) publishMessage(routingKey string, msg amqp.Publishing, direct bool) error {
	ch, err := trans.getChannel()
	if err != nil {
		return err
	}
	defer ch.Close()

	exchange := "" //trans.exchange
	if direct {
		exchange = ""
	}
	trans.logger.Debug("Publish message to ", routingKey)
	return ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		msg,
	)
}

func (trans *AMQPTransport) declareQueue(name string) error {
	ch, err := trans.getChannel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	// TODO: Binding ? needed ? declare exchange first.
	//ch.QueueBind(name, name, trans.exchange, false, nil)
	return err
}

func (trans *AMQPTransport) getChannel() (*amqp.Channel, error) {
	conn, err := trans.getConn()
	if err != nil {
		return nil, err
	}
	return conn.Channel()
}

func (trans *AMQPTransport) getConn() (*amqp.Connection, error) {
	if trans.conn != nil {
		// TODO: Reconnect when connection drop.
		return trans.conn, nil
	}
	conn, err := amqp.Dial(trans.url)
	if err != nil {
		return nil, err
	}
	trans.conn = conn
	return conn, nil
}
