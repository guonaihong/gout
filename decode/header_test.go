package decode

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type rspHeader struct {
	Date          string `header:"date"`
	Connection    string `header:"connection"`
	ContentLength string `header:"Content-Length"`
	ContentType   string `header:"Content-Type"`
}

type headerTest struct {
	in   *http.Response
	need interface{}
	got  interface{}
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
				},
			},
			need: &rspHeader{
				Date:          "1234",
				Connection:    "close",
				ContentLength: "1234",
				ContentType:   "text",
			},
			got: &rspHeader{},
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
