package main

import (
	"bytes"
	"io"
	"reflectsvc/misc"
)

const SEP = "/* ************************** */"

// SimpleService provides operations on strings.
type SimpleService interface {
	Reverse(string) (string, error)
	Reflect(request reflectRequest) reflectResponse
	Convert(request ConvertRequest) (string, error)
	Xml2Json(request xml2JsonRequest) x2jProxyData
	Validate(request validateRequest) validateRequest
	// Success(string) string
}

// simpleService is a concrete implementation of SimpleService
type simpleService struct{}

/*
func (simpleService) Success(s string) string {
	return " "
}
*/

func (simpleService) Validate(v validateRequest) (vr validateRequest) {
	return v
}

func (simpleService) Xml2Json(req xml2JsonRequest) (xjProxy x2jProxyData) {
	if FlagDebug {
		xLog.Printf("enter Xml2Json send request %s", req.MagicInternalGuid)
	}
	xjProxy.Code = 500
	xjProxy.Status = "500 ERROR"
	xjProxy.Body = nil

	buf := bytes.NewBufferString(req.Json())
	rsp, err := x2jProxy(req.Headers, buf)

	if nil != err {
		xLog.Printf("could not proxy json request to %s\n with data\n%s\n because %s",
			FlagDest, req.Json(), err.Error())
		if nil != rsp {
			xLog.Printf("response: %v", rsp)
			xjProxy.Code = rsp.StatusCode
			xjProxy.Status = rsp.Status
		} else {
			xjProxy.Status = "No response from remote server"
			return xjProxy
		}
	}
	defer misc.DeferError(rsp.Body.Close)
	xjProxy.Body, err = io.ReadAll(rsp.Body)
	if nil != err {
		xLog.Printf("json request to %s with data\n%s\n"+
			"\tcould not read response body because %s",
			FlagDest, req.Json(), err.Error())
		xjProxy.Status = "failure"
		xjProxy.Code = 501
	}
	xjProxy.Status = rsp.Status
	xjProxy.Code = rsp.StatusCode

	if FlagDebug || FlagVerbose {
		xLog.Printf("\n%s\n%s\n%s", SEP, string(xjProxy.Body), SEP)
	}

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
