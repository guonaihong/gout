package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

type data struct {
	Id   int    `json:"id" xml:"id"`
	Data string `json:"data xml:"data""`
}

func server() {
	router := gin.Default()

	router.POST("/test.xml", func(c *gin.Context) {
		var d3 data
		err := c.BindXML(&d3)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		c.XML(200, d3)
	})

	router.Run()
}

func main() {
	var rsp data
	go server()
	time.Sleep(time.Millisecond * 500) //sleep下等服务端真正起好

	code := 200

	err := gout.POST(":8080/test.xml").
		Debug(true).
		SetXML(data{Id: 3, Data: "test data"}).
		BindXML(&rsp).
		Code(&code).
		Do()

	if err != nil || code != 200 {
		fmt.Printf("%v:%d\n", err, code)
	}
}
