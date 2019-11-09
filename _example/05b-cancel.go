package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

func server() {
	router := gin.New()
	router.GET("/cancel", func(c *gin.Context) {
		ctx := c.Request.Context()
		select {
		case <-ctx.Done():
			fmt.Printf("cancel done\n")
		}
	})

	router.Run()
}

func main() {
	go server()
	time.Sleep(time.Millisecond)

	// 给http请求 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)

	cancel() //取消

	err := gout.GET(":8080/cancel").
		WithContext(ctx).
		Do()

	fmt.Printf("err = %s\n", err)
}
