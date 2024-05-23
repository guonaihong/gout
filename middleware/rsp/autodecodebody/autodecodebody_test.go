package autodecodebody

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
)

func TestAutoDecodeBody(t *testing.T) {
	// https://developer.mozilla.org/zh-CN/docs/web/http/headers/content-encoding
	data := "test auto decode boyd function"
	for _, rsp := range []http.Response{
		// gzip
		/*
			func() (rv http.Response) {
				var buf bytes.Buffer

				rv.Header = make(map[string][]string)
				rv.Header.Set("Content-Encoding", "gzip")
				zw := gzip.NewWriter(&buf)
				// Setting the Header fields is optional.
				zw.Name = "a-new-hope.txt"
				zw.Comment = "an epic space opera by George Lucas"
				zw.ModTime = time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC)
				_, err := zw.Write([]byte(data))
				if err != nil {
					log.Fatal(err)
				}

				if err := zw.Close(); err != nil {
					log.Fatal(err)
				}

				rv.Body = ioutil.NopCloser(&buf)
				return
			}(),
		*/
		// deflate
		func() (rv http.Response) {

			rv.Header = make(map[string][]string)
			rv.Header.Set("Content-Encoding", "deflate")
			var b bytes.Buffer
			w := zlib.NewWriter(&b)
			_, err := w.Write([]byte(data))
			assert.NoError(t, err)
			w.Close()

			rv.Body = ioutil.NopCloser(&b)
			return
		}(),
		// br
		func() (rv http.Response) {
			rv.Header = make(map[string][]string)
			rv.Header.Set("Content-Encoding", "br")
			var b bytes.Buffer
			w := brotli.NewWriter(&b)
			_, err := w.Write([]byte(data))
			assert.NoError(t, err)
			w.Flush()
			w.Close()
			rv.Body = ioutil.NopCloser(&b)
			return
		}(),
	} {

		err := AutoDecodeBody(&rsp)
		assert.NoError(t, err)
		all, err := ioutil.ReadAll(rsp.Body)
		assert.Equal(t, all, []byte(data))
		assert.NoError(t, err)
	}
}
