package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflectsvc/misc"
	"strconv"
	"sync"
	"time"
)

/*
	{
	    "data": {
	        "result": "ACCEPTED",
	        "evse_id": "USCPIE6579991*1",
	        "location_id": "USCPIL6579991"
	    },
	    "status_code": 1000,
	    "timestamp": "2024-01-12T23:06:27Z"
	}
*/
type successResponse struct {
	StatusCode int                 `json:"status_code"`
	Timestamp  string              `json:"timestamp"`
	Data       successResponseData `'json:"data"`
}

type successResponseData struct {
	Result string `json:"result"`
}

type successRequest struct {
	Data       successRequestData `json:"data,omitempty"`
	StatusCode int                `json:"status_code,omitempty"`
	Timestamp  time.Time          `json:"timestamp,omitempty"`
	// Header     http.Header        `json:"headers,omitempty"`
}

type successRequestData struct {
	Result     string `json:"result,omitempty"`
	EvseID     string `json:"evse_id,omitempty"`
	LocationID string `json:"location_id,omitempty"`
}

/*
{
  "data": {
    "result": "ACCEPTED",
    "evse_id": "USCPIE6579991*1",
    "location_id": "USCPIL6579991"
  },
  "status_code": 1000,
  "timestamp": "2024-01-12T23:06:27Z"
}

*/

const CPTIME = "2006-0102T150405Z"

func makeSuccessResponse() (r successResponse) {

	return successResponse{
		StatusCode: 1000,
		Timestamp:  time.Now().UTC().Format(CPTIME),
		Data: successResponseData{
			Result: "success",
		},
	}
}

func makeSuccessRequest() (r successRequest) {
	return successRequest{
		Data:       successRequestData{EvseID: "USCPIL1", Result: "success", LocationID: "USCPIL2"},
		StatusCode: 1000,
		Timestamp:  time.Now().UTC(),
		// Header:     nil,
	}
}

func makeSuccessEndpoint(_ SimpleService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return makeSuccessResponse(), nil
	}
}

func encodeSuccessResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	v, ok := response.(successResponse)
	if !ok {
		xLog.Printf("What? Got bad successResponse from makeSuccessEndpoint?")
		v = makeSuccessResponse()
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

var successDebugCount int64
var successLock sync.Mutex

func decodeSuccessRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer misc.DeferError(xLogBuffer.Flush)
	var req successRequest

	if FlagDebug {
		var fn, guid string
		body, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()
		{
			successLock.Lock()
			guid = strconv.FormatInt(rand.Int63(), 36)
			fn = fmt.Sprintf("%s_success%03d.log",
				time.Now().UTC().Format(misc.DATE_POG),
				successDebugCount)
			successDebugCount++
			successLock.Unlock()
		}
		xLog.Printf("enter decodeSuccessRequest -- %s -- saving request as %s",
			guid, fn)
		xf, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		defer misc.DeferError(xf.Close)
		var hostname string
		{
			hostnamearray, ok := r.Header["X-Forwarded-Host"]
			if ok && len(hostnamearray) > 0 && misc.IsStringSet(&hostnamearray[0]) {
				hostname = hostnamearray[0]
			} else {
				hostname = "===X=Forwarded-Host-Header-Absent==="
			}
		}
		_, _ = fmt.Fprintf(xf, "host{path} [%s{%s}]\n", hostname, r.URL.String())
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

	if nil != err {
		xLog.Printf("xml.Unmarshal failed because %s\nbody[ %s ]", err.Error(),
			string(body))
		// return nil, err
		req = makeSuccessRequest()
		req.Data.Result = string(body)
	}
	// req.Header = r.Header
	return req, nil
}
