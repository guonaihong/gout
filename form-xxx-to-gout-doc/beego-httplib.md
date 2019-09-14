
### 从beego/httplib 迁移到gout

### 发送大片的数据
#### httplib
```go
package main

import (
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"log"
)

func main() {
	req := httplib.Post("http://127.0.0.1:1234")
	bt, err := ioutil.ReadFile("hello.txt")
	if err != nil {
		log.Fatal("read file err:", err)
	}
	req.Body(bt)
}

```
#### gout
```go
package main

import (
        "github.com/guonaihong/gout"
        "io/ioutil"
        "log"
)

func main() {
        bt, err := ioutil.ReadFile("hello.txt")
        if err != nil {
                log.Fatal("read file err:", err)
        }
        gout.POST("http://127.0.0.1:1234").SetBody(bt).Do()
}

```
### 设置 header 信息
#### httplib
```go
req := httplib.Post("http://beego.me/")
req.Header("Accept-Encoding","gzip,deflate,sdch")
req.Header("Host","beego.me")
req.Header("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36
```

#### gout
```go
gout.POST("http://beego.me/").SetHeader(
	gout.H{"Accept-Encoding":"gzip,deflate,sdch",
	              "Host":"beego.me",
	               "User-Agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36"}).Do()
```
### 发送大片的数据

#### httplib
```go
package main

import (
        "github.com/astaxie/beego/httplib"
        "log"
)

func main() {
        b := httplib.Post("http://127.0.1:1234/")

        b.Param("username", "test")
        b.Param("password", "123456")
        b.PostFile("uploadfile1", "t.go")
        b.PostFile("uploadfile2", "t2.go")

        str, err := b.String()

        if err != nil {
                log.Fatal(err)
        }
        log.Printf("str = %s\n", str)

}

```

##### gout
```go
package main

import (
        "github.com/guonaihong/gout"
        "log"
)

func main() {
        var str string
        err := gout.POST("http://127.0.0.1:1234").SetForm(
                gout.H{"username": "test",
                        "password":    "123456",
                        "uploadfile1": gout.FormFile("./t.go"),
                        "uploadfile2": gout.FormFile("./t2.go"),
                }).BindBody(&str).Do()

        if err != nil {
                log.Fatal(err)
        }
        log.Printf("str = %s\n", str)
}


```
### 获取返回结果
#### httplib
  
> 获取返回结果
上面这些都是在发送请求之前的设置，接下来我们开始发送请求，然后如何来获取数据呢？主要有如下几种方式：

  - 返回 Response 对象，req.Response() 方法

> 这个是 http.Response 对象，用户可以自己读取 body 的数据等。

  *    返回 bytes,req.Bytes() 方法

> 直接返回请求 URL 返回的内容

  - 返回 string，req.String() 方法

> 直接返回请求 URL 返回的内容

  - 保存为文件，req.ToFile(filename) 方法

> 返回结果保存到文件名为 filename 的文件中

  - 解析为 JSON 结构，req.ToJSON(&result) 方法

> 返回结构直接解析为 JSON 格式，解析到 result 对象中

  - 解析为 XML 结构，req.ToXml(&result) 方法

> 返回结构直接解析为 XML 格式，解析到 result 对象中

#### gout
> 获取返回结果
上面这些都是在发送请求之前的设置，接下来我们开始发送请求，然后如何来获取数据呢？主要有如下几种方式：

> gout使用链式调用，返回结果都是通过BindXXX函数获得

  * BindBody
BindBody支持，string, []byte两种类型，不需要如httplib设计出两个接口
> 直接返回请求 URL 返回的内容
   * BindJSON
> 返回结构直接解析为 JSON 格式，解析到 result 对象中

 * BindXML
  解析为 XML 结构，BindXML(&result) 方法

> 返回结构直接解析为 XML 格式，解析到 result 对象中
