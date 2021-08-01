package encode

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testXML struct {
	I int     `xml:"i"`
	F float64 `xml:"f"`
	S string  `xml:"s"`
}

func TestNewXMLEncode(t *testing.T) {
	x := NewXMLEncode(nil)
	assert.Nil(t, x)
}

func TestXMLEncode_Name(t *testing.T) {
	assert.Equal(t, NewXMLEncode("").Name(), "xml")
}

func TestXMLEncode_Encode(t *testing.T) {
	need := testXML{
		I: 100,
		F: 3.14,
		S: "test encode xml",
	}

	out := bytes.Buffer{}

	x := `
<testXML><i>100</i><f>3.14</f><s>test encode xml</s></testXML>
	`
	data := []interface{}{need, &need, x}
	for _, v := range data {
		x := NewXMLEncode(v)
		out.Reset()

		assert.NoError(t, x.Encode(&out))

		got := testXML{}

		err := xml.Unmarshal(out.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, got, need)
	}

	// test fail
	for _, v := range []interface{}{`<testxml>`} {
		x := NewXMLEncode(v)
		out.Reset()
		err := x.Encode(&out)
		assert.Error(t, err)
	}
}
