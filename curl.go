package gout

import (
	"github.com/guonaihong/gout/export"
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

func (c *Curl) GenAndSend() *Curl {
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

func (c *Curl) Do() (err error) {
	req, err := c.df.Req.request()
	if err != nil {
		return err
	}

	client := c.df.out.Client

	if c.generateAndSend {
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		err = c.df.bind(req, resp)
		if err != nil {
			return err
		}
	}

	w := io.Writer(os.Stdout)
	if c.w != nil {
		w = c.w
	}

	return export.GenCurl(req, c.longOption, w)
}
