package color

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_ColorCore_NewFormatter(t *testing.T) {
}

func Test_ColorCore_Marshal(t *testing.T) {
	str := `{
      "str": "foo",
      "num": 100,
      "bool": false,
      "null": null,
      "array": ["foo", "bar", "baz"],
      "obj": { "a": 1, "b": 2 }
    }`

	NoColor = false
	defer func() { NoColor = true }()
	var obj map[string]interface{}
	json.Unmarshal([]byte(str), &obj)

	// Marshall the Colorized JSON
	s, _ := Marshal(obj)
	fmt.Println(string(s))
}
