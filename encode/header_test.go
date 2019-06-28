package encode

import (
	"net/http"
	"testing"
)

type testH map[string]interface{}

func TestHeaderStringSlice(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	err := Encode([]string{
		"header1", "value1",
		"header2", "value2",
		"header3", "value3",
	}, NewHeaderEnocde(req))

	if err != nil {
		t.Errorf("encode http header err([]string):%s\n", err)
		return
	}

	if v := req.Header.Get("header1"); v != "value1" {
		t.Errorf("got (%s) want value1\n", v)
	}

	if v := req.Header.Get("header2"); v != "value2" {
		t.Errorf("got (%s) want value2\n", v)
	}

	if v := req.Header.Get("header3"); v != "value3" {
		t.Errorf("got (%s) want value3\n", v)
	}
}

func TestHeaderMap(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	err := Encode(testH{
		"header1": 1,
		"header2": "value2",
		"header3": 3.14,
	}, NewHeaderEnocde(req))

	if err != nil {
		t.Errorf("encode http header err(gout.H):%s\n", err)
		return
	}

	if v := req.Header.Get("header1"); v != "1" {
		t.Errorf("got %s want 1\n", v)
	}

	if v := req.Header.Get("header2"); v != "value2" {
		t.Errorf("got %s want value2\n", v)
	}

	if v := req.Header.Get("header3"); v != "3.14" {
		t.Errorf("got %s want 3.14\n", v)
	}
}
