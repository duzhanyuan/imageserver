// Package nfntresize provides an Encoder that uses nfnt resize library.
package nfntresize

import (
	"fmt"
	"image"

	"github.com/nfnt/resize"
	"github.com/pierrre/imageserver"
)

const (
	globalParam = "nfntresize"
)

// Processor processes an Image with nfnt resize.
type Processor struct {
}

// Process implements Processor.
func (proc *Processor) Process(im image.Image, params imageserver.Params) (image.Image, error) {
	if !params.Has(globalParam) {
		return im, nil
	}
	params, err := params.GetParams(globalParam)
	if err != nil {
		return nil, err
	}
	if params.Empty() {
		return im, nil
	}
	im, err = process(im, params)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = fmt.Sprintf("%s.%s", globalParam, err.Param)
		}
		return nil, err
	}
	return im, nil
}

func process(im image.Image, params imageserver.Params) (image.Image, error) {
	width, err := getDimension("width", params)
	if err != nil {
		return nil, err
	}
	height, err := getDimension("height", params)
	if err != nil {
		return nil, err
	}
	if width == 0 && height == 0 {
		return im, nil
	}

	interp, err := getInterpolation(params)
	if err != nil {
		return nil, err
	}

	mode, err := getModeFunc(params)
	if err != nil {
		return nil, err
	}

	im = mode(width, height, im, interp)
	return im, nil
}

func getDimension(name string, params imageserver.Params) (uint, error) {
	if !params.Has(name) {
		return 0, nil
	}
	dimension, err := params.GetInt(name)
	if err != nil {
		return 0, &imageserver.ParamError{Param: name, Message: err.Error()}
	}
	if dimension < 0 {
		return 0, &imageserver.ParamError{Param: name, Message: "must be greater than or equal to 0"}
	}
	return uint(dimension), nil
}

func getInterpolation(params imageserver.Params) (resize.InterpolationFunction, error) {
	if !params.Has("interpolation") {
		return resize.Lanczos3, nil
	}
	interpolation, err := params.GetString("interpolation")
	if err != nil {
		return 0, &imageserver.ParamError{Param: "interpolation", Message: err.Error()}
	}
	switch interpolation {
	case "nearest_neighbor":
		return resize.NearestNeighbor, nil
	case "bilinear":
		return resize.Bilinear, nil
	case "bicubic":
		return resize.Bicubic, nil
	case "mitchell_netravali":
		return resize.MitchellNetravali, nil
	case "lanczos2":
		return resize.Lanczos2, nil
	case "lanczos3":
		return resize.Lanczos3, nil
	default:
		return 0, &imageserver.ParamError{Param: "interpolation", Message: "invalid value"}
	}
}

type modeFunc func(uint, uint, image.Image, resize.InterpolationFunction) image.Image

func getModeFunc(params imageserver.Params) (modeFunc, error) {
	if !params.Has("mode") {
		return resize.Resize, nil
	}
	mode, err := params.GetString("mode")
	if err != nil {
		return nil, &imageserver.ParamError{Param: "mode", Message: err.Error()}
	}
	switch mode {
	case "resize":
		return resize.Resize, nil
	case "thumbnail":
		return resize.Thumbnail, nil
	default:
		return nil, &imageserver.ParamError{Param: "mode", Message: "invalid value"}
	}
}
