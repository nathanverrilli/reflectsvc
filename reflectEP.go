package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
	"strings"
)

// For each method, we define request and response structs
type reflectRequest struct {
	S string `json:"s"`
}

type reflectResponse struct {
	S   string `json:"s"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

func makeReflectEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(reflectRequest)
		v, err := svc.Reflect(req.S)
		if err != nil {
			return reflectResponse{v, err.Error()}, nil
		}
		return reflectResponse{v, ""}, nil
	}
}

func decodeReflectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var buf strings.Builder
	if 0 < len(r.Header) {
		for key, val := range r.Header {
			buf.WriteString(fmt.Sprintf("[%11s]==%s\n", key, val))
		}
	}
	n, err := io.Copy(&buf, r.Body)
	if nil != err {
		xLog.Printf("NewDecoder read %d bytes but failed because %s", n, err.Error())
	}
	return reflectRequest{S: buf.String()}, nil
}
