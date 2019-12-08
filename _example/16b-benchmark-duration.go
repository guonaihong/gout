package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

const (
	benchTime       = 4 * time.Second
	benchConcurrent = 30
)

func server() {
	router := gin.New()
	router.POST("/", func(c *gin.Context) {
		c.String(200, "hello world:gout")
	})

	router.Run()
}

func main() {
	go server()
	time.Sleep(300 * time.Millisecond)

	err := gout.
		POST(":8080").
		SetJSON(gout.H{"hello": "world"}). //设置请求body内容
		Filter().                          //打开过滤器
		Bench().                           //选择bench功能
		Concurrent(benchConcurrent).       //并发数
		Durations(benchTime).              //压测时间
		Do()

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
