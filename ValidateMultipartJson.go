package main

import (
	"fmt"
	"reflectsvc/misc"
	"strings"
)

type outOfBandData struct {
	Source string   `json:"source"`
	Data   []string `json:"data"`
}

func (o *outOfBandData) String() string {
	var sb strings.Builder
	sb.WriteString("oob data: source [ ")
	if misc.IsStringSet(&o.Source) {
		sb.WriteString(o.Source)
	}
	sb.WriteString(" ] data [ ")
	for _, s := range o.Data {
		sb.WriteRune('"')
		if misc.IsStringSet(&s) {
			sb.WriteString(s)
		}
		sb.WriteRune('"')
	}
	sb.WriteString(" ]\n")
	return sb.String()
}

type duple struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

func (d *duple) String() string {
	return fmt.Sprintf("[%s][%s]", d.Key, d.Val)
}

type imageData struct {
	Height    int64   `json:"height"`
	Width     int64   `json:"width"`
	BPP       int     `json:"bitsPerPixel"`
	ImageMeta []duple `json:"imageMeta"`
}

func (id *imageData) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Height: %d  Width: %d  BPP: %d\n\tImageMeta:"))
	for _, d := range id.ImageMeta {
		sb.WriteRune(' ')
		sb.WriteString(d.String())
	}
	sb.WriteRune('\n')
	return sb.String()
}

type validateRequest struct {
	SpokeId   string          `json:"spokeId"`
	RequestID string          `json:"requestId"`
	OutOfBand []outOfBandData `json:"outOfBand"`
	Images    []imageData     `json:"imageData"`
}

func (v *validateRequest) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SpokeId: %s\tRequestId: %s\n", v.SpokeId, v.RequestID))
	for _, ob := range v.OutOfBand {
		sb.WriteRune('\t')
		sb.WriteString(ob.String())
		sb.WriteRune('\n')
	}
	for _, id := range v.Images {
		sb.WriteRune('\t')
		sb.WriteString(id.String())
		sb.WriteRune('\n')
	}
	return sb.String()
}
