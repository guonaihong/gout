package dataflow

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/guonaihong/gout/debug"
	"github.com/stretchr/testify/assert"
)

func Test_SetJSONNotEscape(t *testing.T) {
	ts := createGeneral("")
	fmt.Println("url", ts.URL)
	var buf bytes.Buffer
	//POST(ts.URL).Debug(true).SetJSONNotEscape(map[string]any{"url": "http://www.com?a=b&c=d"}).Do()
	POST(ts.URL).Debug(debug.ToWriter(&buf, false)).SetJSONNotEscape(map[string]any{"url": "http://www.com?a=b&c=d"}).Do()
	assert.True(t, bytes.Contains(buf.Bytes(), []byte("&")), buf.String())
	buf.Reset()
	//POST(ts.URL).Debug(true).SetJSON(map[string]any{"url": "http://www.com?a=b&c=d"}).Do()
	POST(ts.URL).Debug(debug.ToWriter(&buf, false)).SetJSON(map[string]any{"url": "http://www.com?a=b&c=d"}).Do()
	assert.False(t, bytes.Contains(buf.Bytes(), []byte("&")), buf.String())
}
