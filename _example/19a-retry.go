package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"time"
)

func main() {
	err := gout.HEAD("127.0.0.1:8080").
		Debug(true).                      //打开debug模式
		Filter().                         //打开过滤器.Filter()的简写是.F()
		Retry().                          //打开重试模式
		Attempt(5).                       //最多重试5次
		WaitTime(500 * time.Millisecond). //基本等待时间
		MaxWaitTime(3 * time.Second).     //最长等待时间
		Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}
