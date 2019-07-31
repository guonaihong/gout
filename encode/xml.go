package encode

import (
	"encoding/xml"
	"io"
)

type XmlEncode struct {
	obj interface{}
}

func NewXmlEncode(obj interface{}) *XmlEncode {
	if obj == nil {
		return nil
	}

	return &XmlEncode{obj: obj}
}

func (x *XmlEncode) Encode(w io.Writer) error {
	encode := xml.NewEncoder(w)
	return encode.Encode(x.obj)
}
