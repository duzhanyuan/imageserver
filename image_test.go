package imageserver_test

import (
	"encoding"
	"encoding/binary"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ encoding.BinaryMarshaler = new(Image)
var _ encoding.BinaryUnmarshaler = new(Image)

func TestImageMarshal(t *testing.T) {
	data, err := testdata.Medium.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	im := new(Image)
	err = im.UnmarshalBinary(data)
	if err != nil {
		t.Fatal(err)
	}
	if !ImageEqual(im, testdata.Medium) {
		t.Fatal("image not equals")
	}
}

func TestImageMarshallErrorFormatMaxLen(t *testing.T) {
	im := &Image{
		Format: strings.Repeat("a", ImageFormatMaxLen+1),
	}
	_, err := im.MarshalBinary()
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestImageMarshallErrorDataMaxLen(t *testing.T) {
	var data []byte
	dataHeader := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	dataHeader.Len = ImageDataMaxLen + 1
	im := &Image{
		Data: data,
	}
	_, err := im.MarshalBinary()
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestImageUnmarshalBinaryErrorEndOfData(t *testing.T) {
	data, err := testdata.Medium.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	index := -1 // Always truncate 1 byte
	for _, offset := range []int{
		4,
		len(testdata.Medium.Format),
		4,
		len(testdata.Medium.Data),
	} {
		index += offset
		errorData := data[0:index]
		im := new(Image)
		err := im.UnmarshalBinary(errorData)
		if err == nil {
			t.Fatal("no error")
		}
		if _, ok := err.(*ImageError); !ok {
			t.Fatalf("unexpected error type: %T", err)
		}
	}
}

func TestImageUnmarshalBinaryErrorFormatMaxLen(t *testing.T) {
	data, err := testdata.Medium.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	formatLenPosition := 0
	binary.LittleEndian.PutUint32(data[formatLenPosition:formatLenPosition+4], uint32(ImageFormatMaxLen+1))
	im := new(Image)
	err = im.UnmarshalBinary(data)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestImageUnmarshalBinaryErrorDataMaxLen(t *testing.T) {
	data, err := testdata.Medium.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	dataLenPosition := 4 + len(testdata.Medium.Format)
	binary.LittleEndian.PutUint32(data[dataLenPosition:dataLenPosition+4], uint32(ImageDataMaxLen+1))
	im := new(Image)
	err = im.UnmarshalBinary(data)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

// TestImageMarshalBugBuffer is a test for a bug with a misused byte buffer pool.
// Successive calls to Image.MarshalBinary() write to the same byte slice.
func TestImageMarshalBugByteBufferPool(t *testing.T) {
	d1, err := testdata.Small.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	d2, err := testdata.Medium.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	im1 := new(Image)
	err = im1.UnmarshalBinary(d1)
	if err != nil {
		t.Fatal(err)
	}
	im2 := new(Image)
	err = im2.UnmarshalBinary(d2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestImageEqual(t *testing.T) {
	type TC struct {
		im1         *Image
		im2         *Image
		equal       bool
		description string
	}
	for _, tc := range []TC{
		{nil, nil, true, "both nil"},
		{nil, testdata.Medium, false, "nil / not nil"},
		{testdata.Medium, nil, false, "not nil / nil"},
		{testdata.Medium, testdata.Medium, true, "same"},
		{testdata.Medium, imageCopy(testdata.Medium), true, "copy"},
		{testdata.Medium, testdata.Animated, false, "different format"},
		{testdata.Medium, testdata.Small, false, "different data"},
	} {
		if ImageEqual(tc.im1, tc.im2) != tc.equal {
			t.Fatalf("invalid result for test: %s", tc.description)
		}
	}
}

func TestImageError(t *testing.T) {
	err := &ImageError{Message: "test"}
	_ = err.Error()
}

func imageCopy(im *Image) *Image {
	value := *im
	return &value
}
