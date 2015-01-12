// Package png provides a PNG Encoder.
package png

import (
	"bytes"
	"image"
	"image/png"

	"github.com/pierrre/imageserver"
	imageserver_processor_native "github.com/pierrre/imageserver/processor/native"
)

// Encoder encodes an Image to PNG.
type Encoder struct {
}

var encoder = &png.Encoder{
	CompressionLevel: png.BestCompression,
}

// Encode implements Encoder.
func (e *Encoder) Encode(nim image.Image, params imageserver.Params) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := encoder.Encode(buf, nim)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func init() {
	imageserver_processor_native.RegisterEncoder("png", &Encoder{})
}
