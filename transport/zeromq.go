package transport

/*
FIXME: zeromq installation.
import "github.com/zeromq/goczmq"

type ZeroMQTransporter struct {
	endpoint  string
	channeler *goczmq.Channeler
}

func NewZeroMQTransporter(endpoint string) *ZeroMQTransporter {
	return &ZeroMQTransporter{
		endpoint:  endpoint,
		channeler: goczmq.NewRouterChanneler(endpoint),
	}
}

func (t *ZeroMQTransporter) Send(endpoint, data []byte) {
	t.sendChan <- []byte{[]byte(endpoint), data}
}

func (t *ZeroMQTransporter) Receive() <-chan [][]byte {
	return t.channeler.RecvChan
}
*/
