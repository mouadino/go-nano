package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"golang.org/x/net/context"
)

type rpcHandler struct {
	hdlr   handler.Handler
	protos protocolMap
}

func newRPCHandler(hdlr handler.Handler) *rpcHandler {
	return &rpcHandler{
		hdlr:   hdlr,
		protos: protocolMap{},
	}
}

func (h *rpcHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	proto, err := h.getProtocol(rw, req)
	if err != nil {
		return
	}

	protoReq, err := h.formatProtocolRequest(proto, rw, req)
	if err != nil {
		return
	}
	protoResp := &protocol.Response{}
	ctx := context.Background()

	h.hdlr.Handle(ctx, protoReq, protoResp)

	h.formatResponse(proto, protoResp, rw)
}

func (h *rpcHandler) getProtocol(rw http.ResponseWriter, req *http.Request) (protocol.Protocol, error) {
	proto := h.protos.Get(req)
	if proto == nil {
		rw.WriteHeader(http.StatusUnsupportedMediaType)
		io.WriteString(rw, "Internal error")
		// TODO: protocol ?
		return nil, fmt.Errorf("unknown protocol")
	}
	return proto, nil
}

func (h *rpcHandler) formatProtocolRequest(proto protocol.Protocol, rw http.ResponseWriter, req *http.Request) (*protocol.Request, error) {
	protoReq, err := proto.DecodeRequest(req.Body, header.Header(req.Header))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		io.WriteString(rw, "Internal error")
		return nil, fmt.Errorf("internal error")
	}
	return protoReq, err
}

func (h *rpcHandler) formatResponse(proto protocol.Protocol, resp *protocol.Response, rw http.ResponseWriter) {
	body, err := proto.EncodeResponse(resp)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		io.WriteString(rw, "Internal error")
	} else {
		status := getHTTPStatus(resp.Error)
		rw.WriteHeader(status)
		for k, vs := range resp.Header {
			for _, v := range vs {
				rw.Header().Add(k, v)
			}
		}
		contentType := getContentType(proto)
		rw.Header().Set("Content-Type", contentType)
		rw.Write(body)
	}
}

func getHTTPStatus(err error) int {
	switch err {
	case nil:
		return http.StatusOK
	case protocol.UnknownMethod:
		return http.StatusMethodNotAllowed
	case protocol.ParamsError:
		return http.StatusBadRequest
	case protocol.InternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
