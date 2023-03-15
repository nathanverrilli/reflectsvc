package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
	"time"
)

type x2j_ProxyData struct {
	Code   int
	Status string
	Body   []byte
}

func (xj x2j_ProxyData) String() string {
	return fmt.Sprintf("status: [%s] status code: [%d]\n---BodyDataStart---\n%s\n---BodyDataEnd\n",
		xj.Status, xj.Code, xj.Body)
}

// For each method, we define request and response structs
type xml2JsonResponse x2j_ProxyData

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
		v := svc.xml2Json(req)
		return xml2JsonResponse(v), nil
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

func x2j_encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	v, ok := response.(xml2JsonResponse)

	if !ok {
		s := fmt.Sprintf("{\"error\":\"%s\"}", v.Status)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(s))

	} else {
		w.WriteHeader(v.Code)
		_, _ = w.Write(v.Body)
	}
	return nil
}
