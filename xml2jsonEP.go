package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
	"time"
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

func x2j_proxy(jsonReader io.Reader) (*http.Response, error) {
	var tr *http.Transport
	var standardRequestHeaders = map[string]string{
		"Accept":         "application/json",
		"Accept-Charset": "utf-8",
		"User-Agent":     "go",
		"DNT":            "1",
	}

	if FlagDestInsecure {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		tr = &http.Transport{}
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFunc()

	hReq, err := http.NewRequestWithContext(ctx, http.MethodPost, FlagDest, jsonReader)
	if nil != err {
		logPrintf("huh? Could not create an httpRequest because %s", err.Error())
		return nil, err
	}
	hReq.Header.Set("Content-Type", "application/json")
	hReq.Header.Set("Accept", "application/json")
	for key, val := range standardRequestHeaders {
		hReq.Header.Set(key, val)
	}
	for ix := range FlagHeaderKey {
		hReq.Header.Set(FlagHeaderKey[ix], FlagHeaderValue[ix])
	}

	httpClient := &http.Client{
		Transport: tr,
	}

	return httpClient.Do(hReq)

}
