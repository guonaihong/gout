package decode

import (
	"encoding/xml"
	"io"
)

type XmlDecode struct {
	obj interface{}
}

func NewXmlDecode(obj interface{}) *XmlDecode {
	if obj == nil {
		return nil
	}
	return &XmlDecode{obj: obj}
}

func (x *XmlDecode) Decode(r io.Reader) error {
	decode := xml.NewDecoder(r)
	return decode.Decode(x.obj)
}

func DecodeXML(r io.Reader, obj interface{}) error {
	decode := xml.NewDecoder(r)
	return decode.Decode(obj)
}
