package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

func setTimeoutExample() {
	// 给http请求 设置超时
	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)

	err := gout.GET(":8080/timeout").
		WithContext(ctx).
		Do()

	fmt.Printf("err = %s\n", err)
}

func main() {
	go server()
	time.Sleep(time.Millisecond)
	setTimeoutExample()
}

func server() {
	router := gin.New()
	router.GET("/timeout", func(c *gin.Context) {
		ctx := c.Request.Context()
		select {
		case <-ctx.Done():
			fmt.Printf("timeout done\n")
		}
	})

	router.Run()
}
