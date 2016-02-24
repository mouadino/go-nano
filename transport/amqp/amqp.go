package amqp

import (
	"bytes"
	"os"
	"path"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
	"github.com/pborman/uuid"
	"github.com/streadway/amqp"
)

type amqpTransport struct {
	url             string
	exchange        string
	conn            *amqp.Connection
	listenQueue     string
	logger          *log.Logger
	pendingRequests map[string]chan []byte
	mu              sync.RWMutex
	listening       bool
	proto           protocol.Protocol
	hdlr            handler.Handler
}

// Exchange option to set exchange name, default "nano".
func Exchange(name string) func(*amqpTransport) {
	return func(t *amqpTransport) {
		t.exchange = name
	}
}

// QueueName option to set queue name, default name of go binary.
func QueueName(name string) func(*amqpTransport) {
	return func(t *amqpTransport) {
		t.listenQueue = name
	}
}

// New returns a Transport that use AMQP to send/receive RPC messages.
func New(url string, options ...func(*amqpTransport)) transport.Transport {
	t := &amqpTransport{
		url:             url,
		exchange:        "nano", // TODO: Is this used ?
		logger:          log.New(),
		listenQueue:     path.Base(os.Args[0]), // FIXME: Should be unique per service.
		pendingRequests: make(map[string]chan []byte),
	}

	for _, opt := range options {
		opt(t)
	}
	return t
}

// FIXME: Add handler not set.
func (trans *amqpTransport) AddHandler(proto protocol.Protocol, hdlr handler.Handler) {
	trans.proto = proto
	trans.hdlr = hdlr
}

func (trans *amqpTransport) Listen() error {
	err := trans.declareQueue(trans.listenQueue)
	if err != nil {
		return err
	}
	trans.logger.Info("Listening on ", trans.url, " ", trans.listenQueue)
	trans.listening = true
	return nil
}

func (trans *amqpTransport) Serve() error {
	trans.consumeMessages()
	return nil
}

func (trans *amqpTransport) consumeMessages() {
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
				go trans.handleResponse(msg)
			} else {
				go trans.handleRequest(msg)
			}
		}
	}
}

func (trans *amqpTransport) handleResponse(msg amqp.Delivery) {
	trans.mu.Lock()
	defer trans.mu.Unlock()
	trans.logger.Debug("Reply received for ", msg.CorrelationId)
	replyCh, ok := trans.pendingRequests[msg.CorrelationId]
	if ok {
		replyCh <- msg.Body
		close(replyCh)
	} else {
		trans.logger.Error("no request waiting for %s", msg)
	}
	msg.Ack(false)
	delete(trans.pendingRequests, msg.CorrelationId)
}

func (trans *amqpTransport) handleRequest(msg amqp.Delivery) {
	trans.logger.Debug("New request received")
	// TODO: Handle error
	req, _ := trans.proto.DecodeRequest(bytes.NewReader(msg.Body))
	resp := &protocol.Response{
		Header: header.Header{},
	}

	trans.hdlr.Handle(resp, req)

	// TODO: Handle error.
	body, _ := trans.proto.EncodeResponse(resp)
	_ = trans.sendReply(msg.ReplyTo, msg.CorrelationId, body)

	msg.Ack(false)
}

func (trans *amqpTransport) Send(endpoint string, req *protocol.Request) (*protocol.Response, error) {
	// TODO: Do the listening somewhere else.
	if !trans.listening {
		err := trans.Listen()
		if err != nil {
			return nil, err
		}
	}

	body, err := trans.proto.EncodeRequest(req)
	if err != nil {
		return nil, err
	}
	replyCh, err := trans.sendRequest(endpoint, body)
	if err != nil {
		return nil, err
	}
	// TODO: Synchronization.
	respBody := <-replyCh
	// TODO: close channel and delete key.
	return trans.proto.DecodeResponse(bytes.NewReader(respBody))
}

func (trans *amqpTransport) sendRequest(routingKey string, body []byte) (<-chan []byte, error) {
	correlationID := uuid.New()
	// TODO: Expire time ?
	// TODO: Persistence ?
	msg := amqp.Publishing{
		// TODO: Use getContentType(trans.proto)
		ContentType:   "text/plain",
		Body:          body,
		CorrelationId: correlationID,
		ReplyTo:       trans.listenQueue,
	}
	trans.logger.Debug("Send request")
	err := trans.publishMessage(routingKey, msg, false)
	if err != nil {
		return nil, err
	}
	// TODO: Should not leak ? remove after some time.
	replyCh := make(chan []byte, 1)
	trans.pendingRequests[correlationID] = replyCh
	return replyCh, nil
}

func (trans *amqpTransport) sendReply(routingKey, correlationID string, message []byte) error {
	msg := amqp.Publishing{
		ContentType:   "text/plain",
		Body:          message,
		CorrelationId: correlationID,
	}
	trans.logger.Debug("Send reply to ", routingKey)
	return trans.publishMessage(routingKey, msg, true)
}

func (trans *amqpTransport) publishMessage(routingKey string, msg amqp.Publishing, direct bool) error {
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

func (trans *amqpTransport) declareQueue(name string) error {
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

func (trans *amqpTransport) getChannel() (*amqp.Channel, error) {
	conn, err := trans.getConn()
	if err != nil {
		return nil, err
	}
	return conn.Channel()
}

func (trans *amqpTransport) getConn() (*amqp.Connection, error) {
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
