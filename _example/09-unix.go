package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"net"
	"net/http"
	"os"
)

func server(path string) *http.Server {
	router := gin.Default()
	type testHeader struct {
		H1 string `header:"h1"`
		H2 string `header:"h2"`
	}

	router.POST("/test/unix", func(c *gin.Context) {

		tHeader := testHeader{}
		err := c.ShouldBindHeader(&tHeader)
		if err != nil {
			c.String(200, "fail")
			return
		}

		c.String(200, "ok")
	})

	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil
	}

	srv := http.Server{Handler: router}
	go func() {
		srv.Serve(listener)
	}()

	return &srv
}

func main() {
	path := "./unix.sock"
	defer os.Remove(path)

	ctx, cancel := context.WithCancel(context.Background())
	srv := server(path)
	defer func() {
		srv.Shutdown(ctx)
		cancel()
	}()

	c := http.Client{}
	s := ""
	err := gout.New(&c).
		Debug(true).
		UnixSocket(path).
		POST("http://xxx/test/unix/").
		SetHeader(gout.H{"h1": "v1", "h2": "v2"}).
		BindBody(&s).Do()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("result = %s\n", s)
}
