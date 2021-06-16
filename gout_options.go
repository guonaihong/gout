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

// 1.start
type insecureSkipVerifyOption bool

func (i insecureSkipVerifyOption) apply(opts *options) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	opts.hc.Transport = tr
}

// 1.忽略ssl验证
func WithInsecureSkipVerify() Option {
	b := true
	return insecureSkipVerifyOption(b)
}

// 2. start
type client http.Client

func (c *client) apply(opts *options) {
	opts.hc = (*http.Client)(c)
}

// 2.自定义http.Client
func WithClient(c *http.Client) Option {
	return (*client)(c)
}

// 3.start
type close3xx struct{}

func (c close3xx) apply(opts *options) {
	opts.hc.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

// 3.关闭3xx自动跳转
func WithClose3xxJump() Option {
	return close3xx{}
}
