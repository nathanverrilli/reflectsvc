package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"reflectsvc/misc"
	"strings"
)

var vr validateRequest

func makeValidateEndpoint(svc SimpleService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(validateRequest)
		return svc.Validate(req), nil
	}
}

func encodeValidateResponse(_ context.Context, writer http.ResponseWriter, request interface{}) error {

	ep, err := json.Marshal(request.(validateRequest))
	if nil != err {
		xLog.Printf("failed to marshal JSON data because %s", err.Error())
		return err
	}
	if FlagDebug {
		xLog.Printf("\n\t/*** *************** ***/n%sn\t/*** *************** ***/n",
			string(ep))
	}
	_, err = writer.Write(ep)
	return err

}

func decodeValidateRequest(_ context.Context, req *http.Request) (interface{}, error) {
	contentType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil || !strings.HasPrefix(contentType, "multipart/") {
		xLog.Printf("expected a multipart message error %s", err.Error())
	}
	if FlagDebug && FlagVerbose {
		debugMapStringString(params)
	}
	mr := multipart.NewReader(req.Body, params["boundary"])
	defer misc.DeferError(req.Body.Close)
	for {
		part, err := mr.NextPart()
		if io.EOF == err {
			break
		}
		if nil != err {
			xLog.Printf("could not read NextPart because %s", err.Error())
		}
		processPart(part)
	}
	return vr, nil
}

func processPart(part *multipart.Part) {
	defer misc.DeferError(part.Close)
	fb, err := io.ReadAll(part)
	if nil != err {
		xLog.Printf("could not read part because %s", err.Error())
		myFatal(1)
	}
	switch part.Header.Get("Content-ID") {
	case "metadata":
		err = json.Unmarshal(fb, &vr)
		if nil != err {
			xLog.Printf("failed to unmarshal JSON data because %s", err.Error())
			myFatal(1)
		}
		xLog.Printf("\n\t***JSON Decode Response***\n%+v\n\t**************************\n",
			vr)
	case "media":
		fn := part.Header.Get("Content-Filename")
		if FlagDebug || FlagVerbose {
			xLog.Printf("filesize == %10d name == %s",
				len(fb), fn)
		}
		fo, err := os.Create("./output/" + fn)
		if nil != err {
			xLog.Printf("failed to create file %s because %s", fn, err.Error())
			myFatal(1)
		}
		defer misc.DeferError(fo.Close)
		_, err = fo.Write(fb)
		if nil != err {
			xLog.Printf("failed to write file %s because %s", fn, err.Error())
			myFatal(1)
		}
	}
}
