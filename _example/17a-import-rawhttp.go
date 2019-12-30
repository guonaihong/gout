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

	router.Run(":1234")
}

func rawhttp() {
	s := `POST /colorjson HTTP/1.1
Host: 127.0.0.1:8080
User-Agent: Go-http-client/1.1
Content-Length: 97
Content-Type: application/json
Accept-Encoding: gzip

{"array":["foo","bar","baz"],"bool":false,"null":null,"num":100,"obj":{"a":1,"b":2},"str":"foo"}
    `
	err := gout.NewImport().RawText(s).Debug(true).SetHost(":1234").Do()
	if err != nil {
		fmt.Printf("err = %s\n", err)
		return
	}
}

func main() {
	go server()
	time.Sleep(time.Millisecond * 200)
	rawhttp()
}
