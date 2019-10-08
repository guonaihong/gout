# gout
gout 是go写的http 客户端，为提高工作效率而开发

[![Build Status](https://travis-ci.org/guonaihong/gout.png)](https://travis-ci.org/guonaihong/gout)
[![codecov](https://codecov.io/gh/guonaihong/gout/branch/master/graph/badge.svg)](https://codecov.io/gh/guonaihong/gout)

## 内容
- [安装](#安装)
- [技能树](#技能树)
- [迁移](#迁移)
- [API示例](#api-示例)
    - [GET POST PUT DELETE PATH HEAD OPTIONS](#get-post-put-delete-path-head-options)
    - [group](#group)
    - [query](#query)
    - [http header](#http-header)
    - [http body](#http-body)
        - [body](#body)
        - [json](#json)
        - [yaml](#yaml)
        - [xml](#xml)
        - [form](#form)
        - [callback](#callback)
    - [proxy](#proxy)
    - [cookie](#cookie)
    - [context](#context)
        - [timeout](#timeout)
        - [cancel](#cancel)
    - [unix socket](#unix-socket)
    

## 安装
```
env GOPATH=`pwd` go get github.com/guonaihong/gout
```
## 技能树
![gout.png](https://github.com/guonaihong/images/blob/master/gout/gout.png)

## 迁移
* [httplib](./form-xxx-to-gout-doc/beego-httplib.md)
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
* SetQuery() 设置http 查询字符串

```go
g := gout.New(nil)
code := 0

if err := g.GET(":8080/testquery").SetQuery(/*看下面支持的情况*/).Code(&code).Do(); err != nil {
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
* SetHeader() 设置http header
* BindHeader() 解析响应http header

对gout来说，既支持客户端发送http header,也支持解码服务端返回的http header
```go
type testHeader struct {
    CheckIn string `header:checkin`
    CheckOut string `header:checkout`
}

t := testheader{}

g := gout.New(nil)
code := 0

if err := g.GET(":8080/testquery").Code(&code).SetHeader(/*看下面支持的类型*/).BindHeader(&t).Do(); err != nil {
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
* SetBody 设置string, []byte等类型数据到http body里面
* BindBody bind body到string, []byte等类型变量里面
```go

err := gout.New(nil).POST(url).SetBody(/*支持的类型如下*/).Do()

var s string
err := gout.New(nil).GET(url).BindBody(&s/*支持的指针类型变量如下*/).Do()

```
## 支持的类型有
* int, int8, int16, int32, int64
* uint, uint8, uint16, uint32, uint64
* string
* []byte
* float32, float64

## 明确不支持的类型有
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

g := gout.New(nil)

err := g.POST(":8080/test.json").SetJSON(&d1).BindJSON(&d2).Code(&httpCode).Do()
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

g := gout.New(nil)

err := g.POST(":8080/test.yaml").SetYAML(&d1).BindYAML(&d2).Code(&httpCode).Do()
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

g := gout.New(nil)

err := g.POST(":8080/test.xml").SetXML(&d1).BindXML(&d2).Code(&httpCode).Do()
if err != nil || httpCode != 200{
    fmt.Printf("send fail:%s\n", err)
}

```

### form
* SetForm() 设置http body 为multipart/form-data格式数据

客户端发送multipart/form-data到服务端,curl用法等同go代码
```bash
curl -F mode=A -F text="good" -F voice=@./test.pcm -f voice2=@./test2.pcm url
```

* 使用gout.H
```go
code := 0
err := gout.New(nil).
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

err := gout.New(nil).POST(url).SetForm(&t).ShoudBindJSON(&r).Code(&code).Do()
if err != nil {

}
```

### callback
callback主要用在，服务端会返回多种格式body的场景
```go
g := gout.New(nil)

r , errCode := Result{}, 0

g.GET(url).Callback(func(c *gout.Context) error {

    switch c.Code {
        case 200:
            c.BindJSON(&r)
        case 500:
            c.BindBody(&errCode)
    }

    return nil
})
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
	var s string
	err := gout.Def().GET("www.qq.com").SetProxy("http://127.0.0.1:7000").BindBody(&s).Do()
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

        // 发送两个cookie
        err := gout.GET(":1234/cookie").SetCookies(&http.Cookie{Name: "test1", Value: "test1"},
                &http.Cookie{Name: "test2", Value: "test2"}).Do()

        // 发送一个cookie
        err = gout.GET(":1234/cookie/one").SetCookies(&http.Cookie{Name: "test3", Value: "test3"}).Do()
        fmt.Println(err)
}
```

## context
* WithContext设置context，可以取消http请求
### timeout
```go
func main() {
	// 给http请求 设置超时
	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)

	err := gout.GET("127.0.0.1:8080/timeout").WithContext(ctx).Do()

}
```
### cancel
```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second)
		cancel() //调用cancel可取消http请求
	}()

	err := gout.GET("127.0.0.1:8080/cancel").WithContext(ctx).Do()
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
	g := gout.New(&c).UnixSocket("/tmp/test.socket")
	err := g.GET("http://a/test").SetBody("hello world").Do()
	fmt.Println(err)
}
```
