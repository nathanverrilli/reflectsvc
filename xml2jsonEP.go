package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
	"strings"
	"time"
)

type x2jProxyData struct {
	Code   int
	Status string
	Body   []byte
}

func (xj x2jProxyData) String() string {
	return fmt.Sprintf("status: [%s] status code: [%d]\n---BodyDataStart---\n%s\n---BodyDataEnd\n",
		xj.Status, xj.Code, xj.Body)
}

// For each method, we define request and response structs
type xml2JsonResponse x2jProxyData

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
		v := svc.Xml2Json(req)
		return xml2JsonResponse(v), nil
	}
}

func decodeXml2JsonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if FlagDebug {
		xLog.Printf("enter decodeXml2JsonRequest")
	}
	var req xml2JsonRequest
	body, err := io.ReadAll(r.Body)
	if nil != err {
		xLog.Printf("io.ReadAll failed on decodeXml2JsonRequest because %s", err.Error())
		return nil, err
	}
	err = xml.Unmarshal(body, &req)
	req.Headers = r.Header
	if nil != err {
		xLog.Printf("xml.Unmarshal failed because %s", err.Error())
		return nil, err
	}

	return req, nil
}

var standardRequestHeaders = map[string]string{
	"Accept-Charset": "utf-8",
	"DNT":            "1",
}

var proxiedHeaders = []string{"Authorization", "User-Agent", "Ocp-Apim-Subscription-Key"}

func x2jProxy(header http.Header, jsonReader io.Reader) (*http.Response, error) {
	var tr *http.Transport

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
		xLog.Printf("huh? Could not create an httpRequest because %s", err.Error())
		return nil, err
	}
	hReq.Header.Set("Content-Type", "application/json")
	hReq.Header.Set("Accept", "application/json")
	for key, val := range standardRequestHeaders {
		hReq.Header.Set(key, val)
	}

	for _, ph := range proxiedHeaders {
		values, ok := header[ph]
		if ok {
			hReq.Header[ph] = values
		}
	}

	if FlagDebug {
		logHeaders(hReq.Header)
	}

	httpClient := &http.Client{
		Transport: tr,
	}

	return httpClient.Do(hReq)
}

func logHeaders(h http.Header) {
	var sb strings.Builder
	for headerName, headerValues := range h {
		sb.WriteRune('\t')
		sb.WriteString(headerName)
		sb.WriteString(" : ")
		for _, value := range headerValues {
			sb.WriteRune('[')
			sb.WriteString(value)
			sb.WriteString("] ")
		}
		sb.WriteRune('\n')
	}
	xLog.Printf("headers for proxied request\n%s", sb.String())
}

func x2jEncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if FlagDebug {
		xLog.Printf("enter x2jEncodeResponse")
	}
	v, ok := response.(xml2JsonResponse)

	if !ok || nil == v.Body || len(v.Body) <= 0 {
		s := fmt.Sprintf("{\"error\":\"%s\"}", v.Status)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(s))

	} else {
		w.WriteHeader(v.Code)
		_, _ = w.Write(v.Body)
	}
	return nil
}
