package dataflow

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/guonaihong/gout/encode"
	"github.com/stretchr/testify/assert"
)

func TestReqModifyUrl(t *testing.T) {
	src := []string{"127.0.0.1", ":8080/query", "/query", "http://127.0.0.1", "https://127.0.0.1"}
	want := []string{"http://127.0.0.1", "http://127.0.0.1:8080/query", "http://127.0.0.1/query", "http://127.0.0.1", "https://127.0.0.1"}

	for k, v := range src {
		if want[k] != modifyURL(v) {
			t.Errorf("got %s want %s\n", modifyURL(v), want[k])
		}
	}
}

type urlTest struct {
	set  interface{}
	need interface{}
}

func TestReq_isAndGetString(t *testing.T) {
	test := []urlTest{
		{set: "?a=b&c=d", need: "a=b&c=d"},
		{set: "c=d&e=f", need: "c=d&e=f"},
		{set: []byte("c=d&e=f"), need: "c=d&e=f"},
		{set: time.Time{}, need: ""},
	}

	for _, v := range test {
		rv, _ := isAndGetString(v.set)
		assert.Equal(t, v.need, rv)
	}
}

// 测试request()函数调用出错的情况
func TestReq_request_fail(t *testing.T) {

	tests := []func() *Req{
		func() *Req {
			r := Req{}
			r.bodyEncoder = encode.NewBodyEncode(&map[string]string{})
			return &r
		},
		func() *Req {
			r := Req{}
			s := "hello"
			r.form = []interface{}{s}
			return &r
		},
	}

	for _, test := range tests {
		r := test()
		_, err := r.Request()
		assert.Error(t, err)
	}

}

type testValid struct {
	Val string `valid:"required"`
}

func Test_Valid(t *testing.T) {
	total := int32(1)
	router := setupMethod(&total)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	testCases := []string{"bindheader"}
	for _, c := range testCases {
		val := testValid{}
		g := GET(ts.URL + "/someGet")
		var err error
		switch c {
		case "bindheader":
			err = g.BindHeader(&val).Do()
		}

		//fmt.Printf("-->%v\n", err)
		assert.Error(t, err)
	}
}

func TestReadAll(t *testing.T) {
	size := 100
	body := bytes.Repeat([]byte{'a'}, size)
	cases := []struct {
		name    string
		resp    *http.Response
		expBody []byte
	}{
		{
			name: "reads known size",
			resp: &http.Response{
				ContentLength: int64(size),
				Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
			},
			expBody: body,
		},
		{
			name: "reads unknown size",
			resp: &http.Response{
				ContentLength: -1,
				Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
			},
			expBody: body,
		},
		{
			name: "supports empty with size=0",
			resp: &http.Response{
				ContentLength: 0,
				Body:          ioutil.NopCloser(bytes.NewBuffer(nil)),
			},
			expBody: []byte{},
		},
		{
			name: "supports empty with unknown size",
			resp: &http.Response{
				ContentLength: -1,
				Body:          ioutil.NopCloser(bytes.NewBuffer(nil)),
			},
			expBody: []byte{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actBody, err := ReadAll(tc.resp)
			require.NoError(t, err)
			require.Equal(t, tc.expBody, actBody)
		})
	}
}

func BenchmarkReadAll(b *testing.B) {
	sizes := []int{
		100,         // 100 bytes
		100 * 1024,  // 100KB
		1024 * 1024, // 1MB
	}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size: %d", size), func(b *testing.B) {

			// emulate a file or an HTTP response
			generated := bytes.Repeat([]byte{'a'}, size)
			content := bytes.NewReader(generated)
			cases := []struct {
				name string
				resp *http.Response
			}{
				{
					name: "unknown length",
					resp: &http.Response{
						ContentLength: -1,
						Body:          ioutil.NopCloser(content),
					},
				},
				{
					name: "known length",
					resp: &http.Response{
						ContentLength: int64(size),
						Body:          ioutil.NopCloser(content),
					},
				},
			}

			b.ResetTimer()

			for _, tc := range cases {
				b.Run(tc.name, func(b *testing.B) {
					b.Run("io.ReadAll", func(b *testing.B) {
						for i := 0; i < b.N; i++ {
							_, err := content.Seek(0, io.SeekStart) // reset
							require.NoError(b, err)
							data, err := ioutil.ReadAll(tc.resp.Body)
							require.NoError(b, err)
							require.Equalf(b, size, len(data), "size does not match, expected %d, actual %d", size, len(data))
						}
					})
					b.Run("bytes.Buffer+io.Copy", func(b *testing.B) {
						for i := 0; i < b.N; i++ {
							_, err := content.Seek(0, io.SeekStart) // reset
							require.NoError(b, err)
							data, err := ReadAll(tc.resp)
							require.NoError(b, err)
							require.Equalf(b, size, len(data), "size does not match, expected %d, actual %d", size, len(data))
						}
					})
				})
			}
		})
	}
}
