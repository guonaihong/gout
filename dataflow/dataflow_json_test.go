package dataflow

import (
	"fmt"
	"testing"
)

func Test_SetJSONEscape(t *testing.T) {
	ts := createGeneral("")
	fmt.Println("url", ts.URL)
	POST(ts.URL).Debug(true).SetJSONNotEscape(map[string]any{"url": "http://www.com?a=b&c=d"}).Do()
	POST(ts.URL).Debug(true).SetJSON(map[string]any{"url": "http://www.com?a=b&c=d"}).Do()
}
