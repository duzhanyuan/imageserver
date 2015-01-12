package main

import (
	"net/http"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_parser_nfntresize "github.com/pierrre/imageserver/http/parser/nfntresize"
	imageserver_processor "github.com/pierrre/imageserver/processor"
	imageserver_processor_native "github.com/pierrre/imageserver/processor/native"
	_ "github.com/pierrre/imageserver/processor/native/encoder/bmp"
	_ "github.com/pierrre/imageserver/processor/native/encoder/gif"
	_ "github.com/pierrre/imageserver/processor/native/encoder/jpeg"
	_ "github.com/pierrre/imageserver/processor/native/encoder/png"
	_ "github.com/pierrre/imageserver/processor/native/encoder/tiff"
	imageserver_processor_native_nfntresize "github.com/pierrre/imageserver/processor/native/nfntresize"
	imageserver_provider "github.com/pierrre/imageserver/provider"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

func main() {
	server := imageserver.Server(&imageserver_provider.Server{
		Provider: imageserver_testdata.Provider,
	})
	server = &imageserver_processor.Server{
		Server: server,
		Processor: &imageserver_processor_native.Processor{
			Processor: &imageserver_processor_native_nfntresize.Processor{},
		},
	}

	handler := &imageserver_http.Handler{
		Parser: &imageserver_http.ListParser{
			&imageserver_http.SourceParser{},
			&imageserver_http_parser_nfntresize.Parser{},
			&imageserver_http.FormatParser{},
			&imageserver_http.QualityParser{},
		},
		Server: server,
		ErrorFunc: func(err error, request *http.Request) {
			println(err.Error())
		},
	}

	http.Handle("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
