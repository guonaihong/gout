package dataflow

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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

func TestSetQueryStruct(t *testing.T) {
	q := queryDefault()

	router := func() *gin.Engine {
		router := gin.New()
		router.GET("/test.query", func(c *gin.Context) {
			q2 := testQuery{}
			err := c.BindQuery(&q2)
			assert.NoError(t, err)

			testQueryEqual(t, *q, q2)
			//assert.Equal(t, q, &q2)
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	code := 0

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err.Error())
	}

	for _, v := range [][]interface{}{
		[]interface{}{queryDefault()}, //一个结构体
		[]interface{}{ //两个结构体
			testQuery{
				Q1: "v1",
				Q2: 2,
				Q3: 3.14,
				Q4: 3.1415,
			},
			testQuery{
				Q5: time.Date(2019, 7, 28, 14, 36, 0, 0, loc),
				Q6: time.Date(2019, 7, 28, 14, 36, 0, 1000, loc),
				Q7: time.Date(2019, 7, 28, 0, 0, 0, 0, loc),
				testQuery2: testQuery2{
					Q8: 8,
				},
			},
		},
		[]interface{}{ //两个结构体，两个字符串混着用
			testQuery{
				Q1: "v1",
				Q2: 2,
				Q3: 3.14,
			},
			"q4=3.1415",
			testQuery{
				Q5: time.Date(2019, 7, 28, 14, 36, 0, 0, loc),
				Q6: time.Date(2019, 7, 28, 14, 36, 0, 1000, loc),
				Q7: time.Date(2019, 7, 28, 0, 0, 0, 0, loc),
			},
			"q8=8",
		},
	} {
		err := New().GET(ts.URL + "/test.query").SetQuery(v...).Code(&code).Do()
		assert.NoError(t, err)
		if err != nil {
			break
		}
	}

}

func TestQueryRaw(t *testing.T) {
	s := "q1=v1&q2=2&q3=3.14&q4=3.1415&q5=1564295760&q6=1564295760000001000&q7=2019-07-28&q8=8"
	b := []byte(s)

	q := queryDefault()
	router := func() *gin.Engine {
		router := gin.New()
		router.GET("/test.query", func(c *gin.Context) {
			q2 := testQuery{}
			err := c.BindQuery(&q2)

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

	code := 0

	for _, data := range []interface{}{
		s,
		&s,
		b,
		&b,
	} {

		err := GET(ts.URL + "/test.query").SetQuery(data).Code(&code).Do()
		assert.NoError(t, err)
		assert.Equal(t, code, 200)
	}

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

func setupQuery(t *testing.T, q *testQuery) func() *gin.Engine {
	return func() *gin.Engine {
		router := gin.New()
		router.GET("/test.query", func(c *gin.Context) {
			q2 := testQuery{}
			err := c.BindQuery(&q2)
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
	err := g.GET(ts.URL + "/test.query").SetQuery(x).Code(&code).Do()

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
	err := g.GET(ts.URL + "/test.query").SetQuery(x).Code(&code).Do()
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
