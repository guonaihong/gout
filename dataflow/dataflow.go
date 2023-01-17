package dataflow

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/guonaihong/gout/debug"
	"github.com/guonaihong/gout/decode"
	"github.com/guonaihong/gout/encode"
	"github.com/guonaihong/gout/enjson"
	"github.com/guonaihong/gout/middler"
	"github.com/guonaihong/gout/middleware/rsp/autodecodebody"
	"github.com/guonaihong/gout/setting"
	"golang.org/x/net/proxy"
)

const (
	get     = "GET"
	post    = "POST"
	put     = "PUT"
	delete2 = "DELETE"
	patch   = "PATCH"
	head    = "HEAD"
	options = "OPTIONS"
)

// DataFlow is the core data structure,
// including the encoder and decoder of http data
type DataFlow struct {
	Req
	out *Gout
}

// GET send HTTP GET method
func (df *DataFlow) GET(url string, urlStruct ...interface{}) *DataFlow {
	df.Req = reqDef(get, cleanPaths(url), df.out, urlStruct...)
	return df
}

// POST send HTTP POST method
func (df *DataFlow) POST(url string, urlStruct ...interface{}) *DataFlow {
	df.Req = reqDef(post, cleanPaths(url), df.out, urlStruct...)
	return df
}

// PUT send HTTP PUT method
func (df *DataFlow) PUT(url string, urlStruct ...interface{}) *DataFlow {
	df.Req = reqDef(put, cleanPaths(url), df.out, urlStruct...)
	return df
}

// DELETE send HTTP DELETE method
func (df *DataFlow) DELETE(url string, urlStruct ...interface{}) *DataFlow {
	df.Req = reqDef(delete2, cleanPaths(url), df.out, urlStruct...)
	return df
}

// PATCH send HTTP PATCH method
func (df *DataFlow) PATCH(url string, urlStruct ...interface{}) *DataFlow {
	df.Req = reqDef(patch, cleanPaths(url), df.out, urlStruct...)
	return df
}

// HEAD send HTTP HEAD method
func (df *DataFlow) HEAD(url string, urlStruct ...interface{}) *DataFlow {
	df.Req = reqDef(head, cleanPaths(url), df.out, urlStruct...)
	return df
}

// OPTIONS send HTTP OPTIONS method
func (df *DataFlow) OPTIONS(url string, urlStruct ...interface{}) *DataFlow {
	df.Req = reqDef(options, cleanPaths(url), df.out, urlStruct...)
	return df
}

// SetSetting
func (df *DataFlow) SetSetting(s setting.Setting) *DataFlow {
	df.Req.Setting = s
	return df
}

// SetHost set host
func (df *DataFlow) SetHost(host string) *DataFlow {
	if df.Err != nil {
		return df
	}

	df.Req.host = host

	df.Req.g = df.out
	return df
}

// GetHost return value->host or host:port
func (df *DataFlow) GetHost() (string, error) {
	if df.req != nil {
		return df.req.URL.Host, nil
	}

	if len(df.host) > 0 {
		return df.host, nil
	}

	if len(df.url) > 0 {
		url, err := url.Parse(df.url)
		if err != nil {
			return "", err
		}
		return url.Host, nil
	}

	return "", errors.New("not url found")
}

// SetMethod set method
func (df *DataFlow) SetMethod(method string) *DataFlow {
	if df.Err != nil {
		return df
	}

	df.Req.method = method
	df.Req.g = df.out
	return df
}

// SetURL set url
func (df *DataFlow) SetURL(url string, urlStruct ...interface{}) *DataFlow {
	if df.Err != nil {
		return df
	}

	if df.Req.url == "" && df.Req.req == nil {
		df.Req = reqDef(df.method, cleanPaths(url), df.out, urlStruct...)
		return df
	}

	df.Req.url = modifyURL(cleanPaths(url))

	return df
}

func (df *DataFlow) SetRequest(req *http.Request) *DataFlow {
	df.req = req
	return df
}

// SetBody set the data to the http body, Support string/bytes/io.Reader
func (df *DataFlow) SetBody(obj interface{}) *DataFlow {

	df.Req.bodyEncoder = encode.NewBodyEncode(obj)
	return df
}

// SetForm send form data to the http body, Support struct/map/array/slice
func (df *DataFlow) SetForm(obj ...interface{}) *DataFlow {
	df.Req.form = append([]interface{}{}, obj...)
	return df
}

// SetWWWForm send x-www-form-urlencoded to the http body, Support struct/map/array/slice types
func (df *DataFlow) SetWWWForm(obj ...interface{}) *DataFlow {
	df.Req.wwwForm = append([]interface{}{}, obj...)
	return df
}

// SetQuery send URL query string, Support string/[]byte/struct/map/slice types
func (df *DataFlow) SetQuery(obj ...interface{}) *DataFlow {
	df.Req.queryEncode = append([]interface{}{}, obj...)
	return df
}

// SetHeader send http header, Support struct/map/slice types
func (df *DataFlow) SetHeader(obj ...interface{}) *DataFlow {
	df.Req.headerEncode = append([]interface{}{}, obj...)
	return df
}

// SetJSON send json to the http body, Support raw json(string, []byte)/struct/map types
func (df *DataFlow) SetJSON(obj interface{}) *DataFlow {
	df.ReqBodyType = "json"
	df.EscapeHTML = true
	df.Req.bodyEncoder = enjson.NewJSONEncode(obj, true)
	return df
}

// SetJSON send json to the http body, Support raw json(string, []byte)/struct/map types
// 与SetJSON的区一区别就是不转义HTML里面的标签
func (df *DataFlow) SetJSONNotEscape(obj interface{}) *DataFlow {
	df.ReqBodyType = "json"
	df.Req.bodyEncoder = enjson.NewJSONEncode(obj, false)
	return df
}

// SetXML send xml to the http body
func (df *DataFlow) SetXML(obj interface{}) *DataFlow {
	df.ReqBodyType = "xml"
	df.Req.bodyEncoder = encode.NewXMLEncode(obj)
	return df
}

// SetYAML send yaml to the http body, Support struct,map types
func (df *DataFlow) SetYAML(obj interface{}) *DataFlow {
	df.ReqBodyType = "yaml"
	df.Req.bodyEncoder = encode.NewYAMLEncode(obj)
	return df
}

// SetProtoBuf send yaml to the http body, Support struct types
// obj必须是结构体指针或者[]byte类型
func (df *DataFlow) SetProtoBuf(obj interface{}) *DataFlow {
	df.ReqBodyType = "protobuf"
	df.Req.bodyEncoder = encode.NewProtoBufEncode(obj)
	return df
}

func (df *DataFlow) initTransport() {
	if df.out.Client.Transport == nil {
		df.out.Client.Transport = &http.Transport{}
	}
}

func (df *DataFlow) getTransport() (*http.Transport, bool) {
	// 直接return df.out.Client.Transport.(*http.Transport) 等于下面的写法
	// ts := df.out.Client.Transport.(*http.Transport)
	// return ts 编译会报错
	ts, ok := df.out.Client.Transport.(*http.Transport)
	return ts, ok
}

// UnixSocket 函数会修改Transport, 请像对待全局变量一样对待UnixSocket
func (df *DataFlow) UnixSocket(path string) *DataFlow {

	df.initTransport()

	transport, ok := df.getTransport()
	if !ok {
		df.Req.Err = fmt.Errorf("UnixSocket:not found http.transport:%T", df.out.Client.Transport)
		return df
	}

	transport.Dial = func(proto, addr string) (conn net.Conn, err error) {
		return net.Dial("unix", path)
	}

	return df
}

// SetProxy 函数会修改Transport，请像对待全局变量一样对待SetProxy
func (df *DataFlow) SetProxy(proxyURL string) *DataFlow {
	proxy, err := url.Parse(modifyURL(proxyURL))
	if err != nil {
		df.Req.Err = err
		return df
	}

	df.initTransport()

	transport, ok := df.getTransport()
	if !ok {
		df.Req.Err = fmt.Errorf("SetProxy:not found http.transport:%T", df.out.Client.Transport)
		return df
	}

	transport.Proxy = http.ProxyURL(proxy)

	return df
}

// SetSOCKS5 函数会修改Transport,请像对待全局变量一样对待SetSOCKS5
func (df *DataFlow) SetSOCKS5(addr string) *DataFlow {
	dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	if err != nil {
		df.Req.Err = err
		return df
	}

	df.initTransport()

	transport, ok := df.getTransport()
	if !ok {
		df.Req.Err = fmt.Errorf("SetSOCKS5:not found http.transport:%T", df.out.Client.Transport)
		return df
	}

	transport.Dial = dialer.Dial
	return df
}

// SetCookies set cookies
func (df *DataFlow) SetCookies(c ...*http.Cookie) *DataFlow {
	df.Req.cookies = append(df.Req.cookies, c...)
	return df
}

// BindHeader parse http header to obj variable.
// obj must be a pointer variable
// Support string/int/float/slice ... types
func (df *DataFlow) BindHeader(obj interface{}) *DataFlow {
	df.Req.headerDecode = obj
	return df
}

// BindBody parse the variables in http body to obj.
// obj must be a pointer variable
func (df *DataFlow) BindBody(obj interface{}) *DataFlow {
	if obj == nil {
		return df
	}
	df.Req.bodyDecoder = append(df.Req.bodyDecoder, decode.NewBodyDecode(obj))
	return df
}

// BindJSON parse the json string in http body to obj.
// obj must be a pointer variable
func (df *DataFlow) BindJSON(obj interface{}) *DataFlow {
	if obj == nil {
		return df
	}
	df.out.RspBodyType = "json"
	df.Req.bodyDecoder = append(df.Req.bodyDecoder, decode.NewJSONDecode(obj))
	return df

}

// BindYAML parse the yaml string in http body to obj.
// obj must be a pointer variable
func (df *DataFlow) BindYAML(obj interface{}) *DataFlow {
	if obj == nil {
		return df
	}
	df.RspBodyType = "yaml"
	df.Req.bodyDecoder = append(df.Req.bodyDecoder, decode.NewYAMLDecode(obj))
	return df
}

// BindXML parse the xml string in http body to obj.
// obj must be a pointer variable
func (df *DataFlow) BindXML(obj interface{}) *DataFlow {
	if obj == nil {
		return df
	}
	df.RspBodyType = "xml"
	df.Req.bodyDecoder = append(df.Req.bodyDecoder, decode.NewXMLDecode(obj))
	return df
}

// Code parse the http code into the variable httpCode
func (df *DataFlow) Code(httpCode *int) *DataFlow {
	df.Req.httpCode = httpCode
	return df
}

// Callback parse the http body into obj according to the condition (json or string)
func (df *DataFlow) Callback(cb func(*Context) error) *DataFlow {
	df.Req.callback = cb
	return df
}

// Chunked
func (df *DataFlow) Chunked() *DataFlow {
	df.Req.Chunked()
	return df
}

// SetTimeout set timeout, and WithContext are mutually exclusive functions
func (df *DataFlow) SetTimeout(d time.Duration) *DataFlow {
	df.Req.SetTimeout(d)
	return df
}

// WithContext set context, and SetTimeout are mutually exclusive functions
func (df *DataFlow) WithContext(c context.Context) *DataFlow {
	df.Req.Index++
	df.Req.ctxIndex = df.Req.Index
	df.Req.c = c
	return df
}

// SetBasicAuth
func (df *DataFlow) SetBasicAuth(username, password string) *DataFlow {
	df.Req.userName = &username
	df.Req.password = &password
	return df
}

// Request middleware
func (df *DataFlow) RequestUse(reqModify ...middler.RequestMiddler) *DataFlow {
	if len(reqModify) > 0 {
		df.reqModify = append(df.reqModify, reqModify...)
	}
	return df
}

// Response middleware
func (df *DataFlow) ResponseUse(responseModify ...middler.ResponseMiddler) *DataFlow {
	if len(responseModify) > 0 {
		df.responseModify = append(df.responseModify, responseModify...)
	}
	return df
}

// Debug start debug mode
func (df *DataFlow) Debug(d ...interface{}) *DataFlow {
	for _, v := range d {
		switch opt := v.(type) {
		case bool:
			if opt {
				debug.DefaultDebug(&df.Setting.Options)
			}
		case debug.Apply:
			opt.Apply(&df.Setting.Options)
		}
	}

	return df
}

// https://github.com/guonaihong/gout/issues/264
// When calling SetWWWForm(), the Content-Type header will be added automatically,
// and calling NoAutoContentType() will not add an HTTP header
//
// SetWWWForm "Content-Type", "application/x-www-form-urlencoded"
// SetJSON "Content-Type", "application/json"
func (df *DataFlow) NoAutoContentType() *DataFlow {
	df.Req.NoAutoContentType = true
	return df
}

// https://github.com/guonaihong/gout/issues/343
// content-encoding会指定response body的压缩方法，支持常用的压缩，gzip, deflate, br等
func (df *DataFlow) AutoDecodeBody() *DataFlow {
	return df.ResponseUse(middler.WithResponseMiddlerFunc(autodecodebody.AutoDecodeBody))
}

func (df *DataFlow) IsDebug() bool {
	return df.Setting.Debug
}

// Do send function
func (df *DataFlow) Do() (err error) {
	return df.Req.Do()
}

// Filter filter function, use this function to turn on the filter function
func (df *DataFlow) Filter() *filter {
	return &filter{df: df}
}

// F filter function, use this function to turn on the filter function
func (df *DataFlow) F() *filter {
	return df.Filter()
}

// Export filter function, use this function to turn on the filter function
func (df *DataFlow) Export() *export {
	return &export{df: df}
}

// E filter function, use this function to turn on the filter function
func (df *DataFlow) E() *export {
	return df.Export()
}

func (df *DataFlow) SetGout(out *Gout) {
	df.out = out
	df.Req.g = out
}
