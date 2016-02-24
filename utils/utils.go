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

// GetListener returns a net.Listener to a random port using TCP,
// when an empty string is passed argument, GetListener will try to guess
// the IP address to use by getting first ip in network interface associated
// with the default gateway.
func GetListener(addr string) (net.Listener, error) {
	var err error

	if addr == "" {
		addr, err = GetExternalIP()
	}
	if err != nil {
		return nil, err
	}
	return net.Listen("tcp", fmt.Sprintf("%s:0", addr))
}

// TODO: https://github.com/twitter/finagle/blob/develop/finagle-core/src/main/scala/com/twitter/finagle/util/InetSocketAddressUtil.scala#L13

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
	return "", errors.New("fail to guess external ip")
}
