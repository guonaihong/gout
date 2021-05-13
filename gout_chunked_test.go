package gout

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"

	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
)

func testTcpSocket(out *bytes.Buffer, quit chan bool, t *testing.T) (addr string) {
	addr = core.GetNoPortExists()

	addr = ":" + addr
	go func() {
		defer close(quit)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			t.Errorf("%v\n", err)
			return
		}
		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			t.Errorf("%v\n", err)
			return
		}

		io.Copy(out, conn)
		conn.Close()
	}()

	return addr

}

func Test_Use_Chunked(t *testing.T) {
	var out bytes.Buffer
	quit := make(chan bool)

	addr := testTcpSocket(&out, quit, t)
	time.Sleep(time.Second / 100) //等待服务起好

	POST(addr).SetTimeout(time.Second / 100).Chunked().SetBody("11111111111").Do()
	<-quit
	//time.Sleep(time.Second)

	assert.NotEqual(t, bytes.Index(out.Bytes(), []byte("Transfer-Encoding: chunked")), -1)
}
