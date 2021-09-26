package encode

import (
	"strconv"
	"testing"
	"time"

	"github.com/guonaihong/gout/core"
	"github.com/guonaihong/gout/setting"
	"github.com/stretchr/testify/assert"
)

// 测试[]string类型
func TestQueryStringSlice(t *testing.T) {
	q := NewQueryEncode(setting.Setting{})

	err := Encode([]string{"q1", "v1", "q2", "v2", "q3", "v3"}, q)

	assert.NoError(t, err)

	assert.Equal(t, q.End(), "q1=v1&q2=v2&q3=v3")
}

// 测试map[string]interface{}
func TestQueryMap(t *testing.T) {

	q := NewQueryEncode(setting.Setting{})

	err := Encode(testH{"q1": "v1", "q2": "v2", "q3": "v3"}, q)

	assert.NoError(t, err)

	assert.Equal(t, q.End(), "q1=v1&q2=v2&q3=v3")
}

type testQuery struct {
	Q1 string    `query:"q1"`
	Q2 string    `query:"q2"`
	Q3 string    `query:"q3"`
	Q4 time.Time `query:"q4" time_format:"unix"`
	Q5 time.Time `query:"q5" time_format:"unixNano"`
}

// 测试结构体
func TestQueryStruct(t *testing.T) {

	q := NewQueryEncode(setting.Setting{})

	unixTime := time.Date(2019, 07, 27, 20, 42, 53, 0, time.Local)
	unixNano := time.Date(2019, 07, 27, 20, 42, 53, 1000, time.Local)

	err := Encode(testQuery{Q1: "v1", Q2: "v2", Q3: "v3", Q4: unixTime, Q5: unixNano}, q)

	assert.NoError(t, err)

	assert.Equal(t, q.End(), "q1=v1&q2=v2&q3=v3&q4="+strconv.FormatInt(unixTime.Unix(), 10)+"&q5="+strconv.FormatInt(unixNano.UnixNano(), 10))
}

// 结构体带[]string
type queryWithSlice struct {
	A []string `query:"a"`
	B string   `query:"b"`
}

func TestQueryFieldWithSlice(t *testing.T) {

	for _, v := range []interface{}{
		queryWithSlice{A: []string{"1", "2", "3"}, B: "b"},
		core.H{"a": []string{"1", "2", "3"}, "b": "b"},
		core.A{"a", []string{"1", "2", "3"}, "b", "b"},
	} {

		q := NewQueryEncode(setting.Setting{})

		err := Encode(v, q)

		assert.NoError(t, err)

		assert.Equal(t, q.End(), "a=1&a=2&a=3&b=b")
	}
}
