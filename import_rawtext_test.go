package gout

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRawText struct {
	A string `form:"a" json:"a"`
	B string `form:"b" json:"b"`
}

func setup_router(t *testing.T) *gin.Engine {
	router := gin.New()

	router.POST("/", func(c *gin.Context) {

		rt := testRawText{}
		err := c.Bind(&rt)
		assert.NoError(t, err)

		rt2 := testRawText{"a", "b"}
		assert.Equal(t, rt, rt2)
	})

	return router
}

func Test_RawText(t *testing.T) {
	all, err := ioutil.ReadFile("./testdata/raw-http-post-formdata.txt")
	assert.NoError(t, err)

	router := setup_router(t)

	code := 0
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	err = NewImport().RawText(all).Debug(true).SetURL(ts.URL).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

}

// 支持导入json
func Test_RawText_JSON(t *testing.T) {
	all, err := ioutil.ReadFile("./testdata/raw-http-post-json.txt")
	assert.NoError(t, err)

	type obj struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	type result struct {
		Result string `json:"result"`
	}

	type testRawTextJSON struct {
		Array []string `json:"array"`
		Num   int      `json:"num"`
		Str   string   `json:"str"`
		Obj   obj      `json:"obj"`
	}

	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/colorjson", func(c *gin.Context) {

			rt := testRawTextJSON{}
			err := c.ShouldBindJSON(&rt)
			assert.NoError(t, err)

			rt2 := testRawTextJSON{Array: []string{"foo", "bar", "baz"}, Str: "foo", Num: 100, Obj: obj{A: 1, B: 2}}
			assert.Equal(t, rt, rt2)
			c.JSON(200, gin.H{"result": "ok"})
		})

		return router
	}()

	code := 0
	var r, r1 result

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	err = NewImport().RawText(all).Debug(true).SetURL(ts.URL + "/colorjson").Code(&code).BindJSON(&r).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)
	r1 = result{"ok"}
	assert.Equal(t, r, r1)

}

// 修改URL
func Test_RawText_URL(t *testing.T) {
	all, err := ioutil.ReadFile("./testdata/raw-http-post-formdata.txt")
	assert.NoError(t, err)

	type testRawText struct {
		A string `form:"a"`
		B string `form:"b"`
	}

	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/", func(c *gin.Context) {

			rt := testRawText{}
			err := c.Bind(&rt)
			assert.NoError(t, err)

			rt2 := testRawText{"a", "b"}
			assert.Equal(t, rt, rt2)

			rt = testRawText{}
			rt2 = testRawText{"q1", "q2"}
			err = c.ShouldBindQuery(&rt)
			assert.NoError(t, err)
			assert.Equal(t, rt, rt2)
		})

		return router
	}()

	code := 0
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	err = NewImport().RawText(all).Debug(true).SetURL(ts.URL).
		SetQuery(H{"a": "q1", "b": "q2"}).
		Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)
}

// 测试错误情况
func Test_RawText_fail(t *testing.T) {

	router := setup_router(t)

	code := 0
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	for _, v := range []interface{}{new(int), "xxxx"} {
		err := NewImport().RawText(v).Debug(true).SetURL(ts.URL).Code(&code).Do()
		assert.Error(t, err)
	}

}

// 测试支持的类型
func Test_RawText_MoreType(t *testing.T) {
	all, err := ioutil.ReadFile("./testdata/raw-http-post-formdata.txt")
	assert.NoError(t, err)

	router := setup_router(t)

	code := 0
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	for _, v := range []interface{}{all, string(all)} {
		err = NewImport().RawText(v).Debug(true).SetURL(ts.URL).Code(&code).Do()
		assert.NoError(t, err)
		assert.Equal(t, code, 200)
	}
}
