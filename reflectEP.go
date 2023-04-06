package main

import (
	"bytes"
	"context"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
	"os"
	"reflectsvc/misc"
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

	if FlagDebug {
		body, _ := io.ReadAll(r.Body)
		xf, _ := os.OpenFile("lastRequestDebug.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		defer misc.DeferError(xf.Close)
		_, _ = xf.Write(body)
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

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
