package dataflow

import "testing"

func Benchmark_URL_Template(b *testing.B) {

	ts := createMethodEcho()
	tc := testURLTemplateCase{Host: ts.URL, Method: "get"}

	for n := 0; n < b.N; n++ {
		code := 0
		New().GET("{{.Host}}/{{.Method}}", tc).Code(&code).Do()
		if code != 200 {
			panic("code != 200")
		}
	}
}
