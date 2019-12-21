package gout

import (
	"github.com/guonaihong/gout/encode"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReqModifyUrl(t *testing.T) {
	src := []string{"127.0.0.1", ":8080/query", "/query", "http://127.0.0.1", "https://127.0.0.1"}
	want := []string{"http://127.0.0.1", "http://127.0.0.1:8080/query", "http://127.0.0.1/query", "http://127.0.0.1", "https://127.0.0.1"}

	for k, v := range src {
		if want[k] != modifyURL(v) {
			t.Errorf("got %s want %s\n", modifyURL(v), want[k])
		}
	}
}

type urlTest struct {
	set  interface{}
	need interface{}
}

func TestReq_isString(t *testing.T) {
	test := []urlTest{
		{set: "?a=b&c=d", need: "a=b&c=d"},
		{set: "c=d&e=f", need: "c=d&e=f"},
		{set: time.Time{}, need: ""},
	}

	for _, v := range test {
		rv, _ := isString(v.set)
		assert.Equal(t, v.need, rv)
	}
}

// 测试request()函数调用出错的情况
func TestReq_request_fail(t *testing.T) {

	tests := []func() *Req{
		func() *Req {
			r := Req{}
			r.bodyEncoder = encode.NewBodyEncode(&map[string]string{})
			return &r
		},
		func() *Req {
			r := Req{}
			s := "hello"
			r.formEncode = s
			return &r
		},
	}

	for _, test := range tests {
		r := test()
		_, err := r.request()
		assert.Error(t, err)
	}
}
