// Package tiff provides a TIFF Encoder.
package tiff

import (
	"bytes"
	"image"

	"github.com/pierrre/imageserver"
	imageserver_processor_native "github.com/pierrre/imageserver/processor/native"
	"golang.org/x/image/tiff"
)

// Encoder encodes an Image to TIFF.
type Encoder struct {
}

var opts = &tiff.Options{
	Compression: tiff.Deflate,
	Predictor:   true,
}

// Encode implements Encoder.
func (e *Encoder) Encode(nim image.Image, params imageserver.Params) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := tiff.Encode(buf, nim, opts)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func init() {
	imageserver_processor_native.RegisterEncoder("tiff", &Encoder{})
}
