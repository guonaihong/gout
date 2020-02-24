package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

// ================SetJSON=接口用法===============
// gout使用SetJSON函数发送json请求至服务端
// 亮点有SetJSON函数支持多种数据类型，map/struct/array/string/bytes
// 下面的xxxExample对应该类型的用法

func mapExample() {
	fmt.Printf("\n\n1.=============color json======map example=====\n\n")
	err := gout.POST(":8080/colorjson").
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
}

func structExample() {

	type req struct {
		Str  string `json:"str"`
		Num  int    `json:"num"`
		Bool bool   `json:"bool"`
		Null *int   `json:"null"`
	}

	fmt.Printf("\n\n2.=============color json======struct example=====\n\n")
	err := gout.POST(":8080/colorjson").
		Debug(true).
		SetJSON(req{Str: "foo",
			Num:  100,
			Bool: false,
			Null: nil,
		}).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

var query = `
{
  "query": {
    "bool": {
      "must": [
        {
          "exists": {
            "field": "voice"
          }
        },
        {
          "match": {
              "errcode": 3
          }
        },
        {
          "range": {
            "time": {
              "lt": "2020-01-13T23:16:04+08:00",
              "gt": "2020-01-13T00:00:00+08:00"
            }
          }
        }
      ]
    }
  }
}
`

func stringExample() {
	fmt.Printf("\n\n3.=============color json======string example=====\n\n")
	err := gout.POST(":8080/colorjson").
		Debug(true).
		SetJSON(query).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

func bytesExample() {
	fmt.Printf("\n\n4.=============color json======bytes example=====\n\n")
	err := gout.POST(":8080/colorjson").
		Debug(true).
		SetJSON(query).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

func main() {
	go server()

	time.Sleep(time.Millisecond * 200)

	mapExample()
	structExample()
	stringExample()
	bytesExample()
}

func server() {
	router := gin.New()
	router.POST("/colorjson", func(c *gin.Context) {
		c.JSON(200, gin.H{"str2": "str2 val", "int2": 2})
	})

	router.Run()
}
