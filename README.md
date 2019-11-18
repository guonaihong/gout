# gout
gout 是go写的http 客户端，为提高工作效率而开发

[![Build Status](https://travis-ci.org/guonaihong/gout.png)](https://travis-ci.org/guonaihong/gout)
[![codecov](https://codecov.io/gh/guonaihong/gout/branch/master/graph/badge.svg)](https://codecov.io/gh/guonaihong/gout)

## 演示
![gout-example.gif](https://github.com/guonaihong/images/blob/master/gout/gout-example.gif?raw=true)

## 内容
- [安装](#安装)
- [技能树](#技能树)
- [迁移文档](#迁移文档)
- [example](#example)
- [API示例](#api-示例)
    - [GET POST PUT DELETE PATH HEAD OPTIONS](#get-post-put-delete-path-head-options)
    - [group](#group)
    - [query](#query)
    - [http header](#http-header)
		- [req header](#req-header)
		- [rsp header](#rsp-header)
    - [http body](#http-body)
        - [body](#body)
            - [SetBody](#setbody)
            - [BindBody](#bindbody)
        - [json](#json)
        - [yaml](#yaml)
        - [xml](#xml)
        - [form-data](#form-data)
        - [x-www-form-urlencoded](#x-www-form-urlencoded)
        - [callback](#callback)
    - [proxy](#proxy)
    - [cookie](#cookie)
    - [context](#context)
        - [timeout](#timeout)
        - [cancel](#cancel)
    - [unix socket](#unix-socket)
    - [http2 doc](#http2-doc)
    - [debug mode](#debug-mode)
        - [color](#color)
        - [customize](#customize)
        - [no-color](#no-color)
 - [特色功能示例](#特色功能示例)
    - [forward gin data](#forward-gin-data)

## 安装
```
env GOPATH=`pwd` go get github.com/guonaihong/gout
```
## 技能树
<details>

![gout.png](https://github.com/guonaihong/images/blob/master/gout/gout.png)

</details>

## 迁移文档
主要方便下面的用户迁移到gout
* [httplib](./to-gout-doc/beego-httplib.md)
* [resty](./to-gout-doc/resty-doc.md)
# example
 [examples](./_example) 目录下面的例子，都是可以直接跑的。如果觉得运行例子还是不明白用法，可以把你迷惑的地方写出来，让后提[issue](https://github.com/guonaihong/gout/issues/new)
# API 示例
## GET POST PUT DELETE PATH HEAD OPTIONS
```go
// 创建一个实例
// 也可以直接调用包里面的GET, POST方法
// 比如gout.GET(url)

g := gout.New(nil)

// 发送GET方法
g.GET(url).Do()

// 发送POST方法
g.POST(url).Do()

// 发送PUT方法
g.PUT(url).Do()

// 发送DELETE方法
g.DELETE(url).Do()

// 发送PATH方法
g.PATCH(url).Do()

// 发送HEAD方法
g.HEAD(url).Do()

// 发送OPTIONS
g.OPTIONS(url).Do()
```
## group
路由组
```go
g := New(nil)

v1 := g.Group(ts.URL + "/v1")
err := v1.POST("/login").Next(). // http://127.0.0.1:80/v1/login
    POST("/submit").Next().      // http://127.0.0.1:80/v1/submit
    POST("/read").Do()           // http://127.0.0.1:80/v1/read

if err != nil {
}

v2 := g.Group(ts.URL + "/v2")
err = v2.POST("/login").Next(). // http://127.0.0.1:80/v2/login
    POST("/submit").Next().     // http://127.0.0.1:80/v2/submit
    POST("/read").Do()          // http://127.0.0.1:80/v2/read

if err != nil {
}
```
## query

### easy example
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
code := 0

err := gout.
	GET(":8080/testquery").
	SetQuery( /*看下面支持的情况*/ ).
	Code(&code). //解析http code，如不关心服务端返回状态吗，不设置该函数即可
	Do()
if err != nil {

}

/*
SetQuery支持的类型有
* string
* map[string]interface{}，可以使用gout.H别名
* struct
* array, slice(长度必须是偶数)
*/

// string
SetQuery("check_in=2019-06-18&check_out=2018-06-18")

// gout.H 或者 map[string]interface{}
SetQuery(gout.H{
    "check_in":"2019-06-18",
    "check_out":"2019-06-18",
})

// struct
type testQuery struct {
    CheckIn string `query:checkin`
    CheckOut string `query:checkout`
}

SetQuery(&testQuery{CheckIn:2019-06-18, CheckOut:2019-06-18})

// array or slice
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
< Date: Mon, 18 Nov 2019 12:54:46 GMT
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
< Date: Mon, 18 Nov 2019 12:59:37 GMT
*/

```
### SetHeader和BindHeader支持的更多类型
```go
type testHeader struct {
    CheckIn string `header:checkin`
    CheckOut string `header:checkout`
}

t := testheader{}

code := 0

if err := gout.GET(":8080/testquery").Code(&code).SetHeader(/*看下面支持的类型*/).BindHeader(&t).Do(); err != nil {
}
```
* BindHeader支持的类型有
```go
// struct
type testHeader struct {
    CheckIn string `header:checkin`
    CheckOut string `header:checkout`
}
```
 结构体
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
* SetBody 设置string, []byte等类型数据到http body里面
```go
// 设置string变量至请求的http body
err := gout.POST(url).SetBody("hello world"/*更多支持类型请看下面*/).Do()

// 设置实现io.Reader接口的变量至 请求的http body
err = gout.POST(url).SetBody(bytes.NewBufferString("hello world")).Code(&code).Do()
```
#### bindBody
* BindBody bind body到string, []byte等类型变量里面
```go
// 解析http body到string类型变量里面
var s string
err := gout.GET(url).BindBody(&s/*更多支持指针类型变量请看下面*/).Do()

// 解析http body至实现io.Writer接口的变量里面
var b bytes.Buffer{}
err = gout.GET(url).BindBody(&b).Code(&code).Do()
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

* SetJSON()  设置请求http body为json
* BindJSON()  解析响应http body里面的json到结构体里面

发送json到服务端，然后把服务端返回的json结果解析到结构体里面
```go
type data struct {
    Id int `json:"id"`
    Data string `json:"data"`
}


var d1, d2 data
var httpCode int

err := gout.POST(":8080/test.json").SetJSON(&d1).BindJSON(&d2).Code(&httpCode).Do()
if err != nil || httpCode != 200{
    fmt.Printf("send fail:%s\n", err)
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
code := 0
err := gout.
    POST(":8080/test").
    SetForm(gout.H{"mode": "A",
        "text":   "good",
        "voice":  gout.FormFile("test.pcm"),
        "voice2": gout.FormMem("pcm")}).Code(&code).Do()

if err != nil {
    fmt.Printf("%s\n", err)
}   

if code != 200 {
}   
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
func main() {
    err := gout.POST(":8080/post").
		Debug(true).
		SetWWWForm(gout.H{
			"int":     3,
			"float64": 3.14,
			"string":  "test-www-Form",
		}).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}
/*
output:
> POST /post HTTP/1.1
> Content-Type: application/x-www-form-urlencoded
>

float64=3.14&int=3&string=test-www-Form

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
	err := gout.New(c).GET("www.qq.com").SetProxy("http://127.0.0.1:7000").BindBody(&s).Do()
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
                Value: "test1"},
            //设置cookie2
            &http.Cookie{
                Name:  "test2",
                Value: "test2"}).
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
                Value: "test3"}).
        Do()
    fmt.Println(err)

}

```

## context
* WithContext设置context，可以取消http请求
### timeout
```go
package main

import (
    "context"
    "github.com/guonaihong/gout"
    "time"
)

func main() {
    // 给http请求 设置超时
    ctx, _ := context.WithTimeout(context.Background(), time.Second*1)

    err := gout.
        GET("127.0.0.1:8080/timeout"). //设置GET请求以及url地址
        WithContext(ctx).              //设置context,　如果传过来的ctx超时，这个http请求就会被取消
        Do()                           //结束函数

    if err != nil {
    }
}

```
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
    err := gout.GET("127.0.0.1:1234").Debug(IOSDebug()).SetBody("test hello").BindBody(&s).Do()
    fmt.Printf("err = %v\n", err)
}

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
	Size int    `query:"size" form:"size"`　// query tag是gout设置查询字符串需要的
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

	// 当前服务
	r.GET("/query", func(c *gin.Context) {
		q := testQuery{}
		err := c.ShouldBindQuery(&q)
		if err != nil {
			return
		}

		// Send to the next service

		code := 0 // http code
		err := gout.
			GET("127.0.0.1:1234/query"). //发起GET请求
			SetQuery(q).                 //设置查询字符串
			Code(&code).                 //解析http code，如果不关系服务端的code　可以不设置该函数
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
