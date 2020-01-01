package gout

import (
	"io"
	"os"
)

type Curl struct {
	w               io.Writer
	df              *DataFlow
	longOption      bool
	generateAndSend bool
}

func (c *Curl) LongOption() *Curl {
	c.longOption = true
	return c
}

func (c *Curl) GenerateAndSend() *Curl {
	c.generateAndSend = true
	return c
}

func (c *Curl) SetOutput(w ...io.Writer) *Curl {
	if len(w) == 0 {
		c.w = os.Stdout
		return c
	}

	c.w = w[0]
	return c
}

func (c *Curl) Do() error {
	return nil
}
