// Package bmp provides a BMP Encoder.
package bmp

import (
	"bytes"
	"image"

	"github.com/pierrre/imageserver"
	imageserver_processor_native "github.com/pierrre/imageserver/processor/native"
	"golang.org/x/image/bmp"
)

// Encoder encodes an Image to BMP.
type Encoder struct {
}

// Encode implements Encoder.
func (e *Encoder) Encode(nim image.Image, params imageserver.Params) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := bmp.Encode(buf, nim)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func init() {
	imageserver_processor_native.RegisterEncoder("bmp", &Encoder{})
}
