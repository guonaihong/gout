package color

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Color_Spurple(t *testing.T) {
	NoColor = false

	need := "\x1b[35;1mhello\x1b[0m"
	got := New(true).Spurple("hello")

	fmt.Printf("got(%s) need(%s)\n", got, need)
	assert.Equal(t, got, need)
	assert.Equal(t, New(false).Spurple("hello"), "hello")
}

func Test_Color_Sbluef(t *testing.T) {
	NoColor = false

	assert.Equal(t, New(true).Sbluef("hello"), "\x1b[34;1mhello\x1b[0m")
	assert.Equal(t, New(false).Sbluef("hello"), "hello")
}

func Test_Color_color(t *testing.T) {
	NoColor = false
	assert.Equal(t, New(true, Blue).color("hello"), "\x1b[34;1mhello\x1b[0m")
	assert.Equal(t, New(false).color("hello"), "hello")
}

func Test_Color_New(t *testing.T) {
	n := New(true)
	assert.NotNil(t, n)
}
