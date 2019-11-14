package decode

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type Tstruct struct {
	X int
	Y int
}

func TestHeaderDecode(t *testing.T) {
	h := headerDecode{}

	type tHeader struct {
		Limit      int           `header:"limit"`
		F64        float64       `header:"f64"`
		F32        float32       `header:"f32"`
		CreateTime time.Time     `header:"createTime" time_format:"unixNano"`
		UnixTime   time.Time     `header:"unixTime" time_format:"unix"`
		Duration   time.Duration `header:"duration"`
		Bool       bool          `header:"bool"`

		LimitZero      int           `header:"limitZero"`
		F64Zero        float64       `header:"f64Zero"`
		F32Zero        float32       `header:"f32Zero"`
		CreateTimeZero time.Time     `header:"createTimeZero" time_format:"unixNano"`
		UnixTimeZero   time.Time     `header:"unixTimeZero" time_format:"unix"`
		DurationZero   time.Duration `header:"durationZero"`
		BoolZero       bool          `header:"boolZero"`
		UintZero       uint          `header:"uintZero"`

		Tstruct Tstruct  `header:"struct"`
		Tslice  []string `header:"tslice"`
	}

	var theader tHeader

	okFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("limit", "1000")
		w.Header().Add("f64", "64")
		w.Header().Add("f32", "32.1")
		w.Header().Add("createTime", "1562400033000000123")
		w.Header().Add("unixTime", "1562400033")
		w.Header().Add("duration", "1h1s")
		w.Header().Add("bool", "true")

		// 测试slice
		w.Header().Add("tslice", "1")
		w.Header().Add("tslice", "2")
		w.Header().Add("tslice", "3")
		w.Header().Add("tslice", "4")

		// 测试0值
		w.Header().Add("limitZero", "")
		w.Header().Add("f64Zero", "")
		w.Header().Add("f32Zero", "")
		w.Header().Add("createTimeZero", "")
		w.Header().Add("unixTimeZero", "")
		w.Header().Add("durationZero", "")
		w.Header().Add("boolZero", "")
		w.Header().Add("uintZero", "")

		w.Header().Add("tstruct", `{"x":3, "y":4}`)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	okFunc(w, req)
	resp := w.Result()

	assert.NoError(t, h.Decode(resp, &theader))

	// 测试slice
	assert.Equal(t, []string{"1", "2", "3", "4"}, theader.Tslice)

	assert.Equal(t, 1000, theader.Limit)
	assert.Equal(t, 64.0, theader.F64)
	assert.Equal(t, float32(32.1), theader.F32)
	assert.Equal(t, int64(1562400033000000123), theader.CreateTime.UnixNano())
	assert.Equal(t, int64(1562400033), theader.UnixTime.Unix())
	assert.Equal(t, int(time.Hour)+int(time.Second), int(theader.Duration))
	assert.Equal(t, true, theader.Bool)

	assert.Equal(t, Tstruct{X: 3, Y: 4}, theader.Tstruct)

	// 测试0值
	assert.Equal(t, 0, theader.LimitZero)
	assert.Equal(t, 0.0, theader.F64Zero)
	assert.Equal(t, float32(0.0), theader.F32Zero)
	assert.Equal(t, time.Time{}.UnixNano(), theader.CreateTimeZero.UnixNano())
	assert.Equal(t, time.Time{}.Unix(), theader.UnixTimeZero.Unix())
	assert.Equal(t, 0, int(theader.DurationZero))
	assert.Equal(t, false, theader.BoolZero)
	assert.Equal(t, uint(0), theader.UintZero)

	failFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("fail", `{fail:fail}`)
	}

	type failStruct struct {
		Fail map[string]interface{} `header:"fail"`
	}

	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	w = httptest.NewRecorder()
	failFunc(w, req)
	resp = w.Result()

	err := h.Decode(resp, &failStruct{})
	assert.Error(t, err)
}

func Test_Core_setBase_fail(t *testing.T) {
	//测试出错
	err := setBase("x", emptyField, reflect.ValueOf(make(chan struct{})))
	assert.Error(t, err)
}

func Test_Core_setTimeDuration_fail(t *testing.T) {
	//测试出错
	err := setTimeDuration("xx", 0, emptyField, reflect.Value{})
	assert.Error(t, err)
}

func Test_Core_setSlice_fail(t *testing.T) {
	//测试出错
	err := setSlice([]string{"1", "2"}, emptyField, reflect.ValueOf([]chan struct{}{make(chan struct{})}))
	assert.Error(t, err)
}

//TODO
func Test_Core_setTimeField(t *testing.T) {
	//测试
}

type emptySet struct{}

func (e *emptySet) Set(value reflect.Value, sf reflect.StructField, tagValue string) error {
	return nil
}

//测试空返回错误
func Test_Core_decode_empty(t *testing.T) {
	var p *int
	err := []error{
		decode(&emptySet{}, nil, "test empty"),
		decode(&emptySet{}, p, "test empty"),
	}

	for _, e := range err {
		assert.Error(t, e)
	}
}

func Test_Core_setForm_fail(t *testing.T) {
}
