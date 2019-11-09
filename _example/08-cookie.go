package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"net/http"
	"time"
)

func server() {
	router := gin.Default()

	router.GET("/cookie", func(c *gin.Context) {

		cookie1, err := c.Request.Cookie("test1")
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		cookie2, err := c.Request.Cookie("test2")
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		fmt.Printf("cookie1 = %v, cookie2 = %v\n", cookie1, cookie2)
	})

	router.GET("/cookie/one", func(c *gin.Context) {

		cookie3, err := c.Request.Cookie("test3")
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		fmt.Printf("cookie3 = %v\n", cookie3)
	})

	router.Run()
}

func main() {
	go server()                        // 起测试服务
	time.Sleep(time.Millisecond * 500) //sleep下等服务端真正起好

	// 发送两个cookie
	fmt.Printf("\n\n===1. Send two cookies=========\n")
	err := gout.GET(":8080/cookie").
		Debug(true).
		SetCookies(&http.Cookie{Name: "test1", Value: "test1"},
			&http.Cookie{Name: "test2", Value: "test2"}).
		Do()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 发送一个cookie
	fmt.Printf("\n\n===1. Send a cookies=========\n")
	err = gout.GET(":8080/cookie/one").
		Debug(true).
		SetCookies(&http.Cookie{Name: "test3", Value: "test3"}).
		Do()
	fmt.Println(err)
}
