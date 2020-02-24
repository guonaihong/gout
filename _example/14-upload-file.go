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

func uploadExample() {
	fmt.Printf("=====3.=====upload file========\n")
	fd, err := os.Open("./14-upload-file.go")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer fd.Close()

	err = gout.POST(":8080/upload/file").
		SetBody(fd). // upload file
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

func main() {
	go server()
	time.Sleep(time.Millisecond * 500) //sleep下等服务端真正起好

	uploadExample()
}

func server() {
	router := gin.New()
	router.POST("/upload/file", func(c *gin.Context) {
		var file bytes.Buffer
		io.Copy(&file, c.Request.Body)
		fmt.Printf("server:file size = %d\n", file.Len())
	})

	router.Run()
}
