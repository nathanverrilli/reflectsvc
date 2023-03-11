package main

import (
	"context"
	"encoding/json"
	"errors"
	httpTransport "github.com/go-kit/kit/transport/http"
	"log"
	"net/http"
)

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

func main() {

	initLog("reflectsvc.log")
	defer closeLog()
	initFlags()

	svc := stringService{}

	reflectHandler := httpTransport.NewServer(
		makeReflectEndpoint(svc),
		decodeReflectRequest,
		encodeResponse)

	reverseHandler := httpTransport.NewServer(
		makeReverseEndpoint(svc),
		decodeReverseRequest,
		encodeResponse)

	parsifalHandler := httpTransport.NewServer(
		makeParsifalEndpoint(svc),
		decodeParsifalRequest,
		encodeResponse)

	http.Handle("/reverse", reverseHandler)
	http.Handle("/parsifal", parsifalHandler)
	http.Handle("/reflect", reflectHandler)

	log.Fatal(http.ListenAndServe(":"+FlagPort, nil))
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	defer flushLog()
	return json.NewEncoder(w).Encode(response)
}
