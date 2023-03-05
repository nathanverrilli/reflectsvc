package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"reflectsvc/misc"
)

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer misc.DeferError(xLogBuffer.Flush)
	var request countRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		xLog.Printf("NewDecoder failed because %s", err.Error())
	}
	return request, nil
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}
