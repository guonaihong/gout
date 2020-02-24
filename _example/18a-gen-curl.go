package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

func genFormCurl() {
	// 1.formdata
	err := gout.GET(":1234").
		SetForm(gout.A{"text", "good", "mode", "A", "voice", gout.FormFile("./t8.go")}).
		Export().Curl().Do()
	// output:
	// curl -X GET -F "text=good" -F "mode=A" -F "voice=@./voice" "http://127.0.0.1:1234"
}
func genJSONCurl() {
	// 2.json body
	err = gout.GET(":1234").
		SetJSON(gout.H{"key1": "val1", "key2": "val2"}).
		Export().Curl().Do()
	// output:
	// curl -X GET -H "Content-Type:application/json" -d "{\"key1\":\"val1\",\"key2\":\"val2\"}" "http://127.0.0.1:1234"

	fmt.Printf("%v\n", err)
}

func main() {
	genFormCurl()
	genJSONCurl()
}
