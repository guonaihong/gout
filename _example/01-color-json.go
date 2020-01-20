package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"time"
)

func server() {
	router := gin.New()
	router.POST("/colorjson", func(c *gin.Context) {
		c.JSON(200, gin.H{"str2": "str2 val", "int2": 2})
	})

	router.Run()
}

func useMap() {
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

func useStruct() {

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

func useString() {
	fmt.Printf("\n\n3.=============color json======string example=====\n\n")
	err := gout.POST(":8080/colorjson").
		Debug(true).
		SetJSON(query).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

func useBytes() {
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

	useMap()
	useStruct()
	useString()
	useBytes()
}
