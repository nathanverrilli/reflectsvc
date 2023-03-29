package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const SEP = "/* ************************** */"

// SimpleService provides operations on strings.
type SimpleService interface {
	Reverse(string) (string, error)
	Reflect(request reflectRequest) reflectResponse
	Convert(request ConvertRequest) (string, error)
	Xml2Json(request xml2JsonRequest) x2j_ProxyData
	Validate(request validateRequest) validateRequest
}

// simpleService is a concrete implementation of SimpleService
type simpleService struct{}

func (simpleService) Validate(v validateRequest) (vr validateRequest) {
	return v
}

func (simpleService) Xml2Json(req xml2JsonRequest) (xjProxy x2j_ProxyData) {

	xjProxy.Code = 500
	xjProxy.Status = "500 ERROR"
	xjProxy.Body = nil

	if FlagDebug {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("hello from xml2JSon\n\tremote endpoint: %s\n",
			FlagDest))
		if len(FlagHeaderKey) > 0 {
			sb.WriteString("call with these configured headers:\n")
			for ix := range FlagHeaderKey {
				sb.WriteString(fmt.Sprintf("\t[%2d]  [%s] == [%s]\n",
					ix, FlagHeaderKey[ix], FlagHeaderValue[ix]))
			}
		}
		sb.WriteString(fmt.Sprintf("%s\n%s\n%s", SEP, req.Json(), SEP))
		logPrintf(sb.String())
	}
	buf := bytes.NewBufferString(req.Json())
	rsp, err := x2j_proxy(buf)
	if nil != err {
		logPrintf("could not proxy json request to %s\n with data\n%s\n because %s",
			FlagDest, req.Json(), err.Error())
		if nil != rsp {
			logPrintf("response: %v", rsp)
		}
		xjProxy.Code = rsp.StatusCode
		xjProxy.Status = rsp.Status
	}
	xjProxy.Body, err = io.ReadAll(rsp.Body)
	if nil != err {
		logPrintf("json request to %s with data\n%s\n"+
			"\tcould not read response body because %s",
			FlagDest, req.Json(), err.Error())
		xjProxy.Status = "failure"
		xjProxy.Code = 501
	}
	xjProxy.Status = rsp.Status
	xjProxy.Code = rsp.StatusCode
	return xjProxy
}

func (simpleService) Reflect(request reflectRequest) reflectResponse {
	if FlagDebug {
		xLog.Printf("reflecting request:\n\t/* *** */\n%s\n\t/* *** */\n", string(request.Body))
	}
	return reflectResponse{Body: request.Body}
}

func (simpleService) Convert(req ConvertRequest) (string, error) {
	xLog.Printf("\n%s\n%s\n%s\n%s\n", SEP, req.String(), req.Json(), SEP)
	return req.Json(), nil
}

func (simpleService) Reverse(s string) (string, error) {
	var r string
	if "" == s {
		return "", ErrEmpty
	}
	for _, c := range s {
		r = string(c) + r
	}
	xLog.Printf("reversed a string %s to %s", s, r)
	return r, nil
}
