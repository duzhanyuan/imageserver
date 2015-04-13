package main

import (
	"image"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/pierrre/imageserver/testdata"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

var testHost string

func testMain(m *testing.M) int {
	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	go http.Serve(listener, newImageHTTPHandler())
	testHost = listener.Addr().String()
	return m.Run()
}

func newTestURL() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   testHost,
	}
}

type testCase struct {
	args               map[string]string
	expectedStatusCode int
	expectedFormat     string
	expectedWidth      int
	expectedHeight     int
}

func TestServer(t *testing.T) {
	for _, tc := range []testCase{
		{
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			args: map[string]string{
				"source": testdata.SmallFileName,
			},
		},
		{
			args: map[string]string{
				"source": testdata.MediumFileName,
			},
		},
		{
			args: map[string]string{
				"source": testdata.LargeFileName,
			},
		},
		{
			args: map[string]string{
				"source": testdata.HugeFileName,
			},
		},
		{
			args: map[string]string{
				"source": testdata.AnimatedFileName,
			},
		},
		{
			args: map[string]string{
				"source": testdata.MediumFileName,
			},
		},
		{
			args: map[string]string{
				"source": testdata.MediumFileName,
				"format": "foobar",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			args: map[string]string{
				"source": testdata.MediumFileName,
				"format": "png",
			},
			expectedFormat: "png",
		},
		{
			args: map[string]string{
				"source": testdata.MediumFileName,
				"format": "gif",
			},
			expectedFormat: "gif",
		},
		{
			args: map[string]string{
				"source":  testdata.MediumFileName,
				"format":  "jpeg",
				"quality": "-10",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			args: map[string]string{
				"source":  testdata.MediumFileName,
				"format":  "jpeg",
				"quality": "50",
			},
			expectedFormat: "jpeg",
		},
		{
			args: map[string]string{
				"source": testdata.MediumFileName,
				"width":  "-100",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			args: map[string]string{
				"source": testdata.MediumFileName,
				"width":  "100",
			},
			expectedWidth: 100,
		},
		{
			args: map[string]string{
				"source": testdata.MediumFileName,
				"height": "-100",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			args: map[string]string{
				"source": testdata.MediumFileName,
				"height": "200",
			},
			expectedHeight: 200,
		},
	} {
		testServerRunTestCase(t, tc)
	}
}

func testServerRunTestCase(t *testing.T, tc testCase) {
	t.Logf("%#v", tc)
	u := newTestURL()
	if tc.args != nil {
		query := make(url.Values)
		for k, v := range tc.args {
			query.Add(k, v)
		}
		u.RawQuery = query.Encode()
	}
	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if tc.expectedStatusCode != 0 && resp.StatusCode != tc.expectedStatusCode {
		t.Fatalf("unexpected http status: %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		if tc.expectedStatusCode != 0 {
			return
		}
		t.Fatalf("http status not OK: %d", resp.StatusCode)
	}
	im, format, err := image.Decode(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if tc.expectedFormat != "" && format != tc.expectedFormat {
		t.Fatalf("unexpected format: %s", format)
	}
	if tc.expectedWidth != 0 && im.Bounds().Dx() != tc.expectedWidth {
		t.Fatalf("unexpected width: %d", im.Bounds().Dx())
	}
	if tc.expectedHeight != 0 && im.Bounds().Dy() != tc.expectedHeight {
		t.Fatalf("unexpected height: %d", im.Bounds().Dy())
	}
}
