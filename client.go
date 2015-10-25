package nano

import (
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
	"github.com/mouadino/go-nano/utils"
)

type remoteClient struct {
	endpoint string
	proto    protocol.Protocol
}

func Client(endpoint string) *remoteClient {
	return &remoteClient{
		endpoint: endpoint,
		proto:    protocol.NewJSONRPCProtocol(transport.NewHTTPTransport()),
	}
}

func (client *remoteClient) Call(method string, params ...interface{}) (interface{}, error) {
	req := protocol.Request{
		Method: method,
		Params: utils.ParamsFormat(params...),
	}
	resp, err := client.proto.SendRequest(client.endpoint, &req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
