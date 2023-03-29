package main

import (
	"context"
	"encoding/json"
	"errors"
	httpTransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"reflectsvc/misc"
	"time"
)

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

func main() {
	var err error

	initLog("reflectsvc.log")
	defer closeLog()
	initFlags()
	// setup for specific services

	svc := simpleService{}

	validateHandler := httpTransport.NewServer(
		makeValidateEndpoint(svc),
		decodeValidateRequest,
		encodeResponse)

	reflectHandler := httpTransport.NewServer(
		makeReflectEndpoint(svc),
		decodeReflectRequest,
		encodeReflectResponse)

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
		x2j_encodeResponse)

	http.Handle("/reverse", reverseHandler)
	http.Handle("/parsifal", convertHandler)
	http.Handle("/convert", convertHandler)
	http.Handle("/reflect", reflectHandler)
	http.Handle("/validate", validateHandler)
	http.Handle("/xml2json", xml2JsonHandler)

	service := ":" + FlagPort

	srv := http.Server{
		Addr:                         service,
		Handler:                      nil,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  2 * time.Minute,
		ReadHeaderTimeout:            2 * time.Minute,
		WriteTimeout:                 2 * time.Minute,
		IdleTimeout:                  2 * time.Minute,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}

	if !misc.IsStringSet(&FlagCert) || !misc.IsStringSet(&FlagKey) {
		xLog.Printf("reverting to HTTP\n\tCertification file is %s\n\tKey file is %s",
			misc.Ternary(misc.IsStringSet(&FlagCert), FlagCert, "missing (use --certfile to set)"),
			misc.Ternary(misc.IsStringSet(&FlagKey), FlagKey, "missing (use --keyfile to set)"))
		flushLog()
		// err = http.ListenAndServe(service, nil)
		err = srv.ListenAndServe()
	} else {
		xLog.Printf("using HTTPS\n\tCertification file is %s\n\tKey file is %s",
			FlagCert, FlagKey)
		flushLog()
		//err = http.ListenAndServeTLS(service, FlagCert, FlagKey, nil)
		err = srv.ListenAndServeTLS(FlagCert, FlagKey)
	}

	if nil != err {
		xLog.Printf("http listener service failed because %s", err.Error())
		myFatal()
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
