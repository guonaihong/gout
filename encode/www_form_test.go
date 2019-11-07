package encode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WWWForm_New(t *testing.T) {
	assert.Nil(t, NewWWWFormEncode(nil))
	assert.NotNil(t, NewWWWFormEncode("h"))
}
