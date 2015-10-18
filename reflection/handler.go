package reflection

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mouadino/go-nano/interfaces"
)

var publicMethod = regexp.MustCompile("^[A-Z]")

type StructHandler struct {
	svc     interface{}
	methods map[string]MethodHandler
}

func FromStruct(svc interface{}) *StructHandler {
	methods := map[string]MethodHandler{}
	svcType := reflect.TypeOf(svc)
	for i := 0; i < svcType.NumMethod(); i++ {
		method := svcType.Method(i)
		if publicMethod.MatchString(method.Name) {
			methods[strings.ToLower(method.Name)] = MethodHandler{svc, svcType.Method(i)}
		}
	}
	return &StructHandler{
		svc:     svc,
		methods: methods,
	}
}

func (h *StructHandler) Handle(resp interfaces.ResponseWriter, req *interfaces.Request) error {
	name := req.Method
	fh, ok := h.methods[name]
	if !ok {
		// TODO: resp.WriteError(...) !?
		fmt.Printf("unknown method %s\n", name)
		return fmt.Errorf("unknown method %s", name)
	}
	return fh.Handle(resp, req)
}

type MethodHandler struct {
	svc    interface{}
	method reflect.Method
}

func (h *MethodHandler) Handle(resp interfaces.ResponseWriter, req *interfaces.Request) error {
	in := make([]reflect.Value, len(req.Params)+1)
	in[0] = reflect.ValueOf(h.svc)
	i := 1
	// TODO: Order parameters.
	for _, v := range req.Params {
		in[i] = reflect.ValueOf(v)
		i++
	}
	fmt.Printf("calling %s with %s\n", h.method, in)
	ret := h.method.Func.Call(in)
	fmt.Printf("returning %s\n", ret)
	data := make([]interface{}, len(ret))
	for i, v := range ret {
		data[i] = v.Interface()
	}
	if len(data) == 1 {
		resp.Write(data[0])
	} else {
		resp.Write(data)
	}
	return nil
}
