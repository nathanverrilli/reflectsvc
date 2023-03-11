package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"reflectsvc/misc"
)

type ParsifalResponse struct {
	Success string `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ParsifalRequest struct {
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

func decodeParsifalRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer misc.DeferError(xLogBuffer.Flush)
	var request ParsifalRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		xLog.Printf("NewDecoder failed on (request ParsifalRequest) because %s", err.Error())
	}
	return request, nil
}

func makeParsifalEndpoint(svc SimpleService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ParsifalRequest)
		_, err := svc.ParsifalUpload(req)
		if err != nil {
			if err != nil {
				return ParsifalResponse{"FAILURE", err.Error()}, nil
			}
		}
		return ParsifalResponse{"Success", ""}, nil
	}
}
