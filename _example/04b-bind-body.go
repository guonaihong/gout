package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

// 可以使用BindBody解析返回的非结构化http body(结构化指json/xml/yaml等)
// 可以做到基础类型自动绑定, 下面是string/[]byte/int的example
func bindString() {
	// 1.解析string
	fmt.Printf("\n\n=========1. bind string=====\n")
	s := ""
	err := gout.GET(":8080/rsp/body/string").
		Debug(true).
		BindBody(&s).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("need(string) got(%s)\n", s)
}

func bindBytes() {
	// 2.解析[]byte
	fmt.Printf("\n\n=========2. bind []byte=====\n")
	var b []byte
	err := gout.GET(":8080/rsp/body/bytes").
		Debug(true).
		BindBody(&b).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("need(bytes) got(%s)\n", b)
}

func bindInt() {
	// 3.解析int
	fmt.Printf("\n\n=========3. bind int=====\n")
	i := 0
	err := gout.GET(":8080/rsp/body/int").
		Debug(true).
		BindBody(&i).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("need(65535) got(%d)\n", i)
	//BindBody支持的更多基础类型有int, int8, int16, int32, int64
	//uint, uint8, uint16, uint32, uint64
	//float32, float64
}

func main() {
	go server()

	time.Sleep(time.Millisecond)

	bindString()
	bindBytes()
	bindInt()
}

func server() {
	router := gin.New()
	router.GET("/rsp/body/bytes", func(c *gin.Context) {
		c.String(200, "bytes")
	})

	router.GET("/rsp/body/string", func(c *gin.Context) {
		c.String(200, "string")
	})

	router.GET("/rsp/body/int", func(c *gin.Context) {
		c.String(200, "65535")
	})

	router.Run()
}
