package gout

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type testGroup struct {
	send  bool
	total int32
}

func setupGroupRouter(t *testGroup) *gin.Engine {

	router := gin.Default()

	postFunc := func(c *gin.Context) {
		atomic.AddInt32(&t.total, 1)
	}

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", postFunc)
		v1.POST("/submit", postFunc)
		v1.POST("/read", postFunc)
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", postFunc)
		v2.POST("/submit", postFunc)
		v2.POST("/read", postFunc)
	}

	v2_1 := v2.Group("/v1")
	{
		v2_1.POST("/login", func(c *gin.Context) {
			t.send = true
			atomic.AddInt32(&t.total, 1)
		})
	}
	return router
}

func TestGroupNew(t *testing.T) {
	tstGroup := testGroup{}

	router := setupGroupRouter(&tstGroup)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	out := New(nil)

	// http://127.0.0.1:80/v1
	v1 := out.Group(ts.URL + "/v1")
	err := v1.POST("/login").Next().
		POST("/submit").Next().
		POST("/read").Do()

	if err != nil {
		t.Errorf("error:%s\n", err)
	}

	// http://127.0.0.1:80/v2
	v2 := out.Group(ts.URL + "/v2")
	err = v2.POST("/login").Next().
		POST("/submit").Next().
		POST("/read").Do()

	if err != nil {
		t.Errorf("error:%s\n", err)
	}

	v2_1 := v2.Group("/v1")
	err = v2_1.POST("/login").Do()

	if err != nil {
		t.Errorf("error:%s\n", err)
	}

	if tstGroup.total != 7 {
		t.Errorf("got %d want 7\n", tstGroup.total)
	}

	if !tstGroup.send {
		t.Errorf("/v2/v1/login fail\n")
	}
}

type data struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}

type BindTest struct {
	InBody   interface{}
	OutBody  interface{}
	httpCode int
}

func TestShouldBindJSON(t *testing.T) {
	var d3 data
	router := func() *gin.Engine {
		router := gin.Default()

		router.POST("/test.json", func(c *gin.Context) {
			c.ShouldBindJSON(&d3)
			c.JSON(200, d3)
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	g := New(nil)

	tests := []BindTest{
		{InBody: data{Id: 9, Data: "早上测试结构体"}, OutBody: data{}},
		{InBody: H{"id": 10, "data": "早上测试map"}, OutBody: data{}},
	}

	for k, _ := range tests {
		t.Logf("outbody type:%T:%p\n", tests[k].OutBody, &tests[k].OutBody)

		err := g.POST(ts.URL + "/test.json").ToJSON(&tests[k].InBody).ShouldBindJSON(&tests[k].OutBody).Code(&tests[k].httpCode).Do()
		if err != nil {
			t.Errorf("send fail:%s\n", err)
		}

		assert.Equal(t, tests[k].httpCode, 200)

		if tests[k].OutBody.(map[string]interface{})["id"].(float64) != float64(d3.Id) {
			t.Errorf("got:%#v(P:%p), want:%#v(T:%T)\n", tests[k].OutBody, &tests[k].OutBody, tests[k].InBody, tests[k].InBody)
		}

		/*
			if !reflect.DeepEqual(&tests[k].InBody, &tests[k].OutBody) {
				t.Errorf("got:%#v(P:%p), want:%#v(T:%T)\n", tests[k].OutBody, &tests[k].OutBody, tests[k].InBody, tests[k].InBody)
			}
		*/

	}
}

func TestShouldBindHeader(t *testing.T) {
	router := func() *gin.Engine {
		router := gin.Default()

		router.GET("/test.header", func(c *gin.Context) {
			c.Writer.Header().Add("sid", "sid-ok")
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	g := New(nil)

	type testHeader struct {
		Sid  string `header:"sid"`
		Code int
	}

	var tHeader testHeader
	g.GET(ts.URL + "/test.header").ShouldBindHeader(&t).Code(&tHeader.Code).Do()
	assert.Equal(t, tHeader.Code, 200)
	assert.Equal(t, tHeader.Sid, "sid-ok")
}
