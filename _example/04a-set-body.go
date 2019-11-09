package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"strings"
	"time"
)

func server() {
	router := gin.New()
	router.POST("/req/body", func(c *gin.Context) {
		c.String(200, "ok")
	})

	router.Run()
}

func main() {
	go server()

	time.Sleep(time.Millisecond)

	// 1.发送string
	fmt.Printf("\n\n=====1.=====send string========\n")
	err := gout.POST(":8080/req/body").
		Debug(gout.DebugColor()).
		SetBody("send string"). // string
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// 2.发送[]byte
	fmt.Printf("=====2.=====send string========\n")
	err = gout.POST(":8080/req/body").
		Debug(gout.DebugColor()).
		SetBody([]byte("send bytes")). // []byte
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// 3.发送实现io.Reader接口的变量
	fmt.Printf("=====3.=====io.Reader========\n")
	err = gout.POST(":8080/req/body").
		Debug(gout.DebugColor()).
		SetBody(strings.NewReader("io.Reader")). // io.Reader
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// 4.发基础类型的变量
	fmt.Printf("=====4.=====base type========\n")
	err = gout.POST(":8080/req/body").
		Debug(gout.DebugColor()).
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
