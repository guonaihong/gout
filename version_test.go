package gout

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Version(t *testing.T) {
	assert.NotEqual(t, len(Version), 0)
}
