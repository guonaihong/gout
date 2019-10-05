package encode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewYAMLEncode(t *testing.T) {
	y := NewYAMLEncode(nil)
	assert.Nil(t, y)
}
