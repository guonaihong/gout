package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/filter"
	"time"
)

var first bool

// mock 服务端函数
func server() {
	router := gin.New()

	router.GET("/code", func(c *gin.Context) {
		if first {
			c.JSON(209, gout.H{"resutl": "1.return 209"})
			first = false
		} else {
			c.JSON(200, gout.H{"resut": "2.return 200"})
		}
	})

	router.Run()
}

func useRetryFuncCode() {
	s := ""
	err := gout.GET(":8080/code").Debug(true).BindBody(&s).F().
		Retry().Attempt(3).WaitTime(time.Millisecond * 10).MaxWaitTime(time.Millisecond * 50).
		Func(func(c *gout.Context) error {
			if c.Error != nil || c.Code == 209 {
				return filter.ErrRetry
			}

			return nil

		}).Do()

	fmt.Printf("err = %v\n", err)
}

func main() {
	first = true
	go server()
	time.Sleep(time.Millisecond * 200)
	useRetryFuncCode()
}
