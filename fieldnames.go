package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"reflectsvc/misc"
	"strings"
)

type JsonFieldType int

const (
	JsonString JsonFieldType = iota
	JsonInteger
	JsonNumeric
	JsonDate
	JsonBoolean
)

type remapField struct {
	JsonName  string
	XMLName   string
	FieldType JsonFieldType
	OmitEmpty bool
	MustBe    []string
}

// FlagRemapMap is not really an argument flag, but it used similarly.
// This is a conversion of incoming XML field names to outgoing JSON
// field names as part of the xml2json endpoint. These to->from strings
// are held in the file specifed by `--fieldNames <file>`. <file> should
// be a plain unicode file. Lines beginning with `#` are ignored (comments).
// Empty lines are ignored. Field name replacements are specified as
// [`From XML Field Name`][`toJsonFieldName`].
var FlagRemapMap map[string]remapField

func loadFieldTranslations(fn string) (remap map[string]remapField) {
	remap = make(map[string]remapField, 64)
	if !misc.IsStringSet(&fn) {
		return remap
	}
	f, err := os.Open(fn)
	if nil != err {
		xLog.Printf("could not open field translation file %s because %s",
			fn, err.Error())
		return remap
	}
	defer misc.DeferError(f.Close)
	rdr := csv.NewReader(bufio.NewReader(f))
	rdr.Comma = ';'
	rdr.Comment = '#'
	rdr.ReuseRecord = true

	var record []string
	for record, err = rdr.Read(); nil != record && nil == err; record, err = rdr.Read() {
		// if the first row is field designators, ignore them
		if "xmlname" == strings.ToLower(record[0]) {
			continue
		}
		var rm remapField
		rm.XMLName = record[0]
		rm.JsonName = record[1]
		switch strings.ToLower(record[2]) {
		case "string":
			rm.FieldType = JsonString
		case "numeric", "number":
			rm.FieldType = JsonNumeric
		case "decimal", "integer":
			rm.FieldType = JsonInteger
		case "boolean", "bool":
			rm.FieldType = JsonBoolean
		case "date":
			rm.FieldType = JsonDate
		default:
			rm.FieldType = JsonString
			xLog.Printf("Huh? Found an unrecognized JSON field "+
				"type  %s in the field translation file (treating it as a string)", record[2])
		}
		switch strings.ToLower(record[3]) {
		case "true":
			rm.OmitEmpty = true
		case "false":
			rm.OmitEmpty = false
		default:
			rm.OmitEmpty = false
			xLog.Printf("Huh? Found a non-true / non-false type for OmitEmpty "+
				" %s in the field translation file (treating it as false", record[2])
		}
		remap[rm.XMLName] = rm
	}
	if nil != err && io.EOF != err {
		xLog.Printf("got non-EOF error while parsing field translation file: %s", err.Error())
	}

	return remap
}
