package gout

import (
	"bytes"
	"testing"

	"github.com/guonaihong/gout/dataflow"
	"github.com/guonaihong/gout/debug"
)

// 测试trace的入口函数
func Test_Trace(t *testing.T) {
	t.Run("TraceJSONToWriter", func(t *testing.T) {
		//大约是这样的字符串
		//{"DnsDuration":0,"ConnDuration":1490250,"TLSDuration":0,"RequestDuration":33166,"WaitResponeDuration":258084,"ResponseDuration":10708,"TotalDuration":1810834}
		ts := createMethodEcho()
		defer ts.Close()
		var buf bytes.Buffer
		err := dataflow.New().GET(ts.URL).Debug(debug.TraceJSONToWriter(&buf)).Do()
		if err != nil {
			t.Fatal(err)
		}
		pos := bytes.Index(buf.Bytes(), []byte("ConnDuration"))
		if pos == -1 {
			t.Fatal("not found ConnDuration")
		}
	})
	t.Run("TraceJSON", func(t *testing.T) {
		ts := createMethodEcho()
		defer ts.Close()
		err := dataflow.New().GET(ts.URL).Debug(debug.TraceJSON()).Do()
		if err != nil {
			t.Fatal(err)
		}
	})
}
