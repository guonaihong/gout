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

func rawhttp() {
	s := `POST /colorjson HTTP/1.1
Host: 127.0.0.1:8080
User-Agent: Go-http-client/1.1
Content-Length: 97
Content-Type: application/json
Accept-Encoding: gzip

{"array":["foo","bar","baz"],"bool":false,"null":null,"num":100,"obj":{"a":1,"b":2},"str":"foo"}
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 27 Dec 2019 05:36:27 GMT
Content-Length: 29

{"int2":2,"str2":"str2 val"}
    `
	err := gout.NewImport().RawText(s).Debug(true).SetURL(":1234/colorjson").Do()
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
