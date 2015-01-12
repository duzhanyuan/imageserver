package gif

import (
	"testing"

	imageserver_processor_native "github.com/pierrre/imageserver/processor/native"
	imageserver_processor_native_encoder_test "github.com/pierrre/imageserver/processor/native/encoder/_test"
)

var _ imageserver_processor_native.Encoder = &Encoder{}

func TestEncoder(t *testing.T) {
	imageserver_processor_native_encoder_test.TestEncoder(t, &Encoder{}, "gif")
}
