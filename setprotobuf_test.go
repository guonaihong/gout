package gout

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/guonaihong/gout/testdata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/gin-gonic/gin"
)

func setupEcho(total *int32, t *testing.T) *gin.Engine {

	router := gin.Default()

	cb := func(c *gin.Context) {
		atomic.AddInt32(total, 1)
		_, err := io.Copy(c.Writer, c.Request.Body)
		assert.NoError(t, err)
	}

	router.GET("/echo", cb)

	return router
}

func Test_SetProtoBuf(t *testing.T) {
	total := int32(0)
	router := setupEcho(&total, t)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	data1, err := proto.Marshal(&testdata.Req{Seq: 1, Res: "fk"})
	assert.NoError(t, err)

	testCount := 0
	for i, v := range []interface{}{
		data1,
		&testdata.Req{Seq: 1, Res: "fk"},
	} {
		code := 0
		var gotRes []byte
		err := GET(ts.URL + "/echo").SetProtoBuf(v).BindBody(&gotRes).Code(&code).Do()
		assert.NoError(t, err)
		assert.Equal(t, code, 200)
		got := &testdata.Req{}

		err = proto.Unmarshal(gotRes, got)
		assert.NoError(t, err)

		assert.Equal(t, got.Seq, int32(1), fmt.Sprintf("fail index:%d", i))
		assert.Equal(t, got.Res, "fk")
		testCount++
	}

	assert.Equal(t, total, int32(testCount))

}
