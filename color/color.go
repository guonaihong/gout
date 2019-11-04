package color

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
	"strings"
)

var (
	// 本行代码来自github.com/fatih/color, 感谢fatih的付出
	NoColor = os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
)

type attr int

const (
	FgBlack attr = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

const (
	Gray = 30
	Blue = 34
)

type Color struct {
	openColor bool
	attr      attr
}

func New(openColor bool, c ...attr) *Color {
	attr := attr(30)
	if len(c) > 0 {
		attr = c[0]
	}
	return &Color{openColor: openColor, attr: attr}
}

func (c *Color) set(buf *strings.Builder, attr attr) {
	if NoColor || !c.openColor {
		return
	}

	fmt.Fprintf(buf, "\x1b[%d;1m", attr)
}

func (c *Color) unset(buf *strings.Builder) {
	if NoColor || !c.openColor {
		return
	}

	fmt.Fprintf(buf, "\x1b[0m")
}

func (c *Color) color(a ...interface{}) string {
	var buf strings.Builder

	c.set(&buf, c.attr)

	fmt.Fprint(&buf, a...)
	c.unset(&buf)

	return buf.String()
}

func (c *Color) colorf(format string, a ...interface{}) string {
	var buf strings.Builder

	c.set(&buf, c.attr)

	fmt.Fprintf(&buf, format, a...)
	c.unset(&buf)

	return buf.String()
}

func (c *Color) Sgrayf(format string, a ...interface{}) string {
	c.attr = Gray
	return c.colorf(format, a...)
}

func (c *Color) Sbluef(format string, a ...interface{}) string {
	c.attr = Blue
	return c.colorf(format, a...)
}
