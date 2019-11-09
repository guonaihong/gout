package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"io"
	"os"
	"time"
)

func server() {
	router := gin.New()
	router.POST("/send/file", func(c *gin.Context) {
		var file bytes.Buffer
		io.Copy(&file, c.Request.Body)
		fmt.Printf("server:file size = %d\n", file.Len())
	})

	router.Run()
}

func main() {
	go server()
	time.Sleep(time.Millisecond * 500) //sleep下等服务端真正起好

	fmt.Printf("=====3.=====send file========\n")
	fd, err := os.Open("./14-send-file.go")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer fd.Close()

	err = gout.POST(":8080/send/file").
		SetBody(fd). // send file
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fi, err := fd.Stat()
	if err != nil {
		return
	}

	fmt.Printf("client:file size:%d\n", fi.Size())
}
