package gout

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guonaihong/gout/dataflow"
)

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
		"_":          1611822443760,
	}

	var out bytes.Buffer
	SaveDebug := func() dataflow.DebugOpt {
		return DebugFunc(func(o *DebugOption) {
			o.Write = &out
			o.Debug = true
		})
	}

	// 默认忽略空值
	err := GET(ts.URL).Debug(SaveDebug()).SetQuery(query).Do()
	assert.NoError(t, err)
	// 没有authid字段
	assert.Equal(t, bytes.Index(out.Bytes(), []byte("authid")), -1)

	// 重置bytes.Buffer
	out.Reset()
	NotIgnoreEmpty()

	err = GET(ts.URL).Debug(SaveDebug()).SetQuery(query).Do()
	assert.NoError(t, err)
	// 有authid字段
	assert.NotEqual(t, bytes.Index(out.Bytes(), []byte("authid")), -1)

	IgnoreEmpty()
}
