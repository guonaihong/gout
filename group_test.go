package gout

import (
	"github.com/gin-gonic/gin"
	"sync/atomic"
	"testing"
	"time"
)

func TestGroupNew(t *testing.T) {
	total := int32(0)

	var send bool

	go func() {
		router := gin.Default()

		postFunc := func(c *gin.Context) {
			atomic.AddInt32(&total, 1)
		}

		// Simple group: v1
		v1 := router.Group("/v1")
		{
			v1.POST("/login", postFunc)
			v1.POST("/submit", postFunc)
			v1.POST("/read", postFunc)
		}

		// Simple group: v2
		v2 := router.Group("/v2")
		{
			v2.POST("/login", postFunc)
			v2.POST("/submit", postFunc)
			v2.POST("/read", postFunc)
		}

		v2_1 := v2.Group("/v1")
		{
			v2_1.POST("/login", func(c *gin.Context) {
				send = true
				atomic.AddInt32(&total, 1)
			})
		}
		router.Run(":8080")
	}()

	time.Sleep(250 * time.Millisecond)

	out := New(nil)

	// http://127.0.0.1:80/v1
	v1 := out.Group(":8080/v1")
	err := v1.POST("/login").Next().
		POST("/submit").Next().
		POST("/read").Do()

	if err != nil {
		t.Errorf("error:%s\n", err)
	}

	// http://127.0.0.1:80/v2
	v2 := out.Group(":8080/v2")
	err = v2.POST("/login").Next().
		POST("/submit").Next().
		POST("/read").Do()

	if err != nil {
		t.Errorf("error:%s\n", err)
	}

	v2_1 := v2.Group("/v1")
	err = v2_1.POST("/login").Do()

	if err != nil {
		t.Errorf("error:%s\n", err)
	}

	if total != 7 {
		t.Errorf("got %d want 7\n", total)
	}

	if !send {
		t.Errorf("/v2/v1/login fail\n")
	}
}

type data struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}

type BindTest struct {
	InBody   interface{}
	OutBody  interface{}
	httpCode int
}

func TestShouldBindJSON(t *testing.T) {
	var d3 data
	go func() {
		router := gin.Default()

		router.POST("/test.json", func(c *gin.Context) {
			c.ShouldBindJSON(&d3)
			c.JSON(200, d3)
		})

		// Listen and serve on 0.0.0.0:8080
		router.Run(":8080")
	}()

	time.Sleep(250 * time.Millisecond)
	g := New(nil)

	tests := []BindTest{
		{InBody: data{Id: 9, Data: "早上测试结构体"}, OutBody: data{}},
		{InBody: H{"id": 10, "data": "早上测试map"}, OutBody: data{}},
	}

	for k, _ := range tests {
		t.Logf("outbody type:%T:%p\n", tests[k].OutBody, &tests[k].OutBody)

		err := g.POST(":8080/test.json").ToJSON(&tests[k].InBody).ShouldBindJSON(&tests[k].OutBody).Code(&tests[k].httpCode).Do()
		if err != nil {
			t.Errorf("send fail:%s\n", err)
		}

		if tests[k].httpCode != 200 {
			t.Errorf("got:%d, want:200\n", tests[k].httpCode)
		}

		if tests[k].OutBody.(map[string]interface{})["id"].(float64) != float64(d3.Id) {
			t.Errorf("got:%#v(P:%p), want:%#v(T:%T)\n", tests[k].OutBody, &tests[k].OutBody, tests[k].InBody, tests[k].InBody)
		}

		/*
			if !reflect.DeepEqual(&tests[k].InBody, &tests[k].OutBody) {
				t.Errorf("got:%#v(P:%p), want:%#v(T:%T)\n", tests[k].OutBody, &tests[k].OutBody, tests[k].InBody, tests[k].InBody)
			}
		*/

	}
}
