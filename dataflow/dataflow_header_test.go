package dataflow

import (
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testHeader2 struct {
	H8 uint8 `header:"h8"`
}

type testSetHeader struct {
	H1 string    `header:"h1"`
	H2 int       `header:"h2"`
	H3 float32   `header:"h3"`
	H4 float64   `header:"h4"`
	H5 time.Time `header:"h5" time_format:"unix"`
	H6 time.Time `header:"h6" time_format:"unixNano"`
	H7 time.Time `header:"h7" time_format:"2006-01-02"`
	testHeader2
}

func createHeader(t *testing.T, h testSetHeader) *gin.Engine {
	router := gin.New()
	router.GET("/test.header", func(c *gin.Context) {
		h2 := testSetHeader{}
		err := c.BindHeader(&h2)
		assert.NoError(t, err)

		assert.Equal(t, h, h2)
	})

	return router
}

func Test_SetHeaderMulti_SpecialCases(t *testing.T) {
	router := createHeader(t, testSetHeader{})
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	code := 0
	for _, h := range []interface{}{nil} {

		err := GET(ts.URL + "/test.header").SetHeader(h).Code(&code).Do()
		assert.NoError(t, err)
	}
}

func Test_SetHeaderMulti(t *testing.T) {
	h := testSetHeader{
		H1: "v1",
		H2: 2,
		H3: 3.14,
		H4: 3.1415,
		H5: time.Date(2019, 7, 28, 14, 36, 0, 0, time.Local),
		H6: time.Date(2019, 7, 28, 14, 36, 0, 1000, time.Local),
		H7: time.Date(2019, 7, 28, 0, 0, 0, 0, time.Local),
		testHeader2: testHeader2{
			H8: 8,
		},
	}

	h2 := core.A{
		"h1", "v1",
		"h2", 2,
		"h3", 3.14,
		"h4", 3.1415,
	}

	h3 := testSetHeader{
		H5: h.H5,
		H6: h.H6,
		H7: h.H7,
		testHeader2: testHeader2{
			H8: 8,
		},
	}

	router := createHeader(t, h)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	code := 0

	err := GET(ts.URL+"/test.header").SetHeader(h2, h3).Code(&code).Do()

	assert.NoError(t, err)
}

func Test_SetHeaderMap(t *testing.T) {
}

func Test_SetHeaderStruct(t *testing.T) {
	h := testSetHeader{
		H1: "v1",
		H2: 2,
		H3: 3.14,
		H4: 3.1415,
		H5: time.Date(2019, 7, 28, 14, 36, 0, 0, time.Local),
		H6: time.Date(2019, 7, 28, 14, 36, 0, 1000, time.Local),
		H7: time.Date(2019, 7, 28, 0, 0, 0, 0, time.Local),
		testHeader2: testHeader2{
			H8: 8,
		},
	}

	router := createHeader(t, h)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	code := 0

	err := GET(ts.URL + "/test.header").SetHeader(h).Code(&code).Do()

	assert.NoError(t, err)
}

func Test_BindHeader(t *testing.T) {
	router := func() *gin.Engine {
		router := gin.New()

		router.GET("/test.header", func(c *gin.Context) {
			c.Writer.Header().Add("sid", "sid-ok")
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	type testHeader2 struct {
		Sid  string `header:"sid"`
		Code int
	}

	var tHeader testHeader2
	err := New().GET(ts.URL + "/test.header").BindHeader(&tHeader).Code(&tHeader.Code).Do()
	assert.NoError(t, err)
	assert.Equal(t, tHeader.Code, 200)
	assert.Equal(t, tHeader.Sid, "sid-ok")
}

// 测试设置空header
func Test_BindHeader_empty(t *testing.T) {
	router := func() *gin.Engine {
		router := gin.New()

		router.GET("/test.header", func(c *gin.Context) {
			c.Writer.Header().Add("sid", "sid-ok")
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	g := New()

	type testHeader2 struct {
		Sid  string `header:"sid"`
		Code int
	}

	var tHeader testHeader2
	for _, v := range []interface{}{
		core.A{},
		core.H{},
	} {
		err := g.GET(ts.URL + "/test.header").SetHeader(v).BindHeader(&tHeader).Code(&tHeader.Code).Do()
		assert.NoError(t, err)
		assert.Equal(t, tHeader.Code, 200)
		assert.Equal(t, tHeader.Sid, "sid-ok")
	}
}
