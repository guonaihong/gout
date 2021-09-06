package dataflow

import (
	"path"
	"strings"
)

const (
	httpProto  = "http://"
	httpsProto = "https://"
)

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func join(elem ...string) (rv string) {

	defer func() {
		if strings.HasPrefix(rv, httpProto) || strings.HasPrefix(rv, httpsProto) {
			var proto, domainPath, query string
			protoEndIdx := strings.Index(rv, "//") + 2
			queryIdx := strings.Index(rv, "?")

			proto = rv[:protoEndIdx]
			if queryIdx != -1 {
				domainPath = rv[protoEndIdx:queryIdx]
				query = rv[queryIdx:]
			} else if rv != proto {
				domainPath = rv[protoEndIdx:]
			}
			if domainPath != "" {
				domainPath = path.Clean(domainPath)
			}
			rv = proto + domainPath + query
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
