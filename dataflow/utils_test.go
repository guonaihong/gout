package dataflow

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtil_Join(t *testing.T) {
	urls := []string{
		"http://127.0.0.1:43471/v1",
	}

	want := []string{
		"http://127.0.0.1:43471/v1",
	}

	assert.Equal(t, joinPaths("", urls[0]), want[0])
}

type joinTest struct {
	absolutePath string
	relativePath string
	need         string
}

func TestUtil_join(t *testing.T) {
	//实际调用不会发生 join函数absolutePath和relativePath参数都为空的情况
	assert.Equal(t, join("", ""), ".")
}

func TestUtil_joinPaths(t *testing.T) {
	test := []joinTest{
		{absolutePath: "", relativePath: "", need: ""},
		{absolutePath: "https://www.aa.com", relativePath: "/a", need: "https://www.aa.com/a"},
		{absolutePath: "http://www.bb.com", relativePath: "/b", need: "http://www.bb.com/b"},
		{absolutePath: "www.bb.com", relativePath: "/b", need: "www.bb.com/b"},
		{absolutePath: "www.bb.com", relativePath: "", need: "www.bb.com"},
		{absolutePath: "", relativePath: "/a", need: "/a"},
		{absolutePath: "", relativePath: "http://", need: "http://"},
		{absolutePath: "", relativePath: "https://", need: "https://"},
		{absolutePath: "", relativePath: "http://www.aa.com/urls?", need: "http://www.aa.com/urls?"},
		{absolutePath: "", relativePath: "https://www.aa.com/urls?bb", need: "https://www.aa.com/urls?bb"},
		{absolutePath: "", relativePath: "http://www.aa.com/urls?site=https://bb.com", need: "http://www.aa.com/urls?site=https://bb.com"},
		{absolutePath: "", relativePath: "https://www.aa.com/urls?site=http://bb.com", need: "https://www.aa.com/urls?site=http://bb.com"},
		{absolutePath: "", relativePath: "http://www.aa.com/urls/./a?site=https://bb.com", need: "http://www.aa.com/urls/a?site=https://bb.com"},
		{absolutePath: "", relativePath: "https://www.aa.com/urls/../a?site=https://bb.com", need: "https://www.aa.com/a?site=https://bb.com"},
	}

	for _, v := range test {
		rv := joinPaths(v.absolutePath, v.relativePath)
		assert.Equal(t, v.need, rv)
	}
}
