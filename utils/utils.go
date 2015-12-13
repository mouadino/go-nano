package utils

import (
	"errors"
	"fmt"
	"net"

	"github.com/mouadino/go-nano/protocol"
)

func ParamsFormat(ps ...interface{}) protocol.Params {
	params := map[string]interface{}{}
	for i, v := range ps {
		params[fmt.Sprintf("_%d", i)] = v
	}
	return params
}

// GetExternalIP returns external ip of local node.
func GetExternalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok {
			// TODO: Use default gateway.
			if !ip.IP.IsLoopback() && ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}
	return "", errors.New("not found")
}
