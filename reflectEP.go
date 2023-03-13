package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
	"strings"
	"time"
)

// For each method, we define request and response structs
type reflectRequest struct {
	S string `json:"s"`
}

type reflectResponse struct {
	S   string `json:"s"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

func makeReflectEndpoint(svc SimpleService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(reflectRequest)
		v, err := svc.Reflect(req.S)
		if err != nil {
			return reflectResponse{v, err.Error()}, nil
		}
		return reflectResponse{v, ""}, nil
	}
}

func iMax(i int, j int) (k int) {
	if i >= j {
		return i
	}
	return j
}

func decodeReflectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var buf strings.Builder
	if 0 < len(r.Header) {
		maxKeyLen := 0
		maxValLen := 0
		for key, list := range r.Header {
			maxKeyLen = iMax(maxKeyLen, len(key))
			for _, val := range list {
				maxValLen = iMax(maxValLen, len(val))
			}
		}
		format := fmt.Sprintf("[%%%ds]==[%%-%ds]\n",
			maxKeyLen, maxValLen)
		for key, list := range r.Header {
			for _, val := range list {
				buf.WriteString(fmt.Sprintf(format, key, val))
			}
		}
	}
	buf.WriteRune('\n')
	body, err := io.ReadAll(r.Body)
	if nil != err {
		xLog.Printf("huh? io.ReadAll failed on request body because %s", err.Error())
	}
	buf.Write(body)
	buf.WriteRune('\n')
	buf.WriteString(time.Now().UTC().Format(time.RFC1123))
	return reflectRequest{S: buf.String()}, nil
}
