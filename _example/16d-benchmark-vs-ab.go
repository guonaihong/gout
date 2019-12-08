package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	benchCount      = 100000
	benchConcurrent = 30
)

func server() {
	router := gin.New()
	router.POST("/", func(c *gin.Context) {
		c.String(200, "hello world:gout")
	})

	router.Run()
}

func runGout() {
	fd, err := os.Open("../testdata/voice.pcm")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer fd.Close()

	err = gout.
		POST(":8080").
		SetBody(fd).                 //设置请求body内容
		Filter().                    //打开过滤器
		Bench().                     //选择bench功能
		Concurrent(benchConcurrent). //并发数
		Number(benchCount).          //压测次数
		Do()

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

// sudo apt install apache2-utils
var abCmd = fmt.Sprintf(`ab -c %d -n %d -p ../testdata/voice.pcm http://127.0.0.1:8080/`, benchConcurrent, benchCount)

func runAb() {
	out, err := exec.Command("bash", "-c", abCmd).Output()
	if err != nil {
		log.Fatal("%s\n", err)
	}
	fmt.Printf("%s\n", out)

}

func main() {
	go server()
	time.Sleep(300 * time.Millisecond)

	// 设为false，可看ab性能
	startGout := true
	if startGout {
		runGout()
	} else {
		runAb()
	}
}
