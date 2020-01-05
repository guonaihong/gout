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

func (c *Curl) SetOutput(w io.Writer) *Curl {
	c.w = w
	return c
}

func (c *Curl) Do() (err error) {
	if c.w == nil {
		c.w = os.Stdout
	}

	w := c.w

	req, err := c.df.Req.request()
	if err != nil {
		return err
	}

	client := c.df.out.Client

	if c.generateAndSend {
		// 清空状态，Setxxx函数拆开使用就不会有问题
		defer c.df.Reset()
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		err = c.df.bind(req, resp)
		if err != nil {
			return err
		}
	}

	return export.GenCurl(req, c.longOption, w)
}
