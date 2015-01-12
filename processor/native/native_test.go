package native

import (
	imageserver_processor "github.com/pierrre/imageserver/processor"
)

var _ imageserver_processor.Processor = &Processor{}

var _ ProcessorNative = ProcessorNativeFunc(nil)

var _ Encoder = EncoderFunc(nil)
