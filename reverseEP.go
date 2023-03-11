package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"reflectsvc/misc"
)

// For each method, we define request and response structs
type reverseRequest struct {
	S string `json:"s"`
}

type reverseResponse struct {
	S   string `json:"s"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

func makeReverseEndpoint(svc SimpleService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(reverseRequest)
		v, err := svc.Reverse(req.S)
		if err != nil {
			return reverseResponse{v, err.Error()}, nil
		}
		return reverseResponse{v, ""}, nil
	}
}

func decodeReverseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer misc.DeferError(xLogBuffer.Flush)
	var request reverseRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if nil != err {
		xLog.Printf("NewDecoder failed because %s", err.Error())
	}
	return request, nil
}
