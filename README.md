# gout
gout 是go写的http 客户端，为提高工作效率而开发

[![Build Status](https://travis-ci.org/guonaihong/gout.png)](https://travis-ci.org/guonaihong/gout)
[![codecov](https://codecov.io/gh/guonaihong/gout/branch/master/graph/badge.svg)](https://codecov.io/gh/guonaihong/gout)

## feature
* 支持设置 GET/PUT/DELETE/PATH/HEAD/OPTIONS
* 支持设置请求 http header(可传 struct,map,array,slice 等类型)
* 支持设置 URL query(可传 struct,map,array,slice,string 等类型)
* 支持设置 json,xml,yaml 编码到请求 body 里面(SetJSON/SetXML/SetYAML)
* 支持设置 form-data(可传 struct,map,array,slice 等类型)
* 支持设置 x-www-form-urlencoded(可传 struct,map,array,slice 等类型) 
* 支持设置 io.Reader，uint/uint8/uint16...int/int8...string...[]byte...float32,float64 至请求 body 里面
* 支持解析响应body里面的json,xml,yaml至结构体里(BindJSON/BindXML/BindYAML)
* 支持解析响应body的内容至io.Writer, uint/uint8...int/int8...string...[]byte...float32,float64
* 支持解析响应header至结构体里
* 支持接口性能benchmark，可控制压测一定次数还是时间，可控制压测频率
* 等等更多

## 演示
![gout-example.gif](https://github.com/guonaihong/images/blob/master/gout/gout-example.gif?raw=true)

## 内容
- [Installation](#Installation)
- [Migrate documents](#Migrate-documents)
- [example](#example)
- [quick start](#quick-start)
- [API Examples](#api-examples)
    - [GET POST PUT DELETE PATH HEAD OPTIONS](#get-post-put-delete-path-head-options)
    - [query](#query)
    - [http header](#http-header)
		- [req header](#req-header)
		- [rsp header](#rsp-header)
    - [http body](#http-body)
        - [body](#body)
            - [SetBody](#setbody)
            - [BindBody](#bindbody)
        - [json](#json)
            - [SetJSON](#setjson)
            - [BindJSON](#bindjson)
        - [yaml](#yaml)
        - [xml](#xml)
        - [form-data](#form-data)
        - [x-www-form-urlencoded](#x-www-form-urlencoded)
        - [callback](#callback)
	- [benchmark](#benchmark)
		- [number](#number)
		- [duration](#duration)
		- [rate](#rate)
	- [timeout](#timeout)
    - [proxy](#proxy)
    - [cookie](#cookie)
    - [context](#context)
        - [cancel](#cancel)
    - [unix socket](#unix-socket)
    - [http2 doc](#http2-doc)
    - [debug mode](#debug-mode)
        - [color](#color)
        - [customize](#customize)
        - [no-color](#no-color)
 - [特色功能示例](#特色功能示例)
    - [forward gin data](#forward-gin-data)

- [FAQ](#FAQ)

## Installation
```
env GOPATH=`pwd` go get github.com/guonaihong/gout
```

## Migrate documents
主要方便下面的用户迁移到gout
* [httplib](./to-gout-doc/beego-httplib.md)
* [resty](./to-gout-doc/resty-doc.md)
# example
 [examples](./_example) 目录下面的例子，都是可以直接跑的。如果觉得运行例子还是不明白用法，可以把你迷惑的地方写出来，然后提[issue](https://github.com/guonaihong/gout/issues/new)
 ### 运行命令如下
 ```bash
 cd _example
 # GOPROXY 是打开go module代理，可以更快下载模块
 # 第一次运行需要加GOPROXY下载模块，模块已的直接 go run 01-color-json.go 即可
 env GOPROXY=https://goproxy.cn go run 01-color-json.go
 ```

 # quick start
 ```go
 package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"time"
)

// 用于解析 服务端 返回的http body
type RspBody struct {
	ErrMsg  string `json:"errmsg"`
	ErrCode int    `json:"errcode"`
	Data    string `json:"data"`
}

// 用于解析 服务端 返回的http header
type RspHeader struct {
	Sid  string `header:"sid"`
	Time int    `header:"time"`
}

func main() {
	rsp := RspBody{}
	header := RspHeader{}

	//code := 0
	err := gout.

		// POST请求
		POST("127.0.0.1:8080").

		// 打开debug模式
		Debug(true).

		// 设置查询字符串
		SetQuery(gout.H{"page": 10, "size": 10}).

		// 设置http header
		SetHeader(gout.H{"X-IP": "127.0.0.1", "sid": fmt.Sprintf("%x", time.Now().UnixNano())}).

		// SetJSON设置http body为json
		// 同类函数有SetBody, SetYAML, SetXML, SetForm, SetWWWForm
		SetJSON(gout.H{"text": "gout"}).

		// BindJSON解析返回的body内容
		// 同类函数有BindBody, BindYAML, BindXML
		BindJSON(&rsp).

		// 解析返回的http header
		BindHeader(&header).
		// http code
		// Code(&code).

		// 结束函数
		Do()

		// 判度错误
	if err != nil {
		fmt.Printf("send fail:%s\n", err)
	}
}

/*
> POST /?page=10&size=10 HTTP/1.1
> Sid: 15d9b742ef32c130
> X-Ip: 127.0.0.1
> Content-Type: application/json
>

{
    "text": "gout"
}


*/
 ```
# API examples
## GET POST PUT DELETE PATH HEAD OPTIONS
```go
package main

import (
	"github.com/guonaihong/gout"
)

func main() {
	url := "https://github.com"
	// 发送GET方法
	gout.GET(url).Do()

	// 发送POST方法
	gout.POST(url).Do()

	// 发送PUT方法
	gout.PUT(url).Do()

	// 发送DELETE方法
	gout.DELETE(url).Do()

	// 发送PATH方法
	gout.PATCH(url).Do()

	// 发送HEAD方法
	gout.HEAD(url).Do()

	// 发送OPTIONS
	gout.OPTIONS(url).Do()
}

```
## query

### SetQuery
```go
package main

import (
    "fmt"
    "github.com/guonaihong/gout"
    "time"
)

func main() {
    err := gout.
        //设置GET请求和url，:8080/test.query是127.0.0.1:8080/test.query的简写
        GET(":8080/test.query").
        //打开debug模式
        Debug(true).
        //设置查询字符串
        SetQuery(gout.H{
            "q1": "v1",
            "q2": 2,
            "q3": float32(3.14),
            "q4": 4.56,
            "q5": time.Now().Unix(),
            "q6": time.Now().UnixNano(),
            "q7": time.Now().Format("2006-01-02")}).
        //结束函数
        Do()
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

}

/*
> GET /test.query?q1=v1&q2=2&q3=3.14&q4=4.56&q5=1574081600&q6=1574081600258009213&q7=2019-11-18 HTTP/1.1
>

< HTTP/1.1 200 OK
< Content-Length: 0
*/


```
### SetQuery支持的更多数据类型
```go
package main

import (
	"github.com/guonaihong/gout"
)

func main() {

	code := 0

	err := gout.

		//发送GET请求 :8080/testquery是127.0.0.1:8080/testquery简写
		GET(":8080/testquery").

		// 设置查询字符串
		SetQuery( /*看下面支持的情况*/ ).

		//解析http code，如不关心服务端返回状态吗，不设置该函数即可
		Code(&code).
		Do()
	if err != nil {

	}
}



/*
SetQuery支持的类型有
* string
* map[string]interface{}，可以使用gout.H别名
* struct
* array, slice(长度必须是偶数)
*/

// 1.string
SetQuery("check_in=2019-06-18&check_out=2018-06-18")

// 2.gout.H 或者 map[string]interface{}
SetQuery(gout.H{
    "check_in":"2019-06-18",
    "check_out":"2019-06-18",
})

// 3.struct
type testQuery struct {
    CheckIn string `query:checkin`
    CheckOut string `query:checkout`
}

SetQuery(&testQuery{CheckIn:2019-06-18, CheckOut:2019-06-18})

// 4.array or slice
// ?active=enable&action=drop
SetQuery([]string{"active", "enable", "action", "drop"})`
```

## http header
#### req header
```go
package main

import (
    "fmt"
    "github.com/guonaihong/gout"
    "time"
)

func main() {
    err := gout.
        //设置GET请求和url，:8080/test.header是127.0.0.1:8080/test.header的简写
        GET(":8080/test.header").
        //设置debug模式
        Debug(true).
        //设置请求http header
        SetHeader(gout.H{
            "h1": "v1",
            "h2": 2,
            "h3": float32(3.14),
            "h4": 4.56,
            "h5": time.Now().Unix(),
            "h6": time.Now().UnixNano(),
            "h7": time.Now().Format("2006-01-02")}).
        Do()
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

}

/*
> GET /test.header HTTP/1.1
> H2: 2
> H3: 3.14
> H4: 4.56
> H5: 1574081686
> H6: 1574081686471347098
> H7: 2019-11-18
> H1: v1
>


< HTTP/1.1 200 OK
< Content-Length: 0
*/
```
#### rsp header
```go
package main

import (
    "fmt"
    "github.com/guonaihong/gout"
    "time"
)

// 和解析json类似，如要解析http header需设置header tag
type rspHeader struct {
    Total int       `header:"total"`
    Sid   string    `header:"sid"`
    Time  time.Time `header:"time" time_format:"2006-01-02"`
}

func main() {

    rsp := rspHeader{}
    err := gout.
        // :8080/test.header是 http://127.0.0.1:8080/test.header的简写
        GET(":8080/test.header").
        //打开debug模式
        Debug(true).
        //解析请求header至结构体中
        BindHeader(&rsp). 
        //结束函数
        Do()
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    fmt.Printf("rsp header:\n%#v \nTime:%s\n", rsp, rsp.Time)
}

/*
> GET /test.header HTTP/1.1
>



< HTTP/1.1 200 OK
< Content-Length: 0
< Sid: 1234
< Time: 2019-11-18
< Total: 2048
*/

```
### SetHeader和BindHeader支持的更多类型
```go
package main

import (
    "fmt"
    "github.com/guonaihong/gout"
)

type testHeader struct {
    CheckIn  string `header:checkin`
    CheckOut string `header:checkout`
}

func main() {

    t := testHeader{}

    code := 0

    err := gout.
        GET(":8080/testquery").
        Code(&code).
        SetHeader( /*看下面支持的类型*/ ).
        BindHeader(&t).
        Do()
    if err != nil {
        fmt.Printf("fail:%s\n", err)
    }   
}

```
* BindHeader支持的类型有
结构体
```go
// struct
type testHeader struct {
    CheckIn string `header:checkin`
    CheckOut string `header:checkout`
}
```

* SetHeader支持的类型有
```go
/*
map[string]interface{}，可以使用gout.H别名
struct
array, slice(长度必须是偶数)
*/

// gout.H 或者 map[string]interface{}
SetHeader(gout.H{
    "check_in":"2019-06-18",
    "check_out":"2019-06-18",
})

// struct
type testHeader struct {
    CheckIn string `header:checkin`
    CheckOut string `header:checkout`
}

SetHeader(&testHeader{CheckIn:2019-06-18, CheckOut:2019-06-18})

// array or slice
// -H active:enable -H action:drop
SetHeader([]string{"active", "enable", "action", "drop"})
```

## http body
### body
#### SetBody
```go
// SetBody 设置string, []byte等类型数据到http body里面
// SetBody支持的更多数据类型可看下面
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

func main() {
	err := gout.
		// 设置POST方法和url
		POST(":8080/req/body").
		//打开debug模式
		Debug(true).
		// 设置非结构化数据到http body里面
		// 设置json需使用SetJSON
		SetBody("send string").
		//结束函数
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}

/*
> POST /req/body HTTP/1.1
>

send string

< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Content-Length: 2

*/

```
#### bindBody
```go
// BindBody bind body到string, []byte等类型变量里面
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

func main() {
	s := ""
	err := gout.
		// 设置GET 方法及url
		GET("www.baidu.com").
		// 绑定返回值
		BindBody(&s).
		// 结束函数
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("html size = %d\n", len(s))
}

```
#### 支持的类型有
* io.Reader(SetBody 支持)
* io.Writer(BindBody 支持)
* int, int8, int16, int32, int64
* uint, uint8, uint16, uint32, uint64
* string
* []byte
* float32, float64

#### 明确不支持的类型有
* struct
* array, slice

### json
#### setjson
```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

func main() {
	err := gout.POST(":8080/colorjson").
		//打开debug模式
		Debug(true).
		//设置json到请求body
		SetJSON(
			gout.H{
				"str":   "foo",
				"num":   100,
				"bool":  false,
				"null":  nil,
				"array": gout.A{"foo", "bar", "baz"},
				"obj":   gout.H{"a": 1, "b": 2},
			},
		).
		Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

/*
> POST /colorjson HTTP/1.1
> Content-Type: application/json
>

{
    "array": [
        "foo",
        "bar",
        "baz"
    ],
    "bool": false,
    "null": null,
    "num": 100,
    "obj": {
        "a": 1,
        "b": 2
    },
    "str": "foo"
}
*/

```
#### bindjson
```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

type rsp struct {
	ErrMsg  string `json:"errmsg"`
	ErrCode int    `json:"errcode"`
}

func main() {
	rsp := rsp{}
	err := gout.
		GET(":8080/colorjson").
		//打开debug模式
		Debug(true).
		//绑定响应json数据到结构体
		BindJSON(&rsp).
		//结束函数
		Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

``` 
### yaml
* SetYAML() 设置请求http body为yaml
* BindYAML() 解析响应http body里面的yaml到结构体里面

发送yaml到服务端，然后把服务端返回的yaml结果解析到结构体里面
```go
type data struct {
    Id int `yaml:"id"`
    Data string `yaml:"data"`
}


var d1, d2 data
var httpCode int 


err := gout.POST(":8080/test.yaml").SetYAML(&d1).BindYAML(&d2).Code(&httpCode).Do()
if err != nil || httpCode != 200{
    fmt.Printf("send fail:%s\n", err)
}

```

### xml
* SetXML() 设置请求http body为xml
* BindXML() 解析响应http body里面的xml到结构体里面

发送xml到服务端，然后把服务端返回的xml结果解析到结构体里面
```go
type data struct {
    Id int `xml:"id"`
    Data string `xml:"data"`
}


var d1, d2 data
var httpCode int 


err := gout.POST(":8080/test.xml").SetXML(&d1).BindXML(&d2).Code(&httpCode).Do()
if err != nil || httpCode != 200{
    fmt.Printf("send fail:%s\n", err)
}

```

### form-data
* SetForm() 设置http body 为multipart/form-data格式数据

客户端发送multipart/form-data到服务端,curl用法等同go代码
```bash
curl -F mode=A -F text="good" -F voice=@./test.pcm -f voice2=@./test2.pcm url
```

* 使用gout.H
```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

func main() {

	code := 0
	err := gout.
		POST(":8080/test").
		// 打开debug模式
		Debug(true).
		SetForm(
			gout.H{
				"mode": "A",
				"text": "good",
				// 从文件里面打开
				"voice":  gout.FormFile("test.pcm"),
				"voice2": gout.FormMem("pcm"),
			},
		).
		//解析http code，如不关心可以不设置
		Code(&code).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	if code != 200 {
	}
}

/*
> POST /test HTTP/1.1
> Content-Type: multipart/form-data; boundary=2b0685e5b98e540f80b247d5e7c1283807aa07e62b827543859a6db765a8
>

--2b0685e5b98e540f80b247d5e7c1283807aa07e62b827543859a6db765a8
Content-Disposition: form-data; name="mode"

A
--2b0685e5b98e540f80b247d5e7c1283807aa07e62b827543859a6db765a8
Content-Disposition: form-data; name="text"

good
--2b0685e5b98e540f80b247d5e7c1283807aa07e62b827543859a6db765a8
Content-Disposition: form-data; name="voice"; filename="voice"
Content-Type: application/octet-stream

pcm pcm pcm

--2b0685e5b98e540f80b247d5e7c1283807aa07e62b827543859a6db765a8
Content-Disposition: form-data; name="voice2"; filename="voice2"
Content-Type: application/octet-stream

pcm
--2b0685e5b98e540f80b247d5e7c1283807aa07e62b827543859a6db765a8--


< HTTP/1.1 200 OK
< Server: gurl-server
< Content-Length: 0
*/
 
```

* 使用结构体
```go
type testForm struct {
    Mode string `form:"mode"`
    Text string `form:"text"`
    Voice string `form:"voice" form-file:"true"` //从文件中读取 
    Voice2 []byte `form:"voice2" form-mem:"true"`  //从内存中构造
}

type rsp struct{
    ErrMsg string `json:"errmsg"`
    ErrCode int `json:"errcode"`
}

t := testForm{}
r := rsp{}
code := 0

err := gout.POST(url).SetForm(&t).ShoudBindJSON(&r).Code(&code).Do()
if err != nil {

}
```
### x-www-form-urlencoded
* 使用SetWWWForm函数实现发送x-www-form-urlencoded类型数据
```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

func main() {

	code := 0
	err := gout.
		POST(":8080/post").
		// 打开debug模式
		Debug(true).
		// 设置x-www-form-urlencoded数据
		SetWWWForm(
			gout.H{
				"int":     3,
				"float64": 3.14,
				"string":  "test-www-Form",
			},
		).
		// 关心http code 返回值设置
		Code(&code).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	if code != 200 {
	}
}

/*
> POST /post HTTP/1.1
> Content-Type: application/x-www-form-urlencoded
>

float64=3.14&int=3&string=test-www-Form

< HTTP/1.1 200 OK
< Content-Length: 0
< Server: gurl-server

*/

```

### callback
callback主要用在，服务端会返回多种格式body的场景, 比如404返回的是html, 200返回json。
这时候要用Callback挂载多种处理函数
```go

func main() {
	
	r, str404 := Result{}, ""
	code := 0

	err := gout.GET(":8080").Code(&code).Callback(func(c *gout.Context) (err error) {

		switch c.Code {
		case 200:
			err = c.BindJSON(&r)
		case 404:
			err = c.BindBody(&str404)
		}
		return

	}).Do()

	if err != nil {
		fmt.Printf("err = %s\n", err)
		return
	}

	fmt.Printf("http code = %d, str404(%s), result(%v)\n", code, str404, r)
}

```
## benchmark
### number
下面的例子，起了20并发。对:8080端口的服务，发送3000次请求进行压测，内容为json结构
```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

const (
	benchNumber     = 30000
	benchConcurrent = 20
)

func main() {
	err := gout.
		POST(":8080").                     //压测本地8080端口
		SetJSON(gout.H{"hello": "world"}). //设置请求body内容
		Filter().                          //打开过滤器
		Bench().                           //选择bench功能
		Concurrent(benchConcurrent).       //并发数
		Number(benchNumber).               //压测次数
		Do()

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

```
### duration
下面的例子，起了20并发。对:8080端口的服务，压测持续时间为10s，内容为json结构
```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"time"
)

const (
	benchTime       = 10 * time.Second
	benchConcurrent = 30
)

func main() {
	err := gout.
		POST(":8080").                     //压测本机8080端口
		SetJSON(gout.H{"hello": "world"}). //设置请求body内容
		Filter().                          //打开过滤器
		Bench().                           //选择bench功能
		Concurrent(benchConcurrent).       //并发数
		Durations(benchTime).              //压测时间
		Do()

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

```
### rate
下面的例子，起了20并发。对:8080端口的服务，压测总次数为3000次，其中每秒发送1000次。内容为json结构
```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

const (
	benchNumber     = 3000
	benchConcurrent = 20
)

func main() {
	err := gout.
		POST(":8080").                     //压测本机8080端口
		SetJSON(gout.H{"hello": "world"}). //设置请求body内容
		Filter().                          //打开过滤器
		Bench().                           //选择bench功能
		Rate(1000).                        //每秒发1000请求
		Concurrent(benchConcurrent).       //并发数
		Number(benchNumber).               //压测次数
		Do()

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

```
## timeout
setimeout是request级别的超时方案。相比http.Client级别，更灵活。
```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"time"
)

func main() {
	err := gout.GET(":8080").
		SetTimeout(2 * time.Second).
		Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

```
## proxy
* SetProxy 设置代理服务地址
```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"log"
)

func main() {
	c := &http.Client{}
	s := ""
	err := gout.
		New(c).
		GET("www.qq.com").
		// 设置proxy服务地址
		SetProxy("http://127.0.0.1:7000").
		// 绑定返回数据到s里面
		BindBody(&s).
		Do()

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(s)
}

```
## cookie
* SetCookies设置cookie, 可以设置一个或者多个cookie

```go
package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"net/http"
)

func main() {

	// === 发送多个cookie ====
	err := gout.
		// :8080/cookie是http://127.0.0.1:8080/cookie的简写
		GET(":8080/cookie").
		//设置debug模式
		Debug(true).
		SetCookies(
			//设置cookie1
			&http.Cookie{
				Name:  "test1",
				Value: "test1",
			},
			//设置cookie2
			&http.Cookie{
				Name:  "test2",
				Value: "test2",
			},
		).
		Do()

	if err != nil {
		fmt.Println(err)
		return
	}

	// === 发送一个cookie ===
	err = gout.
		// :8080/cookie/one是http://127.0.0.1:8080/cookie/one的简写
		GET(":8080/cookie/one").
		//设置debug模式
		Debug(true).
		SetCookies(
			//设置cookie1
			&http.Cookie{
				Name:  "test3",
				Value: "test3",
			},
		).
		Do()
	fmt.Println(err)

}

```

## context
* WithContext设置context，可以取消http请求
### cancel
```go
package main

import (
    "context"
    "github.com/guonaihong/gout"
    "time"
)

func main() {
    //　声明一个context
    ctx, cancel := context.WithCancel(context.Background())

    //调用cancel可取消http请求
    go func() {
        time.Sleep(time.Second)
        cancel()
    }() 

    err := gout.
        GET("127.0.0.1:8080/cancel"). //设置GET请求以及需要访问的url
        WithContext(ctx).             //设置context, 外层调用cancel函数就可取消这个http请求
        Do()

    if err != nil {
    }   
}

```

## unix socket
* UnixSocket可以把http底层通信链路由tcp修改为unix domain socket  
下面的例子，会通过domain socket发送http GET请求，http body的内容是hello world
```go
package main

import (
    "fmt"
    "github.com/guonaihong/gout"
    "net/http"
)

func main() {
    c := http.Client{}

    g := gout.
        New(&c).
        UnixSocket("/tmp/test.socket") //设置unixsocket文件位置

    err := g.
        GET("http://a/test").   //设置GET请求
        SetBody("hello world"). //设置body内容
        Do()
    fmt.Println(err)
}
```
## http2 doc
go 使用https访问http2的服务会自动启用http2协议，这里不需要任何特殊处理
* https://http2.golang.org/ (bradfitz建的http2测试网址,里面大约有十来个测试地址，下面的例子选了一个) 
```go
package main

import (
    "fmt"
    "github.com/guonaihong/gout"
)

func main() {
    s := ""
    err := gout.
        GET("https://http2.golang.org/reqinfo"). //设置GET请求和请求url
        Debug(true).                             //打开debug模式，可以看到请求数据和响应数据
        SetBody("hello, ###########").           //设置请求body的内容，如果你的请求内容是json格式，需要使用SetJSON函数
        BindBody(&s).                            //解析响应body内容
        Do()                                     //结束函数

    if err != nil {
        fmt.Printf("send fail:%s\n", err)
    }   
    _ = s 
}

```
## debug mode
### color
该模式主要方便调试用的，默认开启颜色高亮
* Debug(true)

```go
func main() {
	
	err := gout.POST(":8080/colorjson").
		Debug(true).
		SetJSON(gout.H{"str": "foo",
			"num":   100,
			"bool":  false,
			"null":  nil,
			"array": gout.A{"foo", "bar", "baz"},
			"obj":   gout.H{"a": 1, "b": 2},
		}).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}
```
### customize
debug 自定义模式，可传递函数。下面的只有传递IOS_DEBUG环境变量才输出日志
```go
package main

import (
    "fmt"
    "github.com/guonaihong/gout"
    "os"
)

func IOSDebug() gout.DebugOpt {
    return gout.DebugFunc(func(o *gout.DebugOption) {
        if len(os.Getenv("IOS_DEBUG")) > 0 { 
            o.Debug = true
        }
    })  
}

func main() {

    s := ""
    err := gout.
        GET("127.0.0.1:8080").
        // Debug可以支持自定义方法
        // 可以实现设置某个环境变量才输出debug信息
        // 或者debug信息保存到文件里面，可以看下_example/15-debug-save-file.go
        Debug(IOSDebug()).
        SetBody("test hello").
        BindBody(&s).
        Do()

    fmt.Printf("err = %v\n", err)
}

// env IOS_DEBUG=true go run customize.go
```
### no-color
no-color 使用gout.NoColor()关闭颜色高亮
```go
func main() {
	
	err := gout.POST(":8080/colorjson").
		Debug(gout.NoColor()).
		SetJSON(gout.H{"str": "foo",
			"num":   100,
			"bool":  false,
			"null":  nil,
			"array": gout.A{"foo", "bar", "baz"},
			"obj":   gout.H{"a": 1, "b": 2},
		}).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

```
# 特色功能示例
## forward gin data
gout 设计之初就考虑到要和gin协同工作的可能性，下面展示如何方便地使用gout转发gin绑定的数据。
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
)

type testQuery struct {
	Size int    `query:"size" form:"size"` // query tag是gout设置查询字符串需要的
	Page int    `query:"page" form:"page"`
	Ak   string `query:"ak" form:"ak"`
}

//下一个服务节点
func nextSever() {
	r := gin.Default()

	r.GET("/query", func(c *gin.Context) {
		q := testQuery{}
		err := c.ShouldBindQuery(&q)
		if err != nil {
			return
		}
		c.JSON(200, q)
	})
	r.Run(":1234")
}

func main() {
	go nextSever()
	r := gin.Default()

	// 演示把gin绑定到的查询字符串转发到nextServer节点
	r.GET("/query", func(c *gin.Context) {
		q := testQuery{}
		// 绑定查询字符串
		err := c.ShouldBindQuery(&q)
		if err != nil {
			return
		}

		// 开发转发, 复用gin所用结构体变量q
		code := 0 // http code
		err := gout.
			//发起GET请求
			GET("127.0.0.1:1234/query").
			//设置查询字符串
			SetQuery(q).
			//关心http server返回的状态码 设置该函数
			Code(&code).
			Do()
		if err != nil || code != 200 { /* todo Need to handle errors here */
		}
		c.JSON(200, q)
	})

	r.Run()
}

// http client
// curl '127.0.0.1:8080/query?size=10&page=20&ak=test'
```
# FAQ

## gout benchmark性能如何
下面是与apache ab的性能对比 [_example/16d-benchmark-vs-ab.go](_example/16d-benchmark-vs-ab.go)

![gout-vs-ab.png](https://github.com/guonaihong/images/blob/master/gout/gout-vs-ab.png?raw=true)