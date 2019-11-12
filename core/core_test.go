package core

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

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
	assert.Equal(t, reflect.ValueOf(pi).Interface().(*int), (*int)(nil))
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

type testCore struct {
	need interface{}
	set  interface{}
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
