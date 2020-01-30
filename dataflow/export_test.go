package dataflow

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type testCurl struct{}

func (t *testCurl) New(df *DataFlow) interface{} {
	return &testCurl{}
}

func (t *testCurl) LongOption() Curl {
	return t
}

func (t *testCurl) GenAndSend() Curl {
	return t
}

func (t *testCurl) SetOutput(w io.Writer) Curl {
	return t
}

func (t *testCurl) Do() error {
	return nil
}

type testCurlFail struct{}

func (t *testCurlFail) New(df *DataFlow) interface{} {
	return t
}

func Test_Export_curl(t *testing.T) {
	bkcurl, ok := filters["curl"]
	delete(filters, "curl")
	defer func() {
		if ok {
			filters["curl"] = bkcurl
		}
	}()

	// test panic
	for _, v := range []func(){
		func() {
			e := export{}
			e.Curl()
		},
		func() {
			filters["curl"] = &testCurlFail{}
			e := export{}
			e.Curl()
		},
	} {
		assert.Panics(t, v)
	}

	//test ok
	for _, v := range []func(){
		func() {
			filters["curl"] = &testCurl{}
			e := export{}
			e.Curl()
		},
	} {
		assert.NotPanics(t, v)
	}
}
