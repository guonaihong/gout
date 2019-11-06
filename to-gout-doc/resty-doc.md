### 从resty 迁移到gout
### 目录
- [API示例](#api-示例)
    - [query](#query)
        - [resty-query](#resty-query)
        - [gout-query](#gout-query)
    - [header](#header)
        - [gout-header](#gout-header)
        - [gout-header](#gout-header)
    - [http body](#http-body)
        - [json](#json)
            - [resty-json](#resty-json)
            - [gout-json](#gout-json)
        - [body](#body)
            - [resty-body](#resty-body)
            - [gout-body](#gout-body)
    - [test server](#test-server)   
### query
#### resty-query
```go
func restyQuery() {
    client := resty.New()
    // string
    _, err := client.R().SetQueryString("k1=v1&k2=v2").Post("http://127.0.0.1:8080/json")
    fmt.Printf("err = %v\n", err)

    // slice or array
    _, err = client.R().SetQueryParams(map[string]string{"k1": "v1", "k2": "v2"}).Post("http://127.0.0.1:8080/json")
    fmt.Printf("err = %v\n", err)

    // map
    _, err = client.R().SetQueryParams(map[string]string{"k1": "v1", "k2": "v2"}).Post("http://127.0.0.1:8080/json")
    fmt.Printf("err = %v\n", err)

    // struct 不支持
}

```


#### gout-query
```go
func goutQuery() {
    // string
    err := gout.GET(":8080/query").SetQuery("k1=v1&k2=v2").Do()
    fmt.Printf("err = %v\n", err)

    // slice or array
    err = gout.GET(":8080/query").SetQuery([]string{"k1", "v1", "k2", "v2"}).Do()
    fmt.Printf("err = %v\n", err)

    // map
    err = gout.GET(":8080/query").SetQuery(gout.H{"k1": "v1", "k2": "v2"}).Do()
    fmt.Printf("err = %v\n", err)

    // struct
    err = gout.GET(":8080/query").SetQuery(testData{K1: "v1", K2: "v2"}).Do()
    fmt.Printf("err = %v\n", err)
}

```
### header
#### resty-header
```go
func restyHeader() {
    client := resty.New()

    // slice or array
    _, err := client.R().SetHeader("k1", "v1").SetHeader("k2", "v2").Get("http://127.0.0.1:8080/header")
    fmt.Printf("err = %v\n", err)

    // map
    _, err = client.R().SetHeaders(map[string]string{"k1": "v1", "k2": "v2"}).Get("http://127.0.0.1:8080/header")
    fmt.Printf("err = %v\n", err)

    // struct 不支持
}

```
#### gout-header
```go
func goutHeader() {

    // slice or array
    err := gout.GET(":8080/header").SetHeader([]string{"k1", "v1", "k2", "v2"}).Do()
    fmt.Printf("err = %v\n", err)

    // map
    err = gout.GET(":8080/header").SetHeader(gout.H{"k1": "v1", "k2": "v2"}).Do()
    fmt.Printf("err = %v\n", err)

    // struct
    // 注意结构体要加上header标签
    err = gout.GET(":8080/header").SetHeader(testData{K1: "v1", K2: "v2"}).Do()
    fmt.Printf("err = %v\n", err)
}

```

### json

#### resty-json
```go
func restyJSON() {
    r := testJSON{}
    client := resty.New()
    resp, err := client.R().SetBody(map[string]interface{}{"key": "val"}).SetResult(&r).Post("http://127.0.0.1:8080/json")

    fmt.Printf("err = %v:result%v, resp = %v\n", err, r, resp)
}


```
#### gout-json
```go
func goutJSON() {
    r := testJSON{}
    err := gout.POST(":8080/json").SetJSON(gout.H{"key": "val"}).BindJSON(&r).Do()

    fmt.Printf("err = %v:result%v\n", err, r)
}


```
### body
#### resty-body
```go
func restyBody() {
    client := resty.New()
    resp, err := client.R().SetBody("hello").Post("http://127.0.0.1:8080/body")
    fmt.Println(resp.String(), err)
}
```
#### gout-body
```go
func goutBody() {
    s := ""
    err := gout.POST(":8080/body").SetBody("hello").BindBody(&s).Do()
    fmt.Println(s, err)
}

```
### test server
```go
func router() {
    g := gin.Default()
    g.POST("/json", func(c *gin.Context) {
        t := testData{}
        c.BindJSON(&t)
        c.JSON(200, t)
    })

    g.GET("/query", func(c *gin.Context) {
        t := testData{}
        c.BindQuery(&t)
        fmt.Printf("qeury result %v\n", t)
        c.JSON(200, t)
    })

    g.GET("/header", func(c *gin.Context) {
        t := testData{}
        c.BindHeader(&t)
        fmt.Printf("header result %v\n", t)
        c.JSON(200, t)
    })

    g.POST("/body", func(c *gin.Context) {
        all, _ := ioutil.ReadAll(c.Request.Body)
        fmt.Println(all)
        c.String(200, *(*string)(unsafe.Pointer(&all)))
    })

    g.Run()
}

```