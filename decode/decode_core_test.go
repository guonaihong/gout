package decode

import (
	"github.com/stretchr/testify/assert"
)

func TestHeaderDecode(t *testing.T) {
	h := headerDecode

	type tHeader struct {
		Limit      int       `header:"limit"`
		F64        float64   `header:"f64"`
		F32        float32   `header:"f32"`
		CreateTime time.Time `header:"createTime" time_format:"unixNano"`
		UnixTime   time.Time `header:"unixTime" time_format:"unix"`
	}

	var theader tHeader
	req := requestWithBody("GET", "/", "")
	req.Header.Add("limit", "1000")
	req.Header.Add("f64", "64")
	req.Header.Add("f32", "32.1")
	req.Header.Add("createTime", "1562400033000000123")
	req.Header.Add("unixTime", "1562400033")

	assert.NoError(t, h.Decode(req, &theader))

	assert.Equal(t, 1000, theader.Limit)
	assert.Equal(t, 64.0, theader.F64)
	assert.Equal(t, 32.1, theader.F32)
	assert.Equal(t, 1562400033000000123, theader.CreateTime)
	assert.Equal(t, 1562400033, theader.UnixTime)

	req = requestWithBody("GET", "/", "")
	req.Header.Add("fail", `{fail:fail}`)

	type failStruct struct {
		Fail map[string]interface{} `header:"fail"`
	}

	err := h.Decode(req, &failStruct{})
	assert.Error(t, err)
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return req
}
