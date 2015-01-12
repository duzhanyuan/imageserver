package jpeg

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_processor_native "github.com/pierrre/imageserver/processor/native"
	imageserver_processor_native_encoder_test "github.com/pierrre/imageserver/processor/native/encoder/_test"
)

var _ imageserver_processor_native.Encoder = &Encoder{}

func TestEncoder(t *testing.T) {
	testEncoder(t, &Encoder{})
}

func TestEncoderDefaultQuality(t *testing.T) {
	encoder := &Encoder{
		DefaultQuality: 90,
	}
	testEncoder(t, encoder)
}

func TestEncoderQuality(t *testing.T) {
	params := imageserver.Params{
		"quality": 90,
	}
	testEncoderParams(t, &Encoder{}, params)
}

func TestEncoderErrorQuality(t *testing.T) {
	im := imageserver_processor_native_encoder_test.NewImage()
	encoder := &Encoder{}
	for _, quality := range []interface{}{"foo", -1, 101} {
		_, err := encoder.Encode(im, imageserver.Params{"quality": quality})
		if err == nil {
			t.Fatal("no error")
		}
		errParam, ok := err.(*imageserver.ParamError)
		if !ok {
			t.Fatalf("wrong error type: %T", err)
		}
		if errParam.Param != "quality" {
			t.Fatalf("wrong param: %s", errParam.Param)
		}
	}
}

func testEncoder(t *testing.T, encoder *Encoder) {
	imageserver_processor_native_encoder_test.TestEncoder(t, encoder, "jpeg")
}

func testEncoderParams(t *testing.T, encoder *Encoder, params imageserver.Params) {
	imageserver_processor_native_encoder_test.TestEncoderParams(t, encoder, params, "jpeg")
}
