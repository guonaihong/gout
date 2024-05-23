package dataflow

import "testing"

func Benchmark_URL_Template(b *testing.B) {

	ts := createMethodEcho()
	tc := testURLTemplateCase{Host: ts.URL, Method: "get"}

	for n := 0; n < b.N; n++ {
		code := 0
		err := New().GET("{{.Host}}/{{.Method}}", tc).Code(&code).Do()
		if err != nil {
			panic(err)
		}
		if code != 200 {
			panic("code != 200")
		}
	}
}
