package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

type testForm struct {
	Mode string `form:"mode"`
	Text string `form:"text"`
	//Voice []byte `form:"voice" form-mem:"true"` //todo open
}

type testForm2 struct {
	Mode   string `form:"mode"`
	Text   string `form:"text"`
	Voice  string `form:"voice" form-file:"true"` //从文件中读取
	Voice2 []byte `form:"voice2" form-mem:"true"` //从内存中构造
}

func server() {
	router := gin.New()
	router.POST("/test.form", func(c *gin.Context) {

		t2 := testForm{}
		err := c.Bind(&t2)
		if err != nil {
			fmt.Printf("err = %s\n", err)
			return
		}
	})

	router.Run()
}

func main() {
	go server()
	time.Sleep(time.Millisecond * 500) //sleep下等服务端真正起好

	// 1.使用gout.H
	fmt.Printf("\n\n====1. use gout.H==============\n\n")
	code := 0
	err := gout.
		POST(":8080/test.form").
		Debug(true).
		SetForm(gout.H{"mode": "A",
			"text":   "good",
			"voice":  gout.FormFile("../testdata/voice.pcm"),
			"voice2": gout.FormMem("pcm")}).
		Code(&code).
		Do()

	if err != nil || code != 200 {
		fmt.Printf("%s:code = %d\n", err, code)
		return
	}

	// 2.使用结构体里面的数据
	fmt.Printf("\n\n====2. use struct==============\n\n")
	err = gout.
		POST(":8080/test.form").
		Debug(true).
		SetForm(testForm2{
			Mode:   "A",
			Text:   "good",
			Voice:  "../testdata/voice.pcm",
			Voice2: []byte("pcm")}).
		Code(&code).Do()
	if err != nil || code != 200 {

	}
}
