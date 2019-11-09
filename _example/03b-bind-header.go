package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

func server() {
	router := gin.New()
	router.GET("/test.header", func(c *gin.Context) {
		c.Writer.Header().Add("sid", "1234")
		c.Writer.Header().Add("total", "2048")
		c.Writer.Header().Add("time", time.Now().Format("2006-01-02"))
	})

	router.Run()
}

type rspHeader struct {
	Total int       `header:"total"`
	Sid   string    `header:"sid"`
	Time  time.Time `header:"time" time_format:"2006-01-02"`
}

func main() {
	go server()

	time.Sleep(time.Millisecond)

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
