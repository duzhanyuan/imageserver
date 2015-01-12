// Package _test provides utilities for Encoder testing.
package _test

import (
	"bytes"
	"image"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_processor_native "github.com/pierrre/imageserver/processor/native"
)

func TestEncoder(t *testing.T, encoder imageserver_processor_native.Encoder, expectedFormat string) {
	TestEncoderParams(t, encoder, imageserver.Params{}, expectedFormat)
}

func TestEncoderParams(t *testing.T, encoder imageserver_processor_native.Encoder, params imageserver.Params, expectedFormat string) {
	im := NewImage()
	data, err := encoder.Encode(im, params)
	if err != nil {
		t.Fatal(err)
	}
	_, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if format != expectedFormat {
		t.Fatalf("unexpected format: %s", format)
	}
}

func NewImage() image.Image {
	return image.NewRGBA(image.Rect(0, 0, 100, 100))
}
