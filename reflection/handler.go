package reflection

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
)

var publicMethod = regexp.MustCompile("^[A-Z]")

type StructHandler struct {
	svc     interface{}
	methods map[string]MethodHandler
}

func FromStruct(svc interface{}) *StructHandler {
	methods := map[string]MethodHandler{}
	svcType := reflect.TypeOf(svc)
	log.Printf("%s %s", svc, svcType.NumMethod())
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

func (h *StructHandler) Handle(resp transport.ResponseWriter, req *protocol.Request) {
	name := req.Method
	fh, ok := h.methods[name]
	if !ok {
		log.Printf("unknown method %s\n", name)
		resp.WriteError(protocol.UnknownMethod)
		return
	}
	fh.Handle(resp, req)
}

type MethodHandler struct {
	svc    interface{}
	method reflect.Method
}

func (h *MethodHandler) Handle(resp transport.ResponseWriter, req *protocol.Request) {
	defer h.recoverFromError(resp)
	in := make([]reflect.Value, len(req.Params)+1)
	log.Printf("Parameters %s", req.Params)
	in[0] = reflect.ValueOf(h.svc)
	for i := 0; ; i++ {
		v, ok := req.Params[fmt.Sprintf("_%d", i)]
		log.Printf("Get parameter %s: %s", fmt.Sprintf("_%d", i), v)
		if !ok {
			break
		}
		in[i+1] = reflect.ValueOf(v)
	}
	if h.method.Type.NumIn() != len(in) {
		resp.WriteError(protocol.ParamsError)
		return
	}
	log.Printf("calling %s with %s\n", h.method, in)
	// TODO: Returning error !? .NumOut() ... ?
	ret := h.method.Func.Call(in)
	log.Printf("returning %s\n", ret)
	data := make([]interface{}, len(ret))
	for i, v := range ret {
		data[i] = v.Interface()
	}
	if len(data) == 1 {
		resp.Write(data[0])
	} else {
		resp.Write(data)
	}
}

func (h *MethodHandler) recoverFromError(resp transport.ResponseWriter) {
	if err := recover(); err != nil {
		log.Println("Recovered from handler error", err)
		resp.WriteError(protocol.InternalError)
	}
}
