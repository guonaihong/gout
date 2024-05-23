package hcutil

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupUnixSocket(t *testing.T, path string) *http.Server {
	router := gin.New()
	type testHeader struct {
		H1 string `header:"h1"`
		H2 string `header:"h2"`
	}

	router.POST("/test/unix", func(c *gin.Context) {
		tHeader := testHeader{}
		err := c.ShouldBindHeader(&tHeader)

		assert.Equal(t, tHeader.H1, "v1")
		assert.Equal(t, tHeader.H2, "v2")
		assert.NoError(t, err)

		c.String(200, "ok")
	})

	listener, err := net.Listen("unix", path)
	assert.NoError(t, err)

	srv := http.Server{Handler: router}
	go func() {
		// 外层是通过context关闭， 所以这里会返回错误
		assert.Error(t, srv.Serve(listener))
	}()

	return &srv
}

func setupProxy(t *testing.T) *gin.Engine {
	r := gin.New()

	r.GET("/:a", func(c *gin.Context) {
		all, err := ioutil.ReadAll(c.Request.Body)

		assert.NoError(t, err)
		c.String(200, string(all))
	})

	return r
}

type TransportFail struct{}

func (t *TransportFail) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, errors.New("fail")
}

func TestProxy(t *testing.T) {
	router := setupProxy(t)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()
	proxyTs := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer proxyTs.Close()

	var s string
	var err error

	c := http.Client{}
	err = SetProxy(&c, proxyTs.URL)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", ts.URL+"/login", strings.NewReader(proxyTs.URL))
	assert.NoError(t, err)
	resp, err := c.Do(req)
	assert.NoError(t, err)

	res, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	s = string(res)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, s, proxyTs.URL)

	err = SetProxy(&c, "\x7f" /*url.Parse源代码写了遇到\x7f会报错*/)
	// test fail
	assert.Error(t, err)

	// 错误情况1
	c.Transport = &TransportFail{}
	req, err = http.NewRequest("GET", ts.URL+"/login", strings.NewReader(s))
	assert.NoError(t, err)
	_, err = c.Do(req)
	assert.Error(t, err)
}

func TestSetSOCKS5(t *testing.T) {
	// TODO
}

func TestUnixSocket(t *testing.T) {
	path := "./unix.sock"
	defer os.Remove(path)

	ctx, cancel := context.WithCancel(context.Background())
	srv := setupUnixSocket(t, path)
	defer func() {
		assert.NoError(t, srv.Shutdown(ctx))
		cancel()
	}()

	c := http.Client{}
	err := UnixSocket(&c, path)
	assert.NoError(t, err)
	s := ""

	req, err := http.NewRequest("POST", "http://xxx/test/unix/", nil)
	assert.NoError(t, err)
	req.Header.Add("h1", "v1")
	req.Header.Add("h2", "v2")

	resp, err := c.Do(req)

	assert.NoError(t, err)
	all, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	s = string(all)

	// err := New(&c).UnixSocket(path).POST("http://xxx/test/unix/").SetHeader(core.H{"h1": "v1", "h2": "v2"}).BindBody(&s).Do()
	assert.NoError(t, err)
	assert.Equal(t, s, "ok")

	// 错误情况1
	c.Transport = &TransportFail{}

	_, err = c.Do(req)

	assert.Error(t, err)
}
