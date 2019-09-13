### 从beego/httplib 迁移到gout


### 设置 header 信息
* httplib
```go
req := httplib.Post("http://beego.me/")
req.Header("Accept-Encoding","gzip,deflate,sdch")
req.Header("Host","beego.me")
req.Header("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36
```

* gout
```go
gout.POST("http://beego.me/").SetHeader(
	gout.H{"Accept-Encoding":"gzip,deflate,sdch",
	"Host":"beego.me",
	"User-Agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36"}).Do()
```

### 获取返回结果
* httplib
获取返回结果
上面这些都是在发送请求之前的设置，接下来我们开始发送请求，然后如何来获取数据呢？主要有如下几种方式：

** 返回 Response 对象，req.Response() 方法

> 这个是 http.Response 对象，用户可以自己读取 body 的数据等。

** 返回 bytes,req.Bytes() 方法

> 直接返回请求 URL 返回的内容

** 返回 string，req.String() 方法

> 直接返回请求 URL 返回的内容

** 保存为文件，req.ToFile(filename) 方法

> 返回结果保存到文件名为 filename 的文件中

** 解析为 JSON 结构，req.ToJSON(&result) 方法

> 返回结构直接解析为 JSON 格式，解析到 result 对象中

** 解析为 XML 结构，req.ToXml(&result) 方法

> 返回结构直接解析为 XML 格式，解析到 result 对象中

*gout
获取返回结果
上面这些都是在发送请求之前的设置，接下来我们开始发送请求，然后如何来获取数据呢？主要有如下几种方式：

gout使用链式调用，返回结果都是通过BindXXX函数获得

** BindBody
BindBody支持，string, []byte两种类型，不需要如httplib设计出两个接口
> 直接返回请求 URL 返回的内容

> 直接返回请求 URL 返回的内容
** 解析为 JSON 结构，BindJSON(&result) 方法

> 返回结构直接解析为 JSON 格式，解析到 result 对象中

** 解析为 XML 结构，BindXML(&result) 方法

> 返回结构直接解析为 XML 格式，解析到 result 对象中
