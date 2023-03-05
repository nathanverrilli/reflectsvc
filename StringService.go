package main

import (
	"bufio"
	"fmt"
	"os"
	"reflectsvc/misc"
	"strings"
	"time"
)

// StringService provides operations on strings.
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
	Reverse(string) (string, error)
	ParsifalUpload(request ParsifalRequest) (string, error)
	Reflect(string) (string, error)
}

// stringService is a concrete implementation of StringService
type stringService struct{}

func (stringService) Reflect(req string) (string, error) {
	xLog.Printf("\n  ************* \n%s\n ************** \n", req)
	return req, nil
}

func (stringService) ParsifalUpload(request ParsifalRequest) (string, error) {
	fn := "parsifal." + time.Now().UTC().Format(misc.DATE_POG) + ".log.txt"
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		xLog.Printf("error opening file %s: %s", fn, err.Error())
		return "", err
	}
	defer misc.DeferError(f.Close)

	b := bufio.NewWriter(f)
	defer misc.DeferError(b.Flush)
	_, _ = fmt.Fprintf(b, "%+v\n", request)
	if FlagDebug {
		xLog.Printf("%+v\n", request)
	}

	return "success", nil
}

func (stringService) Reverse(s string) (string, error) {
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

func (stringService) Uppercase(s string) (string, error) {
	defer misc.DeferError(xLogBuffer.Flush)
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	defer misc.DeferError(xLogBuffer.Flush)
	return len(s)
}
