package decode

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
		//w.Header().Add("limitZero", "")
		w.Header().Add("f64Zero", "")
		w.Header().Add("f32Zero", "")
		w.Header().Add("createTimeZero", "")
		w.Header().Add("unixTimeZero", "")
		w.Header().Add("durationZero", "")
		w.Header().Add("boolZero", "")

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

	assert.Equal(t, 0, theader.LimitZero)
	assert.Equal(t, 0.0, theader.F64Zero)
	assert.Equal(t, float32(0.0), theader.F32Zero)
	assert.Equal(t, time.Time{}.UnixNano(), theader.CreateTimeZero.UnixNano())
	assert.Equal(t, time.Time{}.Unix(), theader.UnixTimeZero.Unix())
	assert.Equal(t, 0, int(theader.DurationZero))
	assert.Equal(t, false, theader.BoolZero)

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
