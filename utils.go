package gout

import (
	"path"
	"strings"

	"github.com/guonaihong/gout/core"
)

const (
	httpProto  = "http://"
	httpsProto = "https://"
)

// ReadCloseFail required for internal testing, external can be ignored
type ReadCloseFail = core.ReadCloseFail

// H is short for map[string]interface{}
type H = core.H

// A variable is short for []interface{}
type A = core.A

// FormFile encounter the FormFile variable in the form encoder will open the file from the file
type FormFile = core.FormFile

// FormMem when encountering the FormMem variable in the form encoder, it will be read from memory
type FormMem = core.FormMem

// FormType custom form names and stream types are available
type FormType = core.FormType

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func join(elem ...string) (rv string) {

	defer func() {
		if strings.HasPrefix(rv, httpProto) {
			rv = httpProto + path.Clean(rv[len(httpProto):])
			return
		}

		if strings.HasPrefix(rv, httpsProto) {
			rv = httpsProto + path.Clean(rv[len(httpsProto):])
			return
		}

		rv = path.Clean(rv)
	}()

	for i, e := range elem {
		if e != "" {
			return strings.Join(elem[i:], "/")
		}
	}
	return ""
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}
