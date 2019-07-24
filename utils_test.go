package gout

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	urls := []string{
		"http://127.0.0.1:43471/v1",
	}

	want := []string{
		"http://127.0.0.1:43471/v1",
	}

	assert.Equal(t, joinPaths("", urls[0]), want[0])
}
