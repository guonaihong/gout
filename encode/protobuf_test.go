package encode

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/guonaihong/gout/testdata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestProtoBuf_Name(t *testing.T) {
	assert.Equal(t, NewProtoBufEncode("").Name(), "protobuf")
}

func TestProtoBuf_Fail(t *testing.T) {

	out := bytes.Buffer{}
	for _, v := range []interface{}{testdata.Req{}} {
		p := NewProtoBufEncode(v)
		out.Reset()

		err := p.Encode(&out)
		assert.Error(t, err)
	}

}

func TestProtoBuf_Encode(t *testing.T) {
	data1, err1 := proto.Marshal(&testdata.Req{Seq: 1, Res: "fk"})
	assert.NoError(t, err1)

	data2 := &testdata.Req{Seq: 1, Res: "fk"}
	data := []interface{}{data1, data2, string(data1)}

	out := bytes.Buffer{}

	for i, v := range data {
		p := NewProtoBufEncode(v)
		out.Reset()

		assert.NoError(t, p.Encode(&out))

		got := testdata.Req{}

		err := proto.Unmarshal(out.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, got.Seq, int32(1), fmt.Sprintf("fail index:%d", i))
		assert.Equal(t, got.Res, "fk")
	}
}
