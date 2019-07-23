package gout

import (
	"testing"
)

func TestReqModifyUrl(t *testing.T) {
	src := []string{"127.0.0.1", ":8080/query", "/query", "http://127.0.0.1", "https://127.0.0.1"}
	want := []string{"http://127.0.0.1", "http://127.0.0.1:8080/query", "http://127.0.0.1/query", "http://127.0.0.1", "https://127.0.0.1"}

	for k, v := range src {
		if want[k] != modifyUrl(v) {
			t.Errorf("got %s want %s\n", modifyUrl(v), want[k])
		}
	}
}
