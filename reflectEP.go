package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
)

// For each method, we define request and response structs
type reflectRequest struct {
	Body []byte
}

type reflectResponse reflectRequest

func makeReflectEndpoint(svc SimpleService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(reflectRequest)
		return svc.Reflect(req), nil
	}
}

func decodeReflectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var v reflectRequest

	block, err := io.ReadAll(r.Body)
	if nil == err {
		v.Body = block
	}
	return v, err
}

func encodeReflectResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	v, ok := response.(reflectResponse)
	if !ok {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
		_, _ = w.Write(v.Body)
	}
	return nil
}
