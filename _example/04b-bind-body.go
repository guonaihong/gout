package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

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

func main() {
	go server()

	time.Sleep(time.Millisecond)

	// 1.解析string
	fmt.Printf("\n\n=========1. bind string=====\n")
	s := ""
	err := gout.GET(":8080/rsp/body/string").
		Debug(gout.DebugColor()).
		BindBody(&s).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("need(string) got(%s)\n", s)

	// 2.解析[]byte
	fmt.Printf("\n\n=========2. bind string=====\n")
	var b []byte
	err = gout.GET(":8080/rsp/body/bytes").
		Debug(gout.DebugColor()).
		BindBody(&b).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("need(bytes) got(%s)\n", b)

	// 3.解析int
	fmt.Printf("\n\n=========3. bind int=====\n")
	i := 0
	err = gout.GET(":8080/rsp/body/int").
		Debug(gout.DebugColor()).
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
