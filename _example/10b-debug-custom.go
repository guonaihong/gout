package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"io/ioutil"
	"os"
	"time"
)

// 自定义debug example，下面使用环境变量输出日志输出
// 日志输出功能，使用环境变量打开
func IOSDebug() gout.DebugOpt {
	return gout.DebugFunc(func(o *gout.DebugOption) {
		if len(os.Getenv("IOS_DEBUG")) > 0 {
			o.Debug = true
			o.Color = true //打开颜色高亮
		}
	})
}

func customExample() {
	err := gout.POST(":8080/").
		Debug(IOSDebug()).
		SetJSON(gout.H{"str": "foo",
			"num":   100,
			"bool":  false,
			"null":  nil,
			"array": gout.A{"foo", "bar", "baz"},
			"obj":   gout.H{"a": 1, "b": 2},
		}).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

// 运行 example(其中env IOS_DEBUG=on 用于设置环境变量)
// env IOS_DEBUG=on go run 10b-debug-custom.go
func main() {
	go server()                        // 起测试服务
	time.Sleep(time.Millisecond * 500) //sleep下等服务端真正起好

	customExample()
}

func server() {
	r := gin.New()

	r.POST("/", func(c *gin.Context) {
		all, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(200, "fail")
			return
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.String(200, string(all))
	})

	r.Run()
}
