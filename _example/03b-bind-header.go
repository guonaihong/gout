package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

// ============== 解析http header
// 使用BindHeader接口解析http header, 基本数据类型可自动绑定
type rspHeader struct {
	Total int       `header:"total"`
	Sid   string    `header:"sid"`
	Time  time.Time `header:"time" time_format:"2006-01-02"`
}

func bindHeader() {
	rsp := rspHeader{}
	err := gout.GET(":8080/test.header").
		Debug(true).
		BindHeader(&rsp). //解析请求header
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("rsp header:\n%#v \nTime:%s\n", rsp, rsp.Time)
}

func main() {
	go server()

	time.Sleep(time.Millisecond)
	bindHeader()
}

func server() {
	router := gin.New()
	router.GET("/test.header", func(c *gin.Context) {
		c.Writer.Header().Add("sid", "1234")
		c.Writer.Header().Add("total", "2048")
		c.Writer.Header().Add("time", time.Now().Format("2006-01-02"))
	})

	router.Run()
}
