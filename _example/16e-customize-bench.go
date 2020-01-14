package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/filter"
	"sync/atomic"
	"time"
)

const (
	benchNumber     = 30000
	benchConcurrent = 30
)

func server() {
	router := gin.New()
	router.POST("/", func(c *gin.Context) {
		c.String(200, "hello world:gout")
	})

	router.Run()
}

func customize() {
	i := int32(0)

	err := filter.NewBench().
		Concurrent(benchConcurrent).
		Number(benchNumber).
		Loop(func(c *gout.Context) error {
			uid := uuid.New()
			id := atomic.AddInt32(&i, 1)
			c.POST(":8080").Debug(true).SetJSON(gout.H{"sid": uid.String(),
				"appkey": fmt.Sprintf("ak:%d", id),
				"text":   fmt.Sprintf("test text :%d", id)})
			return nil

		}).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

func main() {
	go server()
	time.Sleep(300 * time.Millisecond)

	customize()
}
