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

		Tstruct Tstruct  `header:"tstruct"`
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

type time1 struct {
	Time time.Time `test:"time"`
}

type time2 struct {
	Time time.Time `test:"time"`
}

type time3 struct {
	Time time.Time `test:"time" time_format:"unixNano"`
}

type timeLocationFail struct {
	Time time.Time `test:"time" time_location:"xxx"`
}

func Test_Core_setTimeField_Fail(t *testing.T) {
	//测试时间
	tests := []reflect.Value{

		reflect.ValueOf(time3{}),
		reflect.ValueOf(timeLocationFail{}),
	}

	for _, v := range tests {
		err := setTimeField("xx", 0, v.Type().Field(0), v.Field(0))
		assert.Error(t, err)
	}
}

type decodeTest struct {
	set  interface{}
	need interface{}
	in   string
}

func Test_Core_setTimeField(t *testing.T) {
	//测试时间
	format := "2006-01-02 15:04:05"
	parseInLocation := func(layout, value string, locStr string) time.Time {
		loc, err := time.LoadLocation(locStr)
		assert.NoError(t, err)

		tm, err := time.ParseInLocation(layout, value, loc)
		assert.NoError(t, err)
		return tm
	}

	tm := time.Now()
	tests := []decodeTest{
		{
			&struct {
				T time.Time `time_format:"2006-01-02 15:04:05" time_location:"Asia/Shanghai"`
			}{},
			&struct {
				T time.Time `time_format:"2006-01-02 15:04:05" time_location:"Asia/Shanghai"`
			}{T: parseInLocation(format, tm.Format(format), "Asia/Shanghai")},
			tm.Format(format),
		},
		{
			&struct {
				T time.Time `time_format:"2006-01-02 15:04:05" time_location:"Asia/Chongqing"`
			}{},
			&struct {
				T time.Time `time_format:"2006-01-02 15:04:05" time_location:"Asia/Chongqing"`
			}{T: parseInLocation(format, tm.Format(format), "Asia/Chongqing")},
			tm.Format(format),
		},
	}

	for _, test := range tests {
		val := reflect.ValueOf(test.set)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		err := setTimeField(test.in, 0, val.Type().Field(0), val.Field(0))
		assert.NoError(t, err)
	}
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

// 测试数组类型
func Test_Core_setForm_array(t *testing.T) {
	//测试出错
	m := map[string][]string{
		"testarray-fail": []string{"123"},
	}

	var s [3]string
	err := setForm(m, reflect.ValueOf(s), emptyField, "testarray-fail")
	assert.Error(t, err)

	//测试ok

	m = map[string][]string{
		"testarray-ok": []string{"123", "456"},
	}

	var ok [2]string
	v := reflect.ValueOf(&ok)
	v = v.Elem()
	err = setForm(m, v, emptyField, "testarray-ok")
	assert.NoError(t, err)
	assert.Equal(t, ok, [2]string{"123", "456"})
}
