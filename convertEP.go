package main

import (
	"context"
	_ "encoding/json"
	"encoding/xml"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
	"reflectsvc/misc"
)

type ConvertResponse struct {
	Success string `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ConvertRequest XtractaEvents

func (pr ConvertRequest) String() string {
	return XtractaEvents(pr).String()
}

func (pr ConvertRequest) Json() string {
	return XtractaEvents(pr).Json()
}

/*
	type ConvertRequest struct {
		ShipmentType       string `json:"Shipment Type"`
		TransitType        string `json:"Transit Type"`
		TransitMode        string `json:"Transit Mode"`
		BillType           string `json:"Bill Type"`
		Auditor            string `json:"Auditor"`
		StorageOnBill      string `json:"Storage On Bill"`
		Account            string `json:"Account"`
		Division           string `json:"Division"`
		FileNumber         string `json:"File Number"`
		LastName           string `json:"Last Name"`
		FirstName          string `json:"First Name"`
		BookingAgent       string `json:"Booking Agent"`
		LoadDate           string `json:"Load Date"`
		OriginCountry      string `json:"Origin Country"`
		OriginState        string `json:"Origin State"`
		OriginCity         string `json:"Origin City"`
		DestinationCountry string `json:"Destination Country"`
		DestinationState   string `json:"Destination City"`
		BillingSP          string `json:"Billing SP"`
		SPBillNumber       string `json:"SP Bill Number"`
		SPBillDate         string `json:"SP Bill Date"`
		SPBillAmount       string `json:"SP Bill Amount"`
		Currency           string `json:"Currency"`
		BillReceivedDate   string `json:"Bill received Date"` // "Bill ReceivedDate
		CostCenter         string `json:"Cost Center"`
		CostEENumber       string `json:"cost EE Number"` // "Cost EE Number
		ClientID           string `json:"ClientID"`
		ModelCD            string `json:"ModelCD"`
		SPBillGrossAmount  string `json:"SPBill Gross amount"` // "SP Bill Gross Amount"
	}
*/

func decodeConvertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer misc.DeferError(xLogBuffer.Flush)
	var request ConvertRequest
	body, err := io.ReadAll(r.Body)
	if nil != err {
		xLog.Printf("io.ReadAll failed on ConvertRequest because %s", err.Error())
		return nil, err
	}
	err = xml.Unmarshal(body, &request)
	if nil != err {
		xLog.Printf("xml.Unmarshal failed because %s", err.Error())
		return nil, err
	}
	// xLog.Print(request.String())
	return request, nil
}

func makeConvertEndpoint(svc SimpleService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ConvertRequest)
		_, err := svc.Convert(req)
		if err != nil {
			if err != nil {
				return ConvertResponse{"FAILURE", err.Error()}, nil
			}
		}
		return ConvertResponse{"Success", ""}, nil
	}
}
