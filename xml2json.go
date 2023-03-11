package main

import (
	"bytes"
	"encoding/xml"
	"io"
)

var bufMemAlloc = 2048

type Events struct {
	XMLName xml.Name `xml:"events" json:"events,omitempty"`
	Text    string   `xml:",chardata" json:"text,omitempty"`
	Event   struct {
		Text      string `xml:",chardata" json:"text,omitempty"`
		Sequence  string `xml:"sequence,attr" json:"sequence,omitempty"`
		Generated string `xml:"generated"`
		Document  struct {
			Text                 string   `xml:",chardata" json:"text,omitempty"`
			Revision             string   `xml:"revision,attr" json:"revision,omitempty"`
			WorkflowID           string   `xml:"workflow_id"`
			DocumentID           string   `xml:"document_id"`
			DocumentStatus       string   `xml:"document_status"`
			NumberOfPages        string   `xml:"number_of_pages"`
			ApiDownloadStatus    string   `xml:"api_download_status"`
			FreeForm             string   `xml:"free_form"`
			Classification       string   `xml:"classification"`
			ClassificationClass  string   `xml:"classification_class"`
			ClassificationDesign string   `xml:"classification_design"`
			DocumentURL          string   `xml:"document_url"`
			ImageURL             []string `xml:"image_url"`
			FieldData            struct {
				Text  string `xml:",chardata" json:"text,omitempty"`
				Field []struct {
					Text                      string `xml:",chardata" json:"text,omitempty"`
					FieldID                   string `xml:"field_id"`
					FieldName                 string `xml:"field_name"`
					FieldValue                string `xml:"field_value"`
					FieldExtractionConfidence string `xml:"field_extraction_confidence"`
				} `xml:"field" json:"field,omitempty"`
			} `xml:"field_data" json:"field_data,omitempty"`
		} `xml:"document" json:"document,omitempty"`
	} `xml:"event" json:"event,omitempty"`
}

func x2j(xmlBuffer []byte) []byte {
	var e Events

	err := xml.Unmarshal(xmlBuffer, &e)
	if nil != err {
		xLog.Printf("Unmarshal of %s failed because %s",
			string(xmlBuffer), err.Error())
	}
	b, err := io.ReadAll(xmloutput2jsonoutput(&e))
	if nil != err {
		xLog.Printf("io.Readall failed because %s", err.Error())
		b = nil
	}
	return b
}

func xmloutput2jsonoutput(xmlEvents *Events) (buf bytes.Buffer) {
	buf.Reset()
	buf.Grow(bufMemAlloc)
	buf.WriteRune('{')
	for ix, val := range xmlEvents.Event.Document.FieldData.Field {
		if 0 != ix { // JSON does not allow a trailing comma
			buf.WriteRune(',')
		}
		buf.WriteRune('"')
		buf.WriteString(val.FieldName)
		buf.WriteString("\":\"")
		buf.WriteString(val.FieldValue)
		buf.WriteRune('"')
	}
	buf.WriteRune('}')
	bufMemAlloc = 1024 * (1 + (buf.Len() % 1024))
	return buf
}
