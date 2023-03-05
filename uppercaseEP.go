package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"reflectsvc/misc"
)

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer misc.DeferError(xLogBuffer.Flush)
	var request uppercaseRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		xLog.Printf("NewDecoder failed because %s", err.Error())
	}
	return request, nil
}
