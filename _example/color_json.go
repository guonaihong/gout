package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

func server() {
	router := gin.Default()
	router.POST("/colorjson", func(c *gin.Context) {
		c.JSON(200, gin.H{"str2": "str2 val", "int2": 2})
	})

	router.Run()
}
func main() {
	go server()

	time.Sleep(time.Millisecond * 500)

	err := gout.POST(":8080/colorjson").
		Debug(gout.DebugColor()).
		SetJSON(gout.H{
			"str":     "str val",
			"int":     3,
			"float64": 3.14}).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}
