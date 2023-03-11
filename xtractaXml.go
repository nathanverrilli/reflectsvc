package main

import "encoding/xml"

type XtractaEvents struct {
	XMLName xml.Name     `xml:"events"`
	Text    string       `xml:",chardata"`
	Event   XtractaEvent `xml:"event"`
}

type XtractaEvent struct {
	Text      string          `xml:",chardata"`
	Sequence  string          `xml:"sequence,attr"`
	Generated string          `xml:"generated"`
	Document  XtractaDocument `xml:"document"`
}

type XtractaDocument struct {
	Text                 string           `xml:",chardata" json:"text,omitempty"`
	Revision             string           `xml:"revision,attr" json:"revision,omitempty"`
	WorkflowID           string           `xml:"workflow_id"`
	DocumentID           string           `xml:"document_id"`
	DocumentStatus       string           `xml:"document_status"`
	NumberOfPages        string           `xml:"number_of_pages"`
	ApiDownloadStatus    string           `xml:"api_download_status"`
	FreeForm             string           `xml:"free_form"`
	Classification       string           `xml:"classification"`
	ClassificationClass  string           `xml:"classification_class"`
	ClassificationDesign string           `xml:"classification_design"`
	DocumentURL          string           `xml:"document_url"`
	ImageURL             []string         `xml:"image_url"`
	FieldData            XtractaFieldData `xml:"field_data"`
}

type XtractaFieldData struct {
	Text  string         `xml:",chardata" json:"text,omitempty"`
	Field []XtractaField `xml:"field" json:"field,omitempty"`
}

type XtractaField struct {
	Text                      string `xml:",chardata" json:"text,omitempty"`
	FieldID                   string `xml:"field_id"`
	FieldName                 string `xml:"field_name"`
	FieldValue                string `xml:"field_value"`
	FieldExtractionConfidence string `xml:"field_extraction_confidence"`
}
