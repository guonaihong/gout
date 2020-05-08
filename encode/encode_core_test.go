package encode

import (
	"errors"
	"fmt"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type encodeTest struct {
	set  interface{}
	need interface{}
}

// 测试正确的情况
func testEncodeCore_Encode(t *testing.T) {
	v := struct{ I int }{}
	p := &v
	pp := &p
	ppp := &pp
	testPtr := []encodeTest{
		{set: (*encodeTest)(nil), need: ""},
		{set: p, need: ""},
		{set: pp, need: ""},
		{set: ppp, need: ""},
		{set: core.A{}, need: ""},
	}

	req, err := http.NewRequest("GET", "127.0.0.1", nil)
	assert.NoError(t, err)

	for _, v := range testPtr {
		err := Encode(v.set, NewHeaderEncode(req))
		assert.NoError(t, err)
	}
}

type addFail struct{}

func (a *addFail) Add(key string, v reflect.Value, sf reflect.StructField) error {
	return errors.New("test use")
}

func (a *addFail) Name() string {
	return "fail"
}

func testEncodeCore_Encode_Fail(t *testing.T) {
	testFail := []encodeTest{
		{set: []string{"123"}, need: ""},
		{set: [1]string{"123"}, need: ""},
	}

	req, err := http.NewRequest("GET", "127.0.0.1", nil)
	assert.NoError(t, err)

	for _, v := range testFail {
		err := Encode(v.set, NewHeaderEncode(req))
		assert.Error(t, err)
	}

	testFail2 := []encodeTest{
		{set: core.H{"AA": "BBB"}, need: ""},
		{set: core.A{"aa", "bb"}, need: ""},
		{set: [2]string{"123", "456"}, need: ""},
		{set: 0, need: ""},
	}

	for _, v := range testFail2 {
		err := Encode(v.set, &addFail{})
		assert.Error(t, err)
	}
}

func TestEncodeCore_Encode(t *testing.T) {
	testEncodeCore_Encode(t)
	testEncodeCore_Encode_Fail(t)
}

func TestEncodeCore_valToStr(t *testing.T) {
	p := new(int)
	test := []encodeTest{
		//测试空指针
		{set: (*int)(nil), need: ""},
		//测试指针
		{set: new(int), need: ""},
		//测试双重指针
		{set: &p, need: ""},
	}

	for _, v := range test {
		assert.Equal(t, valToStr(reflect.ValueOf(v.set), emptyField), v.need)
	}
}

func TestEncodeCore_timeToStr(t *testing.T) {
	tm := time.Now()
	test := []encodeTest{
		{
			struct {
				T time.Time `time_format:"unix"`
			}{tm}, fmt.Sprint(tm.Unix()),
		},
		{
			struct {
				T time.Time `time_format:"unixNano"`
			}{tm}, fmt.Sprint(tm.UnixNano()),
		},
		{
			// time.RFC3339
			struct {
				T time.Time `time_format:"2006-01-02T15:04:05Z07:00"`
			}{tm}, tm.Format(time.RFC3339),
		},
		{
			struct {
				T time.Time
			}{tm}, tm.Format(time.RFC3339),
		},
	}

	for _, v := range test {
		val := reflect.ValueOf(v.set)

		assert.Equal(t, timeToStr(val.Field(0), val.Type().Field(0)), v.need)
	}
}

func TestEncodeCore_parseTagAndSet(t *testing.T) {
	test := []encodeTest{
		{
			struct {
				I int `header:"I,omitempty"`
			}{}, 0,
		},
	}

	for _, v := range test {
		val := reflect.ValueOf(v.set)
		err := parseTagAndSet(val.Field(0), val.Type().Field(0), NewHeaderEncode(&http.Request{}))
		assert.NoError(t, err)
	}
}

func TestEncodeCore_Contains(t *testing.T) {
	test := []encodeTest{
		{tagOptions{"json", "xml"}.Contains("xml"), true},
		{tagOptions{"json", "xml"}.Contains("xxx"), false},
	}

	for _, v := range test {
		assert.Equal(t, v.set, v.need)
	}
}

func TestEncodeCore_valueIsEmpty(t *testing.T) {
	test := []encodeTest{
		{uint(0), true},
		{uint8(0), true},
		{uint16(0), true},
		{uint32(0), true},
		{uint64(0), true},
		{int(0), true},
		{int8(0), true},
		{int16(0), true},
		{int32(0), true},
		{int64(0), true},
		{float32(0.0), true},
		{float64(0.0), true},
		{false, true},
		{[]byte{}, true},
		{"", true},
		{map[int]int{}, true},
		{[0]int{}, true},
		{interface{}(nil), true},
		{time.Time{}, true},
		{(*int)(nil), true},
		{nil, true},

		// 不是空值
		{encodeTest{}, false},
	}

	for _, v := range test {
		assert.Equal(t, valueIsEmpty(reflect.ValueOf(v.set)), v.need)
	}
}
