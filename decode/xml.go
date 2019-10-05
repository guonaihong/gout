package decode

import (
	"encoding/xml"
	"io"
)

type XMLDecode struct {
	obj interface{}
}

func NewXMLDecode(obj interface{}) *XMLDecode {
	if obj == nil {
		return nil
	}
	return &XMLDecode{obj: obj}
}

func (x *XMLDecode) Decode(r io.Reader) error {
	decode := xml.NewDecoder(r)
	return decode.Decode(x.obj)
}

func DecodeXML(r io.Reader, obj interface{}) error {
	decode := xml.NewDecoder(r)
	return decode.Decode(obj)
}
