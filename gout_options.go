package gout

import (
	"crypto/tls"
	"net"
	"net/http"
)

// 配置
type options struct {
	hc *http.Client
}

// apply接口
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

// 传递http.Client
func WithClient(c *http.Client) Option {
	return (*client)(c)
}

type fromIP func(netw, addr string) (net.Conn, error)

func (f fromIP) apply(opts *options) {
	tr, ok := opts.hc.Transport.(*http.Transport)
	if !ok {
		tr = &http.Transport{}
	}

	tr.Dial = f
}

// 指定发送的ip
func WithFromIP(IP string) Option {
	dial := func(network, addr string) (net.Conn, error) {
		lAddr, err := net.ResolveTCPAddr(network, IP+":0")
		if err != nil {
			return nil, err
		}

		//被请求的地址
		rAddr, err := net.ResolveTCPAddr(network, addr)
		if err != nil {
			return nil, err
		}

		return net.DialTCP(network, lAddr, rAddr)
	}

	return fromIP(dial)
}
