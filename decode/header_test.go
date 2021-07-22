package decode

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type rspHeader struct {
	Date          string   `header:"date"`
	Connection    string   `header:"connection"`
	ContentLength string   `header:"Content-Length"`
	ContentType   string   `header:"Content-Type"`
	SID           []string `header:"sid"`
	Max           int      `header:"max"`
	Rate          float64  `header:"rate"`
}

type headerTest struct {
	in   *http.Response
	need interface{}
	got  interface{}
}

func Test_Header_Decode_Err(t *testing.T) {
	err := (&headerDecode{}).Decode(nil, nil)
	assert.Error(t, err)
}

func Test_Header_Decode(t *testing.T) {

	tests := []headerTest{
		{
			in: &http.Response{
				Header: http.Header{
					"Date":           []string{"1234"},
					"Connection":     []string{"close"},
					"Content-Length": []string{"1234"},
					"Content-Type":   []string{"text"},
					"Sid":            []string{"sid1", "sid2"},
					"Max":            []string{"1000"},
					"Rate":           []string{"16000"},
				},
			},
			need: &rspHeader{
				Date:          "1234",
				Connection:    "close",
				ContentLength: "1234",
				ContentType:   "text",
				SID:           []string{"sid1", "sid2"},
				Max:           1000,
				Rate:          16000,
			},
			got: &rspHeader{},
		},
		{
			in: &http.Response{
				Header: http.Header{
					"Date":           []string{"1234"},
					"Connection":     []string{"close"},
					"Content-Length": []string{"1234"},
					"Content-Type":   []string{"text"},
					"Sid":            []string{"sid1", "sid2"},
					"Max":            []string{"1000"},
					"Rate":           []string{"16000"},
				},
			},
			need: http.Header{
				"Date":           []string{"1234"},
				"Connection":     []string{"close"},
				"Content-Length": []string{"1234"},
				"Content-Type":   []string{"text"},
				"Sid":            []string{"sid1", "sid2"},
				"Max":            []string{"1000"},
				"Rate":           []string{"16000"},
			},
			got: http.Header{},
		},
		{
			in: &http.Response{
				Header: http.Header{
					"Date":           []string{"1234"},
					"Connection":     []string{"close"},
					"Content-Length": []string{"1234"},
					"Content-Type":   []string{"text"},
					"Sid":            []string{"sid1", "sid2"},
					"Max":            []string{"1000"},
					"Rate":           []string{"16000"},
				},
			},
			need: &http.Header{
				"Date":           []string{"1234"},
				"Connection":     []string{"close"},
				"Content-Length": []string{"1234"},
				"Content-Type":   []string{"text"},
				"Sid":            []string{"sid1", "sid2"},
				"Max":            []string{"1000"},
				"Rate":           []string{"16000"},
			},
			got: &http.Header{},
		},
	}

	for _, v := range tests {
		err := (&headerDecode{}).Decode(v.in, v.got)
		assert.NoError(t, err)
		if err != nil {
			return
		}

		assert.Equal(t, v.got, v.need)
	}
}
