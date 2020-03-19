package core

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type testCore struct {
	need interface{}
	set  interface{}
}

func Test_Core_LoopElem(t *testing.T) {
	s := "hello world"
	p := &s
	pp := &p
	ppp := &pp

	rs := []reflect.Value{
		LoopElem(reflect.ValueOf(s)),
		LoopElem(reflect.ValueOf(p)),
		LoopElem(reflect.ValueOf(pp)),
		LoopElem(reflect.ValueOf(ppp)),
	}

	for _, v := range rs {
		assert.Equal(t, v.Interface().(string), s)
	}

	var pi *int
	assert.Equal(t, LoopElem(reflect.ValueOf(pi)), reflect.ValueOf(pi))
}

type testBytesToStr struct {
	set  []byte
	need string
}

func Test_Core_BytesToString(t *testing.T) {
	tests := []testBytesToStr{
		{set: []byte("hello world"), need: "hello world"},
		{set: []byte("hello"), need: "hello"},
	}

	for _, test := range tests {
		assert.Equal(t, BytesToString(test.set), test.need)
	}
}

func Test_Core_StringToBytes(t *testing.T) {
	tests := []testCore{
		{StringToBytes("hello"), []byte("hello")},
		{StringToBytes("world"), []byte("world")},
	}

	for _, v := range tests {
		assert.Equal(t, v.need, v.set)
	}
}

func Test_Core_NewPtrVal(t *testing.T) {
	tests := []testCore{
		{NewPtrVal(1), 1},
		{NewPtrVal(11.11), 11.11},
	}

	for _, v := range tests {
		val := reflect.ValueOf(v.need)
		val = val.Elem()
		assert.Equal(t, val.Interface(), v.set)
	}
}

func Test_Bench_closeRequest(t *testing.T) {
	b := bytes.NewBuffer([]byte("hello"))
	req, err := http.NewRequest("GET", "hello", b)
	assert.NoError(t, err)

	req.Header.Add("h1", "h2")
	req.Header.Add("h2", "h2")

	req2, err := CloneRequest(req)
	assert.NoError(t, err)

	// 测试http header是否一样
	assert.Equal(t, req.Header, req2.Header)

	b2, err := req2.GetBody()
	assert.NoError(t, err)

	b3 := bytes.NewBuffer(nil)
	io.Copy(b3, b2)

	// 测试body是否一样
	assert.Equal(t, b, b3)
}

func Test_Core_GetBytes(t *testing.T) {
	type getBytes struct {
		need []byte
		set  interface{}
	}

	tests := []getBytes{
		{[]byte("ok"), "ok"},
		{[]byte("ok"), []byte("ok")},
		{nil, new(int)},
	}

	for _, test := range tests {
		v, _ := GetBytes(test.set)
		assert.Equal(t, test.need, v)
	}
}

func Test_Core_GetString(t *testing.T) {
	type getBytes struct {
		need string
		set  interface{}
	}

	tests := []getBytes{
		{"ok", "ok"},
		{"ok", []byte("ok")},
		{"", new(int)},
	}

	for _, test := range tests {
		v, _ := GetString(test.set)
		assert.Equal(t, test.need, v)
	}
}
