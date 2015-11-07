package reflection

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/mouadino/go-nano/protocol"
)

var publicMethod = regexp.MustCompile("^[A-Z]")

type Params []reflect.Value

type StructHandler struct {
	svc     interface{}
	methods map[string]MethodHandler
}

func FromStruct(svc interface{}) *StructHandler {
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

func isRPCMethod(name string) bool {
	return publicMethod.MatchString(name) && !strings.HasPrefix(name, "Nano")
}

func (h *StructHandler) Handle(resp protocol.ResponseWriter, req *protocol.Request) {
	name := req.Method
	fh, ok := h.methods[name]
	if !ok {
		resp.WriteError(protocol.UnknownMethod)
		return
	}
	fh.Handle(resp, req)
}

type MethodHandler struct {
	svc    interface{}
	method reflect.Method
}

func (h *MethodHandler) Handle(resp protocol.ResponseWriter, req *protocol.Request) {
	defer h.recoverFromError(resp)

	params, err := h.parseParams(req)
	if err != nil {
		resp.WriteError(err)
		return
	}
	// TODO: Returning error !? .NumOut() ... ?
	data := h.call(params)
	resp.Write(data)
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

func (h *MethodHandler) recoverFromError(resp protocol.ResponseWriter) {
	if err := recover(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Recovered from handler")
		// TODO: WithTrace() ? https://github.com/Sirupsen/logrus/pull/284
		debug.PrintStack()
		resp.WriteError(protocol.InternalError)
	}
}
