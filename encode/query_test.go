package encode

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
	"time"
)

//todo
func TestQueryString(t *testing.T) {
}

func TestQueryStringSlice(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	q := NewQueryEncode(req)

	err := Encode([]string{"q1", "v1", "q2", "v2", "q3", "v3"}, q)

	assert.NoError(t, err)

	assert.Equal(t, q.End(), "q1=v1&q2=v2&q3=v3")
}

func TestQueryMap(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	q := NewQueryEncode(req)

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

func TestQueryStruct(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	q := NewQueryEncode(req)

	unixTime := time.Date(2019, 07, 27, 20, 42, 53, 0, time.Local)
	unixNano := time.Date(2019, 07, 27, 20, 42, 53, 1000, time.Local)

	err := Encode(testQuery{Q1: "v1", Q2: "v2", Q3: "v3", Q4: unixTime, Q5: unixNano}, q)

	assert.NoError(t, err)

	assert.Equal(t, q.End(), "q1=v1&q2=v2&q3=v3&q4="+strconv.Itoa(int(unixTime.Unix()))+"&q5="+strconv.Itoa(int(unixNano.UnixNano())))
}
