package gout

import (
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
