package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"reflectsvc/misc"
	"strings"
)

type XtractaEvents struct {
	XMLName xml.Name     `xml:"events"`
	Text    string       `xml:",chardata"`
	Event   XtractaEvent `xml:"event"`
	Headers http.Header
}

func (x XtractaEvents) String() string {
	var sb strings.Builder
	sb.WriteString(
		fmt.Sprintf("XMLSpace:    XMLName:  %s  Text: %s\n",
			x.XMLName, x.Text))
	sb.WriteString(x.Event.String())
	return sb.String()
}

type XtractaEvent struct {
	Text      string          `xml:",chardata"`
	Sequence  string          `xml:"sequence,attr"`
	Generated string          `xml:"generated"`
	Document  XtractaDocument `xml:"document"`
}

func (x XtractaEvent) String() string {
	var sb strings.Builder
	sb.WriteString(
		fmt.Sprintf("XtractaEvent Text: %s  Sequence %s  Generated %s\n",
			x.Text, x.Sequence, x.Generated))
	sb.WriteString(x.Document.String())
	return sb.String()
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

func (x XtractaDocument) String() string {
	var sb strings.Builder
	sb.WriteString(
		fmt.Sprintf("XtractaDocument Text: %s  Revision: %s  WorkflowID: %s  DocumentID: %s\n"+
			"\tDocumentStatus %s  NumberOfPages %s  ApiDownloadStatus:  %s\n"+
			"\tFreeForm: %s  Classification %s  ClassificationStatus: %s\n"+
			"\tClassificationDesign: %s  DocumentURL: %s",
			x.Text, x.Revision, x.WorkflowID, x.DocumentID, x.DocumentStatus,
			x.NumberOfPages, x.ApiDownloadStatus, x.FreeForm, x.Classification,
			x.ClassificationClass, x.ClassificationDesign, x.DocumentURL))
	if len(x.ImageURL) > 0 {
		sb.WriteString("\nImage URLs:\n")
		for _, iu := range x.ImageURL {
			sb.WriteString("\n\t")
			sb.WriteString(iu)
		}
		sb.WriteRune('\n')
	}
	sb.WriteString(x.FieldData.String())
	return sb.String()
}

type XtractaFieldData struct {
	Text  string         `xml:",chardata" json:"text,omitempty"`
	Field []XtractaField `xml:"field" json:"field,omitempty"`
}

func (x XtractaFieldData) String() string {
	var sb strings.Builder
	sb.WriteString("FieldData Text: ")
	sb.WriteString(x.Text)
	for _, fld := range x.Field {
		sb.WriteString("\nField [")
		sb.WriteString(fld.String())
		sb.WriteRune(']')
	}
	return sb.String()
}

type XtractaField struct {
	Text                      string `xml:",chardata" json:"text,omitempty"`
	FieldID                   string `xml:"field_id"`
	FieldName                 string `xml:"field_name"`
	FieldValue                string `xml:"field_value"`
	FieldExtractionConfidence string `xml:"field_extraction_confidence"`
}

func (x XtractaField) String() string {
	return fmt.Sprintf("Text: %s  FieldID: %s  FieldName: %s  FieldValue: %s   FieldConfidence: %s",
		x.Text, x.FieldID, x.FieldName, x.FieldValue, x.FieldExtractionConfidence)
}

// Json() Convert XtractaEvents data to JSON data
// --fieldNames permits remapping XML field names to new JSON field names.
// --omitEmpty means that XML fields without field values are omitted.
func (x XtractaEvents) Json() string {
	var sb strings.Builder

	sb.WriteRune('{')

	needComma4Json := false
	for _, fld := range x.Event.Document.FieldData.Field {
		if FlagOmitEmpty && !misc.IsStringSet(&fld.FieldValue) {
			continue
		}
		if needComma4Json {
			sb.WriteRune(',')
		} else {
			needComma4Json = true
		}
		sb.WriteRune('"')
		{ // change the field name based on lookup
			cfn, ok := FlagRemapMap[fld.FieldName]
			if ok {
				sb.WriteString(cfn)
			} else {
				sb.WriteString(fld.FieldName)
			}
		}
		sb.WriteString("\":\"")
		sb.WriteString(fld.FieldValue)
		sb.WriteRune('"')
	}
	sb.WriteRune('}')
	return sb.String()
}
