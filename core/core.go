package core

import (
	"errors"
)

type FormFile string

type FormMem []byte

var ErrUnkownType = errors.New("unkown type")
