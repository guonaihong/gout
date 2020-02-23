package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"strings"
	"time"
)

// 写入非结构化数据至http.body里，使用SetBody接口
// 当然结构化数据主要有json/xml/yaml
func stringExample() {
	// 1.发送string
	fmt.Printf("\n\n=====1.=====send string========\n")
	err := gout.POST(":8080/req/body").
		Debug(true).
		SetBody("send string"). // string
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func bytesExample() {
	// 2.发送[]byte
	fmt.Printf("=====2.=====send string========\n")
	err := gout.POST(":8080/req/body").
		Debug(true).
		SetBody([]byte("send bytes")). // []byte
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func ReaderExample() {
	// 3.发送实现io.Reader接口的变量
	fmt.Printf("=====3.=====io.Reader========\n")
	err := gout.POST(":8080/req/body").
		Debug(true).
		SetBody(strings.NewReader("io.Reader")). // io.Reader
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func baseTypeExample() {
	// 4.发送基础类型的变量
	fmt.Printf("=====4.=====base type========\n")
	err := gout.POST(":8080/req/body").
		Debug(true).
		SetBody(3.14). //float64
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	//SetBody支持的更多基础类型有int, int8, int16, int32, int64
	//uint, uint8, uint16, uint32, uint64
	//float32, float64
}

func main() {
	go server()

	time.Sleep(time.Millisecond)

	stringExample()
	bytesExample()
	ReaderExample()
	baseTypeExample()
}

func server() {
	router := gin.New()
	router.POST("/req/body", func(c *gin.Context) {
		c.String(200, "ok")
	})

	router.Run()
}
