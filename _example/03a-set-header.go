package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

type testHeader struct {
	H1 string    `header:"h1"`
	H2 int       `header:"h2"`
	H3 float32   `header:"h3"`
	H4 float64   `header:"h4"`
	H5 time.Time `header:"h5" time_format:"unix"`
	H6 time.Time `header:"h6" time_format:"unixNano"`
	H7 time.Time `header:"h7" time_format:"2006-01-02"`
}

func server() {
	router := gin.New()
	router.GET("/test.header", func(c *gin.Context) {
		h2 := testHeader{}
		err := c.BindHeader(&h2)
		if err != nil {
			c.String(500, "fail")
			return
		}
	})

	router.Run()
}

func main() {
	go server()

	time.Sleep(time.Millisecond)
	// 1.使用gout.H
	fmt.Printf("======1. SetHeader======use gout.H=====\n")
	err := gout.GET(":8080/test.header").
		Debug(true).
		SetHeader(gout.H{"h1": "v1",
			"h2": 2,
			"h3": float32(3.14),
			"h4": 4.56,
			"h5": time.Now().Unix(),
			"h6": time.Now().UnixNano(),
			"h7": time.Now().Format("2006-01-02")}).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// 2.使用数组变量
	fmt.Printf("======2. SetHeader======use array=====\n")
	err = gout.GET(":8080/test.header").
		Debug(true).
		SetHeader(gout.A{"h1", "v1",
			"h2", 2,
			"h3", float32(3.14),
			"h4", 4.56,
			"h5", time.Now().Unix(),
			"h6", time.Now().UnixNano(),
			"h7", time.Now().Format("2006-01-02")}).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// 3.使用结构体
	// 使用结构体需要设置"header" tag
	fmt.Printf("======3. SetHeader======use struct=====\n")
	err = gout.GET(":8080/test.header").
		Debug(true).
		SetHeader(testHeader{H1: "v1",
			H2: 2,
			H3: float32(3.14),
			H4: 4.56,
			H5: time.Now(),
			H6: time.Now(),
			H7: time.Now()}).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}
