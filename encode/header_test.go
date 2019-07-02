package encode

import (
	"fmt"
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

type testHeader1 struct {
	H4 int64 `header:"h4"`
	H5 int32 `header:"h5"`
}

type testHeader struct {
	H1 string  `header:"h1"`
	H2 int     `header:"h2"`
	H3 float64 `header:"h3"`
	testHeader1
}

func TestHeaderStruct(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	err := Encode(testHeader{
		H1: "test-header-1",
		H2: 2,
		H3: 3.3,
		testHeader1: testHeader1{
			H4: int64(4),
			H5: int32(5),
		},
	},
		NewHeaderEnocde(req),
	)

	if err != nil {
		t.Errorf("encode http header err(struct):%s\n", err)
		return
	}

	needVal := []string{"test-header-1", "2", "3.3", "4", "5"}

	for k, v := range needVal {
		s := fmt.Sprintf("h%d", k+1)
		if v1 := req.Header.Get(s); v1 != v {
			t.Errorf("got %s want %s\n", v1, v)
		}
	}
}
