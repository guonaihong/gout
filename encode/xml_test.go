package encode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewXMLEncode(t *testing.T) {
	x := NewXMLEncode(nil)
	assert.Nil(t, x)
}
