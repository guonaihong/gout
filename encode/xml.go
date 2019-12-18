package encode

import (
	"encoding/xml"
	"io"
)

// XMLEncode xml encoder structure
type XMLEncode struct {
	obj interface{}
}

// NewXMLEncode create a new xml encoder
func NewXMLEncode(obj interface{}) *XMLEncode {
	if obj == nil {
		return nil
	}

	return &XMLEncode{obj: obj}
}

// Encode xml encoder
func (x *XMLEncode) Encode(w io.Writer) error {
	encode := xml.NewEncoder(w)
	return encode.Encode(x.obj)
}

// Name xml Encoder name
func (x *XMLEncode) Name() string {
	return "xml"
}
