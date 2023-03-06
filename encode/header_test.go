package encode

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testH map[string]interface{}

func TestHeaderStringSlice(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	err := Encode([]string{
		"header1", "value1",
		"header2", "value2",
		"header3", "value3",
	}, NewHeaderEncode(req, false))

	assert.NoError(t, err)

	needVal := []string{"value1", "value2", "value3"}

	for k, v := range needVal {
		s := fmt.Sprintf("header%d", k+1)
		assert.Equal(t, req.Header.Get(s), v)
	}
}

func TestHeaderMap(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	err := Encode(testH{
		"header1": 1,
		"header2": "value2",
		"header3": 3.14,
	}, NewHeaderEncode(req, false))

	assert.NoError(t, err)

	needVal := []string{"1", "value2", "3.14"}

	for k, v := range needVal {
		s := fmt.Sprintf("header%d", k+1)
		assert.Equal(t, req.Header.Get(s), v)
	}
}

type testHeader1 struct {
	H4 int64  `header:"h4"`
	H5 int32  `header:"h5"`
	H6 string `header:"-"`
}

type testHeader2 struct {
	H7 string `header:"h7"`
}

type testHeader struct {
	H1 string  `header:"h1"`
	H2 int     `header:"h2"`
	H3 float64 `header:"h3"`
	testHeader1

	H  **testHeader2 // 测试多重指针
	H8 *testHeader2  //测试结构体空指针
}

func TestHeaderStruct(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	p := &testHeader2{H7: "h7"}

	err := Encode(testHeader{
		H1: "test-header-1",
		H2: 2,
		H3: 3.3,
		testHeader1: testHeader1{
			H4: int64(4),
			H5: int32(5),
		},
		H: &p,
	},
		NewHeaderEncode(req, false),
	)

	assert.NoError(t, err)

	needVal := []string{"test-header-1", "2", "3.3", "4", "5", "", "h7"}

	for k, v := range needVal {
		s := fmt.Sprintf("h%d", k+1)
		assert.Equal(t, req.Header.Get(s), v)
	}
}
