package utils

import (
	"fmt"

	"github.com/mouadino/go-nano/protocol"
)

func ParamsFormat(ps ...interface{}) protocol.Params {
	params := map[string]interface{}{}
	for i, v := range ps {
		params[fmt.Sprintf("_%d", i)] = v
	}
	return params
}
