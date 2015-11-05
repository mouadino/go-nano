package client

type Client interface {
	Call(string, ...interface{}) (interface{}, error)
}

type ClientExtension func(Client) Client
