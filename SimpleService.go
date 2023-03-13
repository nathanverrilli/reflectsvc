package main

import (
	"bufio"
	"fmt"
	"os"
	"reflectsvc/misc"
	"time"
)

const SEP = "/* ************************** */"

// SimpleService provides operations on strings.
type SimpleService interface {
	Reverse(string) (string, error)
	Reflect(string) (string, error)
	xml2Json(request xml2JsonRequest) (string, error)
	Convert(request ConvertRequest) (string, error)
}

// simpleService is a concrete implementation of SimpleService
type simpleService struct{}

func (simpleService) xml2Json(request xml2JsonRequest) (string, error) {
	return request.Json(), nil
}

func (simpleService) Reflect(json string) (string, error) {
	return json, nil
}

func (simpleService) Convert(req ConvertRequest) (string, error) {
	xLog.Printf("\n%s\n%s\n%s\n", SEP, req, SEP)
	return req.Json(), nil
}

func (simpleService) Parsifal(request ConvertRequest) (string, error) {
	fn := "parsifal." + time.Now().UTC().Format(misc.DATE_POG) + ".log.txt"
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		xLog.Printf("error opening file %s: %s", fn, err.Error())
		return "", err
	}
	defer misc.DeferError(f.Close)
	b := bufio.NewWriter(f)
	defer misc.DeferError(b.Flush)
	_, _ = fmt.Fprintf(b, "%+v\n%s\n%s",
		request, SEP, request.Json())

	if FlagDebug || FlagVerbose {
		xLog.Print(request.String())
		xLog.Print(request.Json())
	}

	return "success", nil
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
