// Package gif provides a GIF Encoder.
package gif

import (
	"bytes"
	"image"
	"image/gif"

	"github.com/pierrre/imageserver"
	imageserver_processor_native "github.com/pierrre/imageserver/processor/native"
)

// Encoder encodes an Image to GIF.
type Encoder struct {
}

// Encode implements Encoder.
func (e *Encoder) Encode(nim image.Image, params imageserver.Params) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := gif.Encode(buf, nim, &gif.Options{})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func init() {
	imageserver_processor_native.RegisterEncoder("gif", &Encoder{})
}
