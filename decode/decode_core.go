package decode

import (
	"net/http"
	"reflect"
)

type setter interface {
	Set(key string,
		value reflect.Value,
		sf reflect.StructField,
		tagValue string) error
}

type decData map[string][]string

func (d decData) Set(
	key string,
	value reflect.Value,
	sf reflect.StructField,
	tagValue string) {

}

func decode(d decData, obj interface{}) {

	// elem pointer
}
