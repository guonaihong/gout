package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

// 使用SetWWWForm 接口设置x-www-form-urlencoded格式数据
// 下面的代码有 map/array/struct 使用example
type testWWWForm struct {
	Int     int     `form:"int" www-form:"int"`
	Float64 float64 `form:"float64" www-form:"float64"`
	String  string  `form:"string" www-form:"string"`
}

func mapExample() {
	fmt.Printf("====1.===============www-form=====use gout.H==\n\n")
	// 1.第一种方式，使用gout.H
	err := gout.POST(":8080/post").
		Debug(true).
		SetWWWForm(gout.H{
			"int":     3,
			"float64": 3.14,
			"string":  "test-www-Form",
		}).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func arrayExample() {
	fmt.Printf("====2.===============www-form=====use gout.A==\n\n")
	// 2.第一种方式，使用gout.A
	err := gout.POST(":8080/post").
		Debug(true).
		SetWWWForm(gout.A{
			"int", 3,
			"float64", 3.14,
			"string", "test-www-Form",
		}).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func structExample() {
	fmt.Printf("====3.=================www-form=====use struct==\n\n")
	// 3.第一种方式，使用结构体
	need := testWWWForm{
		Int:     3,
		Float64: 3.14,
		String:  "test-www-Form",
	}

	err := gout.POST(":8080/post").
		Debug(true).
		SetWWWForm(need).Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}

func main() {

	go server()

	time.Sleep(time.Millisecond * 500)

	mapExample()
	arrayExample()
	structExample()
}

func server() {
	router := gin.New()

	router.POST("/post", func(c *gin.Context) {

		t := testWWWForm{}
		err := c.ShouldBind(&t)
		if err != nil {
			c.String(200, "demo fail")
			return
		}

		fmt.Printf("\n\nread client data#------->%#v\n\n", t)
		c.String(200, "I am the server response: www-form demo ok")
	})
	router.Run(":8080")
}
