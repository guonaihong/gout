package decode

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHeaderDecode(t *testing.T) {
	h := headerDecode{}

	type tHeader struct {
		Limit      int           `header:"limit"`
		F64        float64       `header:"f64"`
		F32        float32       `header:"f32"`
		CreateTime time.Time     `header:"createTime" time_format:"unixNano"`
		UnixTime   time.Time     `header:"unixTime" time_format:"unix"`
		Duration   time.Duration `header:"duration"`
	}

	var theader tHeader

	okFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("limit", "1000")
		w.Header().Add("f64", "64")
		w.Header().Add("f32", "32.1")
		w.Header().Add("createTime", "1562400033000000123")
		w.Header().Add("unixTime", "1562400033")
		w.Header().Add("duration", "1h1s")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	okFunc(w, req)
	resp := w.Result()

	// todo array slice
	assert.NoError(t, h.Decode(resp, &theader))

	assert.Equal(t, 1000, theader.Limit)
	assert.Equal(t, 64.0, theader.F64)
	assert.Equal(t, float32(32.1), theader.F32)
	assert.Equal(t, int64(1562400033000000123), theader.CreateTime.UnixNano())
	assert.Equal(t, int64(1562400033), theader.UnixTime.Unix())
	assert.Equal(t, int(time.Hour)+int(time.Second), int(theader.Duration))

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
