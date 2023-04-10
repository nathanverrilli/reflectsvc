package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflectsvc/misc"
	"strconv"
	"strings"
	"sync"
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

var xmlDebugCount = 0
var decodeSync sync.Mutex

func decodeXml2JsonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req xml2JsonRequest

	if FlagDebug {
		body, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()
		decodeSync.Lock()
		guid := strconv.FormatInt(rand.Int63(), 36)
		req.MagicInternalGuid = guid
		fn := fmt.Sprintf("%s_xmldbg%03d.log",
			time.Now().UTC().Format(misc.DATE_POG),
			xmlDebugCount)
		xmlDebugCount++
		decodeSync.Unlock()
		xLog.Printf("enter decodeXml2JsonRequest -- %s -- saving request as %s",
			guid, fn)
		xf, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		defer misc.DeferError(xf.Close)
		_, _ = fmt.Fprintf(xf, "request %s\n\t\tHEADERS\n", fn)
		_, _ = xf.Write(debugMapStringArrayString(r.Header))
		_, _ = fmt.Fprintf(xf, "\n\t\tBODY\n")
		_, _ = xf.Write(body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

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
	var responseBody string
	if FlagDebug {
		xLog.Printf("enter x2jEncodeResponse")
	}
	v, ok := response.(xml2JsonResponse)

	if !ok || nil == v.Body || len(v.Body) <= 0 {
		responseBody = fmt.Sprintf("{\"error\":\"%s\"}", v.Status)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(responseBody))
		if nil != err {
			xLog.Printf("could not write header to response because %s", err.Error())
			return err
		}
	} else {
		w.WriteHeader(v.Code)
		//_. err := w.Write(v.Body)
		responseBody = "{\"success\":true}"
		_, err := w.Write([]byte(responseBody))
		if nil != err {
			xLog.Printf("could not write header to response because %s", err.Error())
			return err
		}
	}
	if FlagDebug {

		fn := fmt.Sprintf("%s_xmlrspdbg%03d.log",
			time.Now().UTC().Format(misc.DATE_POG),
			xmlDebugCount)
		xf, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = fmt.Fprintf(xf, "request %s\n", fn)
		_, _ = xf.Write(debugMapStringArrayString(w.Header()))
		_, _ = xf.WriteString("\n")
		_, _ = xf.Write([]byte(responseBody))
		_, _ = xf.WriteString("\n")
		_ = xf.Close()

		xLog.Printf("exiting x2jEncodeResponse")
	}
	return nil
}

func debugMapStringArrayString(m map[string][]string) []byte {
	var sb strings.Builder

	if len(m) <= 0 {
		return []byte("\n no headers \n")
	}

	for key, val := range m {
		sb.WriteString(key)
		sb.WriteString(" [")
		first := true
		for _, val2 := range val {
			if first {
				sb.WriteRune(' ')
				first = false
			} else {
				sb.WriteString(" | ")
			}
			sb.WriteString(val2)
		}
		sb.WriteString(" ]\n")
	}
	return []byte(sb.String())
}
