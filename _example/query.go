package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

// struct里面的form tag是gin用来绑定数据用的
// testQuery
type testQuery struct {
	Q1 string    `query:"q1" form:"q1"`
	Q2 int       `query:"q2" form:"q2"`
	Q3 float32   `query:"q3" form:"q3"`
	Q4 float64   `query:"q4" form:"q4"`
	Q5 time.Time `query:"q5" form:"q5" time_format:"unix" time_location:"Asia/Shanghai"`
	Q6 time.Time `query:"q6" form:"q6" time_format:"unixNano" time_location:"Asia/Shanghai"`
	Q7 time.Time `query:"q7" form:"q7" time_format:"2006-01-02" time_location:"Asia/Shanghai"`
}

func server() {
	router := gin.New()
	router.GET("/test.query", func(c *gin.Context) {
		q2 := testQuery{}
		err := c.BindQuery(&q2)
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
	fmt.Printf("======1. SetQuery======use gout.H=====\n")
	err := gout.GET(":8080/test.query").
		Debug(gout.DebugColor()).
		SetQuery(gout.H{"q1": "v1",
			"q2": 2,
			"q3": float32(3.14),
			"q4": 4.56,
			"q5": time.Now().Unix(),
			"q6": time.Now().UnixNano(),
			"q7": time.Now().Format("2006-01-02")}).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// 2.使用数组变量
	fmt.Printf("======2. SetQuery======use array=====\n")
	err = gout.GET(":8080/test.query").
		Debug(gout.DebugColor()).
		SetQuery(gout.A{"q1", "v1",
			"q2", 2,
			"q3", float32(3.14),
			"q4", 4.56,
			"q5", time.Now().Unix(),
			"q6", time.Now().UnixNano(),
			"q7", time.Now().Format("2006-01-02")}).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// 3.使用结构体
	// 使用结构体需要设置query tag
	fmt.Printf("======2. SetQuery======use struct=====\n")
	err = gout.GET(":8080/test.query").
		Debug(gout.DebugColor()).
		SetQuery(testQuery{Q1: "v1",
			Q2: 2,
			Q3: float32(3.14),
			Q4: 4.56,
			Q5: time.Now(),
			Q6: time.Now(),
			Q7: time.Now()}).
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// 4.使用string
	fmt.Printf("======2. SetQuery======use string=====\n")
	err = gout.GET(":8080/test.query").
		Debug(gout.DebugColor()).
		SetQuery("q1=v1&q2=2&q3=3.14&q4=3.1415&q5=1564295760&q6=1564295760000001000&q7=2019-07-28").
		Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}
