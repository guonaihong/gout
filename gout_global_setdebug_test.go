package gout

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Global_SetDebug(t *testing.T) {
	router := setupDataFlow(t)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		assert.NoError(t, err)
		outC <- buf.String()
	}()

	// reading our temp stdout
	// 只设置timeout
	SetDebug(true) //设置全局超时时间
	err := GET(ts.URL + "/setdebug").Do()
	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	assert.NoError(t, err)
	assert.NotEqual(t, strings.Index(out, "setdebug"), -1)
}
