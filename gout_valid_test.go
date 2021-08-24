package gout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testValid struct {
	Val string `valid:"required"`
}

func Test_Valid(t *testing.T) {
	total := int32(1)
	router := setupMethod(&total)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	testCases := []string{"bindjson", "bindxml", "bindyaml", "bindheader"}
	for _, c := range testCases {
		val := testValid{}
		g := GET(ts.URL + "/someGet")
		var err error
		switch c {
		case "bindjson":
			err = g.BindJSON(&val).Do()
		case "bindxml":
			err = g.BindXML(&val).Do()
		case "bindyaml":
			err = g.BindYAML(&val).Do()
		case "bindheader":
			err = g.BindHeader(&val).Do()
		}

		//fmt.Printf("-->%v\n", err)
		assert.Error(t, err)
	}
}
