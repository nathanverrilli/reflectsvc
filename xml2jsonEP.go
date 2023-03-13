package main

import (
	"context"
	"encoding/xml"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
)

// For each method, we define request and response structs
type xml2JsonResponse struct {
	Success string `json:"success"`
	Error   string `json:"error,omitempty"`
}

type xml2JsonRequest XtractaEvents

func (pr xml2JsonRequest) String() string {
	return XtractaEvents(pr).String()
}

func (pr xml2JsonRequest) Json() string {
	return XtractaEvents(pr).Json()
}

func makeXml2JsonEndpoint(svc SimpleService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(xml2JsonRequest)
		v, err := svc.xml2Json(req)
		if err != nil {
			return xml2JsonResponse{v, err.Error()}, nil
		}
		return xml2JsonResponse{v, ""}, nil
	}
}

func decodeXml2JsonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer flushLog()
	var req xml2JsonRequest
	body, err := io.ReadAll(r.Body)
	if nil != err {
		xLog.Printf("io.ReadAll failed on decodeXml2JsonRequest because %s", err.Error())
		return nil, err
	}
	err = xml.Unmarshal(body, &req)
	if nil != err {
		xLog.Printf("xml.Unmarshal failed because %s", err.Error())
		return nil, err
	}
	// xLog.Print(req.String())
	return req, nil
}
