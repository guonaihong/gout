package encode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJSONEncode(t *testing.T) {
	j := NewJSONEncode(nil)
	assert.Nil(t, j)
}
