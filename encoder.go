package gout

import (
	"io"
)

type Encoder interface {
	Encode(w io.Writer) error
	Name() string
}
