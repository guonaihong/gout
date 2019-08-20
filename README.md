# gout
gout 是go写的http 客户端，为提高工作效率而开发

## 内容
- [安装](#安装)
- [API示例](#api-示例)
    - [GET POST PUT DELETE PATH HEAD OPTIONS](#get-post-put-delete-path-head-options)
    - [json](#json)
    - [yaml](#yaml)
    - [xml](#xml)
    - [form](#form)
    - [query](#query)
    - [http header](#http-header)
    - [group](#group)

## 安装
```
env GOPATH=`pwd` go get github.com/guonaihong/gout
```
## GET POST PUT DELETE PATH HEAD OPTIONS
```go
// 创建一个实例
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
g.PATH(url).Do()

// 发送HEAD方法
g.HEAD(url).Do()

// 发送OPTIONS
g.OPTIONS(url).Do()
```

## json

* ToJSON()  设置请求http body为json
* ShouldBindJSON()  解析响应http body里面的json到结构体里面

发送json到服务端，然后把服务端返回的json结果解析到结构体里面
```go
type data struct {
    Id int `json:"id"`
    Data string `json:"data"`
}


var d1, d2 data
var httpCode int

g := gout.New(nil)

err := g.POST(":8080/test.json").ToJSON(&d1).ShouldBindJSON(&d2).Code(&httpCode).Do()
if err != nil || httpCode != 200{
    fmt.Printf("send fail:%s\n", err)
}
```

## yaml
* ToYAML() 设置请求http body为yaml
* ShouldBindYAML() 解析响应http body里面的yaml到结构体里面

发送yaml到服务端，然后把服务端返回的yaml结果解析到结构体里面
```go
type data struct {
    Id int `yaml:"id"`
    Data string `yaml:"data"`
}


var d1, d2 data
var httpCode int 

g := gout.New(nil)

err := g.POST(":8080/test.yaml").ToYAML(&d1).ShouldBindYAML(&d2).Code(&httpCode).Do()
if err != nil || httpCode != 200{
    fmt.Printf("send fail:%s\n", err)
}

```

## xml
* ToXML() 设置请求http body为xml
* ShouldBindXML() 解析响应http body里面的xml到结构体里面

发送xml到服务端，然后把服务端返回的xml结果解析到结构体里面
```go
type data struct {
    Id int `xml:"id"`
    Data string `xml:"data"`
}


var d1, d2 data
var httpCode int 

g := gout.New(nil)

err := g.POST(":8080/test.xml").ToXML(&d1).ShouldBindXML(&d2).Code(&httpCode).Do()
if err != nil || httpCode != 200{
    fmt.Printf("send fail:%s\n", err)
}

```

## form
* ToForm() 设置http body 为multipart/form-data格式数据

客户端发送multipart/form-data到服务端,curl用法等同go代码
```bash
curl -F mode=A -F text="good" -F voice=@./test.pcm -f voice2=@./test2.pcm url
```

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

err := gout.New(nil).ToForm(&t).ShoudBindJSON(&r).Code(&code).Do()
if err != nil {

}
```


## query
* ToQuery() 设置http 查询字符串

```go
out := gout.New(nil)
code := 0

if err := out.GET(":8080/testquery").ToQuery(/*看下面支持的情况*/).Code(&code).Do(); err != nil {
}

/*
ToQuery支持的类型有
* string
* map[string]interface{}，可以使用gout.H别名
* struct
* array, slice(长度必须是偶数)
*/

// string
ToQuery("check_in=2019-06-18&check_out=2018-06-18")

// gout.H 或者 map[string]interface{}
ToQuery(gout.H{
    "check_in":"2019-06-18",
    "check_out":"2019-06-18",
})

// struct
type testQuery struct {
    CheckIn string `query:checkin`
    CheckOut string `query:checkout`
}

ToQuery(&testQuery{CheckIn:2019-06-18, CheckOut:2019-06-18})

// array or slice
// ?active=enable&action=drop
ToQuery([]string{"active", "enable", "action", "drop"})`
```

## http header
* ToHeader() 设置http header
* ShouldBindHeader() 解析响应http header

对gout来说，既支持客户端发送http header,也支持解码服务端返回的http header
```go
//ShouldBindHeader支持的类型有
// 结构体
type testHeader struct {
    CheckIn string `header:checkin`
    CheckOut string `header:checkout`
}

t := testheader{}

out := gout.New(nil)
code := 0

if err := out.GET(":8080/testquery").Code(&code).ToHeader(/*看下面支持的类型*/).ShoudBindHeader(&t).Do(); err != nil {
}
```
* ToHeader支持的类型有
```go
map[string]interface{}，可以使用gout.H别名
struct
array, slice(长度必须是偶数)
// gout.H 或者 map[string]interface{}
ToHeader(gout.H{
    "check_in":"2019-06-18",
    "check_out":"2019-06-18",
})

// struct
type testHeader struct {
    CheckIn string `header:checkin`
    CheckOut string `header:checkout`
}

ToHeader(&testHeader{CheckIn:2019-06-18, CheckOut:2019-06-18})

// array or slice
// -H active:enable -H action:drop
ToHeader([]string{"active", "enable", "action", "drop"})
```

## group
路由组
```go
out := New(nil)

// http://127.0.0.1:80/v1/login
// http://127.0.0.1:80/v1/submit
// http://127.0.0.1:80/v1/read
v1 := out.Group(ts.URL + "/v1")
err := v1.POST("/login").Next().
    POST("/submit").Next().
    POST("/read").Do()

if err != nil {
}

// http://127.0.0.1:80/v2/login
// http://127.0.0.1:80/v2/submit
// http://127.0.0.1:80/v2/read
v2 := out.Group(ts.URL + "/v2")
err = v2.POST("/login").Next().
    POST("/submit").Next().
    POST("/read").Do()

if err != nil {
}
```
