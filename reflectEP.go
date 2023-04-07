package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
	"os"
	"reflectsvc/misc"
	"sync"
	"time"
)

// For each method, we define request and response structs
type reflectRequest struct {
	Body []byte
}

type reflectResponse reflectRequest

func makeReflectEndpoint(svc SimpleService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(reflectRequest)
		return svc.Reflect(req), nil
	}
}

var reflectDebugCount = 0
var reflectSync sync.Mutex

func decodeReflectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var v reflectRequest

	if FlagDebug {
		body, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()
		reflectSync.Lock()
		fn := fmt.Sprintf("%s_rfldbg%03d.log",
			time.Now().UTC().Format(misc.DATE_POG),
			reflectDebugCount)
		reflectDebugCount++
		reflectSync.Unlock()
		xLog.Printf("enter decodeReflectRequest -- saving request as %s\n", fn)
		xf, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = xf.Write(body)
		_ = xf.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	block, err := io.ReadAll(r.Body)

	if nil == err {
		v.Body = block
	}
	_ = r.Body.Close()
	return v, err
}

func encodeReflectResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	v, ok := response.(reflectResponse)
	if !ok {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
		_, _ = w.Write(v.Body)
	}
	return nil
}
