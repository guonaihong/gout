package dataflow

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var test_autoDecodeBody_data = "test auto decode boyd function"

// 测试服务
func create_AutoDecodeBody() *httptest.Server {
	r := gin.New()
	r.GET("/gzip", func(c *gin.Context) {

		var buf bytes.Buffer

		zw := gzip.NewWriter(&buf)
		// Setting the Header fields is optional.
		zw.Name = "a-new-hope.txt"
		zw.Comment = "an epic space opera by George Lucas"
		zw.ModTime = time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC)
		_, err := zw.Write([]byte(test_autoDecodeBody_data))
		if err != nil {
			log.Fatal(err)
		}

		if err := zw.Close(); err != nil {
			log.Fatal(err)
		}
		c.Header("Content-Encoding", "gzip")
		c.String(200, buf.String())
	})

	r.GET("/br", func(c *gin.Context) {

		c.Header("Content-Encoding", "br")
		var buf bytes.Buffer
		w := brotli.NewWriter(&buf)
		w.Write([]byte(test_autoDecodeBody_data))
		w.Flush()
		w.Close()
		c.String(200, buf.String())
	})

	r.GET("/deflate", func(c *gin.Context) {

		var buf bytes.Buffer
		w := zlib.NewWriter(&buf)
		w.Write([]byte(test_autoDecodeBody_data))
		w.Close()
		c.Header("Content-Encoding", "deflate")
		c.String(200, buf.String())
	})

	r.GET("/compress", func(c *gin.Context) {
		c.Header("Content-Encoding", "compress")
	})
	return httptest.NewServer(http.HandlerFunc(r.ServeHTTP))
}

func Test_AutoDecodeBody(t *testing.T) {
	ts := create_AutoDecodeBody()
	var err error
	for _, path := range []string{"/gzip", "/br", "/deflate"} {
		s := ""
		if path == "/gzip" {
			err = New().GET(ts.URL + path).Debug(true).BindBody(&s).Do()
		} else {
			err = New().GET(ts.URL + path).AutoDecodeBody().Debug(true).BindBody(&s).Do()

		}
		assert.NoError(t, err)
		assert.Equal(t, s, test_autoDecodeBody_data)
	}
}

func Test_AutoDecodeBody_Fail(t *testing.T) {
	ts := create_AutoDecodeBody()
	err := New().GET(ts.URL + "/compress").AutoDecodeBody().Do()
	assert.Error(t, err)
}
