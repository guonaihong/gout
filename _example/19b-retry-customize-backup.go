package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/core"
	"github.com/guonaihong/gout/filter"
	"time"
)

func server() {
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"result": "test ok"})
	})

	router.Run(":1234")
}

func useRetryFunc() {
	// 获取一个没有服务绑定的端口
	port := core.GetNoPortExists()
	s := ""

	err := gout.GET(":" + port).Debug(true).BindBody(&s).F().
		Retry().Attempt(3).WaitTime(time.Millisecond * 10).MaxWaitTime(time.Millisecond * 50).
		Func(func(c *gout.Context) error {
			if c.Error != nil {
				c.SetHost(":1234") //必须是存在的端口
				return filter.ErrRetry
			}
			return nil

		}).Do()
	fmt.Printf("err = %v\n", err)
}

func main() {
	go server()
	time.Sleep(time.Millisecond * 200)
	useRetryFunc()
}
