package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"log"
	"net/http"
)

func main() {
	c := &http.Client{}
	s := ""
	err := gout.
		New(c).
		GET("www.google.com").
		// 设置proxy服务地址
		SetSOCKS5("47.104.175.115:80").
		// 绑定返回数据到s里面
		BindBody(&s).
		Do()

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(s)
}
