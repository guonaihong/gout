package gout

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/guonaihong/gout/dataflow"
)

type queryWithSlice struct {
	A []string `query:"a" form:"a"`
	B string   `query:"b" form:"b"`
}

func testQueryWithSliceServer(t *testing.T) *httptest.Server {

	r := gin.New()

	need := queryWithSlice{A: []string{"1", "2", "3"}, B: "b"}
	r.GET("/query", func(c *gin.Context) {

		got := queryWithSlice{}
		err := c.ShouldBindQuery(&got)
		assert.NoError(t, err)
		assert.Equal(t, need, got)
	})

	return httptest.NewServer(http.HandlerFunc(r.ServeHTTP))
}

// 测试query接口，带slice的情况
func TestQuery_slice(t *testing.T) {

	ts := testQueryWithSliceServer(t)

	for _, v := range []interface{}{
		queryWithSlice{A: []string{"1", "2", "3"}, B: "b"},
		H{"a": []string{"1", "2", "3"}, "b": "b"},
		A{"a", []string{"1", "2", "3"}, "b", "b"},
	} {

		err := GET(ts.URL + "/query").Debug(true).SetQuery(v).Do()
		assert.NoError(t, err)
	}
}

func TestQuery_NotIgnoreEmpty(t *testing.T) {

	total := int32(0)
	router := setupMethod(&total)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	query := H{
		"t":          1296,
		"callback":   "searchresult",
		"q":          "美食",
		"stype":      1,
		"pagesize":   100,
		"pagenum":    1,
		"imageType":  2,
		"imageColor": "",
		"brand":      "",
		"imageSType": "",
		"fr":         1,
		"sortFlag":   1,
		"imageUType": "",
		"btype":      "",
		"authid":     "",
		"_":          int64(1611822443760),
	}

	var out bytes.Buffer
	SaveDebug := func() dataflow.DebugOpt {
		return DebugFunc(func(o *DebugOption) {
			o.Write = &out
			o.Debug = true
		})
	}

	// 默认不忽略空值
	err := GET(ts.URL).Debug(SaveDebug()).SetQuery(query).Do()
	assert.NoError(t, err)
	// 有authid字段
	assert.NotEqual(t, bytes.Index(out.Bytes(), []byte("authid")), -1)

	// 重置bytes.Buffer
	out.Reset()
	// 忽略空值
	IgnoreEmpty()
	// 默认不忽略空值
	err = GET(ts.URL).Debug(SaveDebug()).SetQuery(query).Do()
	assert.NoError(t, err)
	// 没有authid字段
	assert.Equal(t, bytes.Index(out.Bytes(), []byte("authid")), -1)

	NotIgnoreEmpty()
}
