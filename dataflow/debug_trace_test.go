package dataflow

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func checkValue(b *bytes.Buffer) bool {
	info := &TraceInfo{}
	v := reflect.ValueOf(info)
	v = v.Elem()

	debugInfo := b.String()
	have := false
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" {
			continue
		}

		name := sf.Name
		if !strings.Contains(debugInfo, name) {
			return false
		}
		have = true
	}

	if !have {
		return have
	}

	return true
}

func Test_Debug_Trace(t *testing.T) {
	router := func() *gin.Engine {
		router := gin.Default()

		router.POST("/test.json", func(c *gin.Context) {
			c.String(200, "ok")
		})

		return router
	}()
	errs := []error{
		func() error {
			ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
			return New().GET(ts.URL).Debug(Trace()).Do()
		}(),
		func() error {
			ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
			var b bytes.Buffer
			custom := func() DebugOpt {
				return DebugFunc(func(o *DebugOption) {
					o.Color = true
					o.Trace = true
					o.Write = &b
				})
			}
			err := New().GET(ts.URL).Debug(custom()).Do()
			if err != nil {
				return err
			}

			if !checkValue(&b) {
				return errors.New("No caring fields")
			}
			return nil
		}(),
	}

	for id, e := range errs {
		assert.NoError(t, e, fmt.Sprintf("test case id:%d", id))
	}
}
