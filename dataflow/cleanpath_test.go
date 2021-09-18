package dataflow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cleanPathTest struct {
	path string
	need string
}

func TestUtil_cleanPaths(t *testing.T) {
	test := []cleanPathTest{
		{path: "", need: ""},
		{path: "https://www.aa.com/a", need: "https://www.aa.com/a"},
		{path: "http://www.bb.com/b", need: "http://www.bb.com/b"},
		{path: "www.bb.com/b", need: "www.bb.com/b"},
		{path: "www.bb.com", need: "www.bb.com"},
		{path: "/a", need: "/a"},
		{path: "/a/", need: "/a/"},
		{path: "http://", need: "http://"},
		{path: "https://", need: "https://"},
		{path: "http://www.aa.com/urls?", need: "http://www.aa.com/urls?"},
		{path: "https://www.aa.com/urls?bb", need: "https://www.aa.com/urls?bb"},
		{path: "http://www.aa.com/urls?site=https://bb.com", need: "http://www.aa.com/urls?site=https://bb.com"},
		{path: "https://www.aa.com/urls?site=http://bb.com", need: "https://www.aa.com/urls?site=http://bb.com"},
		{path: "http://www.aa.com/urls/./a?site=https://bb.com", need: "http://www.aa.com/urls/a?site=https://bb.com"},
		{path: "https://www.aa.com/urls/../a?site=https://bb.com", need: "https://www.aa.com/a?site=https://bb.com"},
		{path: "https://api.map.baidu.com/weather/v1/?district_id=310100&data_type=all&ak=ffyu0pP8P6Ao0KYr8FTZwDgsOFiA1oYA", need: "https://api.map.baidu.com/weather/v1/?district_id=310100&data_type=all&ak=ffyu0pP8P6Ao0KYr8FTZwDgsOFiA1oYA"},
	}

	for index, v := range test {
		rv := cleanPaths(v.path)
		assert.Equal(t, v.need, rv, fmt.Sprintf("index:%d", index))
	}
}
