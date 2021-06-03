package gout

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type chanWriter chan string

func (w chanWriter) Write(p []byte) (n int, err error) {
	w <- string(p)
	return len(p), nil
}

func Test_WithInsecureSkipVerify(t *testing.T) {

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}))

	errc := make(chanWriter, 10) // but only expecting 1
	ts.Config.ErrorLog = log.New(errc, "", 0)

	defer ts.Close()

	c := ts.Client()
	for _, insecure := range []bool{true, false} {
		var opts []Option
		if insecure {
			opts = []Option{WithClient(c), WithInsecureSkipVerify()}
		}
		client := NewWithOpt(opts...)
		err := client.GET(ts.URL).Do()
		if (err == nil) != insecure {
			t.Errorf("#insecure=%v: got unexpected err=%v", insecure, err)
		}
	}

	select {
	case v := <-errc:
		if !strings.Contains(v, "TLS handshake error") {
			t.Errorf("expected an error log message containing 'TLS handshake error'; got %q", v)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("timeout waiting for logged error")
	}
}
