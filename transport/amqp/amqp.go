package amqp

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/mouadino/go-nano/transport"
	"github.com/pborman/uuid"
	"github.com/streadway/amqp"
)

type AMQPTransport struct {
	url             string
	exchange        string
	conn            *amqp.Connection
	reqs            chan transport.Request
	listenQueue     string
	logger          *log.Logger
	pendingRequests map[string]chan []byte
	listening       bool
}

// Exchange option to set exchange name, default "nano".
func Exchange(name string) func(*AMQPTransport) {
	return func(t *AMQPTransport) {
		t.exchange = name
	}
}

// QueueName option to set queue name, default name of go binary.
func QueueName(name string) func(*AMQPTransport) {
	return func(t *AMQPTransport) {
		t.listenQueue = name
	}
}

// New returns a Transport that use AMQP to send/receive RPC messages.
func New(url string, options ...func(*AMQPTransport)) transport.Transport {
	t := &AMQPTransport{
		url:             url,
		exchange:        "nano", // TODO: Is this used ?
		logger:          log.New(),
		listenQueue:     path.Base(os.Args[0]), // FIXME: Should be unique per service.
		reqs:            make(chan transport.Request),
		pendingRequests: make(map[string]chan []byte),
	}

	for _, opt := range options {
		opt(t)
	}
	return t
}

func (trans *AMQPTransport) Listen() error {
	err := trans.declareQueue(trans.listenQueue)
	if err != nil {
		return err
	}
	trans.logger.Info("Listening on ", trans.url, " ", trans.listenQueue)
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
				replyCh, ok := trans.pendingRequests[msg.CorrelationId]
				if ok {
					replyCh <- msg.Body
				} else {
					trans.logger.Error("no request waiting for %s", msg)
				}
				msg.Ack(true)
				// TODO: Remove trans.pendingRequests.
				// TODO: Mutex.
			} else {
				trans.logger.Debug("New request received")
				trans.reqs <- transport.Request{
					Body: msg.Body,
					Resp: NewAMQPResponseWriter(trans, msg),
				}
			}
		}
	}
}

func (trans *AMQPTransport) Receive() <-chan transport.Request {
	return trans.reqs
}

func (trans *AMQPTransport) Send(endpoint string, message io.Reader) ([]byte, error) {
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

func (trans *AMQPTransport) sendRequest(routingKey string, message io.Reader) (string, error) {
	correlationID := uuid.New()
	body, err := ioutil.ReadAll(message)
	if err != nil {
		return "", err
	}
	// TODO: Expire time ?
	// TODO: Persistence ?
	msg := amqp.Publishing{
		ContentType:   "text/plain",
		Body:          body,
		CorrelationId: correlationID,
		ReplyTo:       trans.listenQueue,
	}
	trans.logger.Debug("Send request")
	err = trans.publishMessage(routingKey, msg, false)
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
