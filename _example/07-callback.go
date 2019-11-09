package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"math/rand"
	"time"
)

type Result struct {
	Errmsg  string `json:"errmsg"`
	ErrCode int    `json:"errcode"`
}

// 模拟 API网管
func server() {
	router := gin.New()

	router.GET("/", func(c *gin.Context) {

		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(2) //使用随机函数模拟某个服务有一定概率出现404
		switch x {
		case 0: // 模拟 404 找不到资源
			c.String(404, "<html> not found </html>")
		case 1:
			c.JSON(200, Result{Errmsg: "ok"})
		}
	})

	router.Run()
}

func main() {
	go server()                        // 等会起测试服务
	time.Sleep(time.Millisecond * 500) //用时间做个等待同步

	r, str404 := Result{}, ""
	code := 0

	err := gout.GET(":8080").Code(&code).Callback(func(c *gout.Context) (err error) {

		switch c.Code {
		case 200: //http code 200是json
			err = c.BindJSON(&r)
		case 404: // http code 404 是html
			err = c.BindBody(&str404)
		}
		return

	}).Do()

	if err != nil {
		fmt.Printf("err = %s\n", err)
		return
	}

	fmt.Printf("http code = %d, str404(%s), result(%v)\n", code, str404, r)
}
