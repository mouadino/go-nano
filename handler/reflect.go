package handler

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mouadino/go-nano/protocol"
)

var publicMethod = regexp.MustCompile("^[A-Z]")

type Configurable interface {
	NanoConfigure(interface{}) error
}

type Startable interface {
	NanoStart() error
	NanoStop() error
}

type Params []reflect.Value

func Reflect(svc interface{}) Handler {
	if hdlr, ok := svc.(Handler); ok {
		return hdlr
	}
	return NewStructHandler(svc)
}

type StructHandler struct {
	svc     interface{}
	methods map[string]MethodHandler
}

func NewStructHandler(svc interface{}) *StructHandler {
	methods := map[string]MethodHandler{}
	svcType := reflect.TypeOf(svc)
	for i := 0; i < svcType.NumMethod(); i++ {
		method := svcType.Method(i)
		if isRPCMethod(method.Name) {
			methods[strings.ToLower(method.Name)] = MethodHandler{svc, svcType.Method(i)}
		}
	}
	return &StructHandler{
		svc:     svc,
		methods: methods,
	}
}

// TODO: https://github.com/golang/go/blob/master/src/net/rpc/server.go#L203
func isRPCMethod(name string) bool {
	return publicMethod.MatchString(name) && !strings.HasPrefix(name, "Nano")
}

func (h *StructHandler) Handle(resp protocol.ResponseWriter, req *protocol.Request) {
	name := req.Method
	fh, ok := h.methods[name]
	if !ok {
		resp.SetError(protocol.UnknownMethod)
		return
	}
	fh.Handle(resp, req)
}

type MethodHandler struct {
	svc    interface{}
	method reflect.Method
}

func (h *MethodHandler) Handle(resp protocol.ResponseWriter, req *protocol.Request) {
	params, err := h.parseParams(req)
	if err != nil {
		resp.SetError(err)
		return
	}
	// TODO: Returning error !? .NumOut() ... ?
	data := h.call(params)
	resp.Set(data)
}

func (h *MethodHandler) parseParams(req *protocol.Request) (Params, error) {
	params := make(Params, len(req.Params)+1)
	params[0] = reflect.ValueOf(h.svc)
	for i := 0; ; i++ {
		v, ok := req.Params[fmt.Sprintf("_%d", i)]
		if !ok {
			break
		}
		params[i+1] = reflect.ValueOf(v)
	}
	if h.method.Type.NumIn() != len(params) {
		return params, protocol.ParamsError
	}
	return params, nil
}

func (h *MethodHandler) call(params Params) interface{} {
	ret := h.method.Func.Call(params)
	data := make([]interface{}, len(ret))
	for i, v := range ret {
		data[i] = v.Interface()
	}
	// XXX Can we do better ?
	if len(data) == 1 {
		return data[0]
	}
	return data
}
