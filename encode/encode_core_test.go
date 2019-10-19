package encode

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type encodeTest struct {
	set  interface{}
	need interface{}
}

func TestEncodeCore_Encode(t *testing.T) {
	//test := []encodeTest{}
}

func TestEncodeCore_valToStr(t *testing.T) {
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
