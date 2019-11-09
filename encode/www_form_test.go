package encode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WWWForm_New(t *testing.T) {
	assert.Nil(t, NewWWWFormEncode(nil))
	assert.NotNil(t, NewWWWFormEncode("h"))
}

func Test_WWWForm_Encode(t *testing.T) {
}

func Test_WWWForm_Add(t *testing.T) {
}

func Test_WWWForm_Name(t *testing.T) {
	assert.Equal(t, NewWWWFormEncode("").Name(), "www-form")
}
