package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Core_Util_Close(t *testing.T) {
	assert.NoError(t, (&ReadCloseFail{}).Close())
}

func Test_Core_Util_Read(t *testing.T) {
	p := make([]byte, 1)
	_, err := (&ReadCloseFail{}).Read(p)
	assert.Error(t, err)
}
