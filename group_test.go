package gout

import (
	"fmt"
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
	Id   int    `json:"id" xml:"id"`
	Data string `json:"data xml:"data""`
}

type BindTest struct {
	InBody   interface{}
	OutBody  interface{}
	httpCode int
}

func TestShouldBindXML(t *testing.T) {
	var d, d2 data
	router := func() *gin.Engine {
		router := gin.Default()

		router.POST("/test.xml", func(c *gin.Context) {
			var d3 data
			err := c.ShouldBindXML(&d3)
			assert.NoError(t, err)
			c.XML(200, d3)
		})
		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	g := New(nil)

	d.Id = 3
	d.Data = "test data"

	code := 200

	err := g.POST(ts.URL + "/test.xml").ToXML(&d).ShouldBindXML(&d2).Code(&code).Do()

	assert.NoError(t, err)
	assert.Equal(t, code, 200)
	assert.Equal(t, d, d2)
}

func TestShouldBindYAML(t *testing.T) {
	var d, d2 data
	router := func() *gin.Engine {
		router := gin.Default()

		router.POST("/test.yaml", func(c *gin.Context) {
			var d3 data
			err := c.ShouldBindYAML(&d3)
			assert.NoError(t, err)
			c.YAML(200, d3)
		})
		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	g := New(nil)

	d.Id = 3
	d.Data = "test yaml data"

	code := 200

	err := g.POST(ts.URL + "/test.yaml").ToYAML(&d).ShouldBindYAML(&d2).Code(&code).Do()

	assert.NoError(t, err)
	assert.Equal(t, code, 200)
	assert.Equal(t, d, d2)
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

type testForm struct {
	Mode string `form:"mode"`
	Text string `form:"text"`
	//Voice []byte `form:"voice" form-mem:"true"` //todo open
}

func TestToForm(t *testing.T) {
	reqTestForm := testForm{
		Mode: "A",
		Text: "good morning",
		//Voice: []byte("pcm data"),
	}

	router := func() *gin.Engine {
		router := gin.Default()
		router.POST("/test.form", func(c *gin.Context) {

			t2 := testForm{}
			err := c.ShouldBind(&t2)
			assert.NoError(t, err)
			//assert.Equal(t, reqTestForm, t2)
			assert.Equal(t, reqTestForm.Mode, t2.Mode)
			assert.Equal(t, reqTestForm.Text, t2.Text)
		})
		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	g := New(nil)
	code := 0
	err := g.POST(ts.URL + "/test.form").ToForm(&reqTestForm).Code(&code).Do()

	assert.NoError(t, err)
}

func TestToHeaderMap(t *testing.T) {
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

type testQuery2 struct {
	Q8 uint8 `query:"q8" form:"q8"`
}

type testQuery struct {
	Q1 string    `query:"q1" form:"q1"`
	Q2 int       `query:"q2" form:"q2"`
	Q3 float32   `query:"q3" form:"q3"`
	Q4 float64   `query:"q4" form:"q4"`
	Q5 time.Time `query:"q5" form:"q5" time_format:"unix" time_location:"Asia/Shanghai"`
	Q6 time.Time `query:"q6" form:"q6" time_format:"unixNano" time_location:"Asia/Shanghai"`
	Q7 time.Time `query:"q7" form:"q7" time_format:"2006-01-02" time_location:"Asia/Shanghai"`
	testQuery2
}

func queryDefault() *testQuery {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err.Error())
	}

	return &testQuery{
		Q1: "v1",
		Q2: 2,
		Q3: 3.14,
		Q4: 3.1415,
		Q5: time.Date(2019, 7, 28, 14, 36, 0, 0, loc),
		Q6: time.Date(2019, 7, 28, 14, 36, 0, 1000, loc),
		Q7: time.Date(2019, 7, 28, 0, 0, 0, 0, loc),
		testQuery2: testQuery2{
			Q8: 8,
		},
	}
}

func TestToQueryStruct(t *testing.T) {
	q := queryDefault()
	router := func() *gin.Engine {
		router := gin.Default()
		router.GET("/test.query", func(c *gin.Context) {
			q2 := testQuery{}
			err := c.ShouldBindQuery(&q2)
			assert.NoError(t, err)

			testQueryEqual(t, *q, q2)
			//assert.Equal(t, q, &q2)
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

func TestQueryString(t *testing.T) {
	s := "q1=v1&q2=2&q3=3.14&q4=3.1415&q5=1564295760&q6=1564295760000001000&q7=2019-07-28&q8=8"
	testQueryStringCore(t, s, false)
	testQueryStringCore(t, s, true)
}

func testQueryEqual(t *testing.T, q1, q2 testQuery) {
	//不用assert.Equal(t, q1, q2)
	//assert.Equal 有个bug

	assert.Equal(t, q1.Q1, q2.Q1)
	assert.Equal(t, q1.Q2, q2.Q2)
	assert.Equal(t, q1.Q3, q2.Q3)
	assert.Equal(t, q1.Q4, q2.Q4)
	assert.Equal(t, q1.Q8, q2.Q8)
	if !q1.Q5.Equal(q2.Q5) {
		t.Errorf("want(%s) got(%s)\n", q1.Q5, q2.Q5)
	}
	if !q1.Q6.Equal(q2.Q6) {
		t.Errorf("want(%s) got(%s)\n", q1.Q6, q2.Q6)
	}
	if !q1.Q7.Equal(q2.Q7) {
		t.Errorf("want(%s) got(%s)\n", q1.Q7, q2.Q7)
	}
}

func testQueryStringCore(t *testing.T, qStr string, isPtr bool) {
	q := queryDefault()
	router := func() *gin.Engine {
		router := gin.Default()
		router.GET("/test.query", func(c *gin.Context) {
			q2 := testQuery{}
			err := c.ShouldBindQuery(&q2)

			//todo
			//fmt.Printf("------->q7(%t)\n", reflect.DeepEqual(q1.Q5, q2.Q5), reflect.DeepEqual(q1.Q6, q2.Q6), reflect.DeepEqual(q1.Q7, q2.Q7))
			//fmt.Printf("%s:%s, %t\n", q1.Q7, q2.Q7, q1.Q7.Equal(q2.Q7))
			assert.NoError(t, err)

			testQueryEqual(t, *q, q2)
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	g := New(nil)
	code := 0

	var err error
	if isPtr {
		err = g.GET(ts.URL + "/test.query").ToQuery(&qStr).Code(&code).Do()
	} else {
		err = g.GET(ts.URL + "/test.query").ToQuery(qStr).Code(&code).Do()
	}

	assert.NoError(t, err)
	assert.Equal(t, code, 200)
}

func setupQuery(t *testing.T, q *testQuery) func() *gin.Engine {
	return func() *gin.Engine {
		router := gin.Default()
		router.GET("/test.query", func(c *gin.Context) {
			q2 := testQuery{}
			err := c.ShouldBindQuery(&q2)
			assert.NoError(t, err)

			testQueryEqual(t, *q, q2)
			//assert.Equal(t, *q, q2)
		})

		return router
	}
}

func testQuerySliceAndArrayCore(t *testing.T, x interface{}) {
	q := queryDefault()
	router := setupQuery(t, q)

	ts := httptest.NewServer(http.HandlerFunc(router().ServeHTTP))
	defer ts.Close()

	g := New(nil)

	code := 0
	err := g.GET(ts.URL + "/test.query").ToQuery(x).Code(&code).Do()

	assert.NoError(t, err)
	assert.Equal(t, code, 200)
}

func testQueryFail(t *testing.T, x interface{}) {
	q := queryDefault()
	router := setupQuery(t, q)

	ts := httptest.NewServer(http.HandlerFunc(router().ServeHTTP))
	defer ts.Close()

	g := New(nil)

	code := 0
	err := g.GET(ts.URL + "/test.query").ToQuery(x).Code(&code).Do()
	assert.Error(t, err)
	assert.NotEqual(t, code, 200)
}

func TestQueryFail(t *testing.T) {
	testQueryFail(t, []string{"q1"})
}

func TestQuerySliceAndArray(t *testing.T) {
	q := queryDefault()
	testQuerySliceAndArrayCore(t, []string{"q1", "v1", "q2", "2", "q3", "3.14", "q4", "3.1415", "q5",
		fmt.Sprint(q.Q5.Unix()), "q6", fmt.Sprint(q.Q6.UnixNano()), "q7", q.Q7.Format("2006-01-02"), "q8", "8"})
	testQuerySliceAndArrayCore(t, [8 * 2]string{"q1", "v1", "q2", "2", "q3", "3.14", "q4", "3.1415", "q5",
		fmt.Sprint(q.Q5.Unix()), "q6", fmt.Sprint(q.Q6.UnixNano()), "q7", q.Q7.Format("2006-01-02"), "q8", "8"})

	testQuerySliceAndArrayCore(t, &[]string{"q1", "v1", "q2", "2", "q3", "3.14", "q4", "3.1415", "q5",
		fmt.Sprint(q.Q5.Unix()), "q6", fmt.Sprint(q.Q6.UnixNano()), "q7", q.Q7.Format("2006-01-02"), "q8", "8"})
	testQuerySliceAndArrayCore(t, &[8 * 2]string{"q1", "v1", "q2", "2", "q3", "3.14", "q4", "3.1415", "q5",
		fmt.Sprint(q.Q5.Unix()), "q6", fmt.Sprint(q.Q6.UnixNano()), "q7", q.Q7.Format("2006-01-02"), "q8", "8"})
}
