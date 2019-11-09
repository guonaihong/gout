package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"io/ioutil"
	"time"
)

func server() {
	r := gin.New()

	r.POST("/", func(c *gin.Context) {
		all, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(200, "fail")
			return
		}

		c.String(200, string(all))
	})

	r.Run()
}

func main() {
	go server()                        // 起测试服务
	time.Sleep(time.Millisecond * 500) //sleep下等服务端真正起好

	err := gout.POST(":8080/").
		Debug(true).
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

	err = gout.POST(":8080/").
		Debug(gout.NoColor()).
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
