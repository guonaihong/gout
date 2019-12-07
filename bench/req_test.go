package bench

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func Test_Bench_closeRequest(t *testing.T) {
	b := bytes.NewBuffer([]byte("hello"))
	req, err := http.NewRequest("GET", "hello", b)
	assert.NoError(t, err)

	req.Header.Add("h1", "h2")
	req.Header.Add("h2", "h2")

	req2, err := cloneRequest(req)
	assert.NoError(t, err)

	// 测试http header是否一样
	assert.Equal(t, req.Header, req2.Header)

	b2, err := req2.GetBody()
	assert.NoError(t, err)

	b3 := bytes.NewBuffer(nil)
	io.Copy(b3, b2)

	// 测试body是否一样
	assert.Equal(t, b, b3)
}
