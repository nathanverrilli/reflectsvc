package main

import (
	"context"
	"encoding/json"
	"errors"
	httpTransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"reflectsvc/misc"
)

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

func main() {
	var err error

	initLog("reflectsvc.log")
	defer closeLog()
	initFlags()

	svc := simpleService{}

	reflectHandler := httpTransport.NewServer(
		makeReflectEndpoint(svc),
		decodeReflectRequest,
		encodeResponse)

	reverseHandler := httpTransport.NewServer(
		makeReverseEndpoint(svc),
		decodeReverseRequest,
		encodeResponse)

	convertHandler := httpTransport.NewServer(
		makeConvertEndpoint(svc),
		decodeConvertRequest,
		encodeResponse)

	xml2JsonHandler := httpTransport.NewServer(
		makeXml2JsonEndpoint(svc),
		decodeXml2JsonRequest,
		encodeResponse)

	http.Handle("/reverse", reverseHandler)
	http.Handle("/parsifal", convertHandler)
	http.Handle("/reflect", reflectHandler)
	http.Handle("/xml2json", xml2JsonHandler)

	service := ":" + FlagPort
	if !misc.IsStringSet(&FlagCert) || !misc.IsStringSet(&FlagKey) {
		xLog.Printf("reverting to HTTP\n\tCertification file is %s\n\tKey file is %s",
			misc.Ternary(misc.IsStringSet(&FlagCert), FlagCert, "missing (use --certfile to set)"),
			misc.Ternary(misc.IsStringSet(&FlagKey), FlagKey, "missing (use --keyfile to set)"))
		flushLog()
		err = http.ListenAndServe(service, nil)
	} else {
		xLog.Printf("using HTTPS\n\tCertification file is %s\n\tKey file is %s",
			FlagCert, FlagKey)
		flushLog()
		err = http.ListenAndServeTLS(service, FlagCert, FlagKey, nil)
	}

	if nil != err {
		xLog.Printf("http listener service failed because %s", err.Error())
		myFatal()
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	defer flushLog()
	return json.NewEncoder(w).Encode(response)
}
