package gout

import (
	"io"
)

type Decoder interface {
	Decode(r io.Reader) error
}
