package handler

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/context"

	"github.com/mouadino/go-nano/protocol"
)

type params []reflect.Value

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type structHandler struct {
	svc     interface{}
	methods map[string]methodHandler
}

func newStructHandler(svc interface{}) *structHandler {
	methods := map[string]methodHandler{}
	svcType := reflect.TypeOf(svc)
	for i := 0; i < svcType.NumMethod(); i++ {
		method := svcType.Method(i)
		if isRPCMethod(method) {
			methods[strings.ToLower(method.Name)] = methodHandler{svc, svcType.Method(i)}
		}
	}
	return &structHandler{
		svc:     svc,
		methods: methods,
	}
}

func isRPCMethod(meth reflect.Method) bool {
	rune, _ := utf8.DecodeRuneInString(meth.Name)
	isExported := unicode.IsUpper(rune)
	retTypeCorrect := meth.Type.NumOut() == 2 && meth.Type.Out(1) == typeOfError

	return isExported && retTypeCorrect
}

func (h *structHandler) Handle(ctx context.Context, req *protocol.Request, resp *protocol.Response) {
	name := req.Method
	fh, ok := h.methods[name]
	if !ok {
		resp.Error = protocol.UnknownMethod
		return
	}
	fh.Handle(ctx, req, resp)
}

type methodHandler struct {
	svc    interface{}
	method reflect.Method
}

func (h *methodHandler) Handle(ctx context.Context, req *protocol.Request, resp *protocol.Response) {
	params, err := h.parseParams(req)
	if err != nil {
		resp.Error = err
		return
	}
	data, err := h.call(params)
	if err != nil {
		resp.Error = err
		return
	}
	resp.Body = data
}

func (h *methodHandler) parseParams(req *protocol.Request) (params, error) {
	params := make(params, len(req.Params)+1)
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

func (h *methodHandler) call(ps params) (interface{}, error) {
	ret := h.method.Func.Call(ps)
	data := make([]interface{}, len(ret))
	for i, v := range ret {
		data[i] = v.Interface()
	}
	var err error
	if data[1] == nil {
		err = nil
	} else {
		err = data[1].(error)
	}

	return data[0], err
}
