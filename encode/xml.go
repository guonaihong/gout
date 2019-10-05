package encode

import (
	"encoding/xml"
	"io"
)

type XMLEncode struct {
	obj interface{}
}

func NewXMLEncode(obj interface{}) *XMLEncode {
	if obj == nil {
		return nil
	}

	return &XMLEncode{obj: obj}
}

func (x *XMLEncode) Encode(w io.Writer) error {
	encode := xml.NewEncoder(w)
	return encode.Encode(x.obj)
}
