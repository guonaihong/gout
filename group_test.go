package gout

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
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
	err := g.GET(ts.URL + "/test.header").ShouldBindHeader(&tHeader).Code(&tHeader.Code).Do()
	assert.NoError(t, err)
	assert.Equal(t, tHeader.Code, 200)
	assert.Equal(t, tHeader.Sid, "sid-ok")
}

func TestToHeaderStruct(t *testing.T) {
	type testHeader2 struct {
		Q8 uint8 `header:"h8"`
	}

	type testHeader struct {
		Q1 string    `header:"h1"`
		Q2 int       `header:"h2"`
		Q3 float32   `header:"h3"`
		Q4 float64   `header:"h4"`
		Q5 time.Time `header:"h5" time_format:"unix"`
		Q6 time.Time `header:"h6" time_format:"unixNano"`
		Q7 time.Time `header:"h7" time_format:"2006-01-02"`
		testHeader2
	}

	h := testHeader{
		Q1: "v1",
		Q2: 2,
		Q3: 3.14,
		Q4: 3.1415,
		Q5: time.Date(2019, 7, 28, 14, 36, 0, 0, time.Local),
		Q6: time.Date(2019, 7, 28, 14, 36, 0, 1000, time.Local),
		Q7: time.Date(2019, 7, 28, 0, 0, 0, 0, time.Local),
		testHeader2: testHeader2{
			Q8: 8,
		},
	}

	router := func() *gin.Engine {
		router := gin.Default()
		router.GET("/test.header", func(c *gin.Context) {
			h2 := testHeader{}
			err := c.ShouldBindHeader(&h2)
			assert.NoError(t, err)

			assert.Equal(t, h, h2)
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	g := New(nil)
	code := 0

	err := g.GET(ts.URL + "/test.header").ToHeader(h).Code(&code).Do()

	assert.NoError(t, err)
}

func TestToQueryStruct(t *testing.T) {
	type testQuery2 struct {
		Q8 uint8 `query:"q8" form:"q8"`
	}

	type testQuery struct {
		Q1 string    `query:"q1" form:"q1"`
		Q2 int       `query:"q2" form:"q2"`
		Q3 float32   `query:"q3" form:"q3"`
		Q4 float64   `query:"q4" form:"q4"`
		Q5 time.Time `query:"q5" form:"q5" time_format:"unix"`
		Q6 time.Time `query:"q6" form:"q6" time_format:"unixNano"`
		Q7 time.Time `query:"q7" form:"q7" time_format:"2006-01-02"`
		testQuery2
	}

	q := testQuery{
		Q1: "v1",
		Q2: 2,
		Q3: 3.14,
		Q4: 3.1415,
		Q5: time.Date(2019, 7, 28, 14, 36, 0, 0, time.Local),
		Q6: time.Date(2019, 7, 28, 14, 36, 0, 1000, time.Local),
		Q7: time.Date(2019, 7, 28, 0, 0, 0, 0, time.Local),
		testQuery2: testQuery2{
			Q8: 8,
		},
	}

	router := func() *gin.Engine {
		router := gin.Default()
		router.GET("/test.query", func(c *gin.Context) {
			q2 := testQuery{}
			err := c.ShouldBindQuery(&q2)
			assert.NoError(t, err)

			assert.Equal(t, q, q2)
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	g := New(nil)
	code := 0

	err := g.GET(ts.URL + "/test.query").ToQuery(q).Code(&code).Do()

	assert.NoError(t, err)
}
