package gout

import (
	"crypto/tls"
	"net/http"
)

type options struct {
	hc *http.Client
}

type Option interface {
	apply(*options)
}

type insecureSkipVerifyOption bool

func (i insecureSkipVerifyOption) apply(opts *options) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	opts.hc.Transport = tr
}

// 忽略ssl验证
func WithInsecureSkipVerify() Option {
	b := true
	return insecureSkipVerifyOption(b)
}

type client http.Client

func (c *client) apply(opts *options) {
	opts.hc = (*http.Client)(c)
}

func WithClient(c *http.Client) Option {
	return (*client)(c)
}
