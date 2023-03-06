package encode

import (
	"bytes"
	"testing"

	"github.com/guonaihong/gout/core"
	"github.com/guonaihong/gout/setting"
	"github.com/stretchr/testify/assert"
)

type testWWWForm struct {
	w    *WWWFormEncode
	in   interface{}
	need string
}

func Test_WWWForm_Encode(t *testing.T) {
	var out bytes.Buffer

	tests := []testWWWForm{
		{NewWWWFormEncode(setting.Setting{}), core.A{"k1", "v1", "k2", 2, "k3", 3.14}, "k1=v1&k2=2&k3=3.14"},
		{NewWWWFormEncode(setting.Setting{}), core.H{"k1": "v1", "k2": 2, "k3": 3.14}, "k1=v1&k2=2&k3=3.14"},
	}

	for _, v := range tests {
		assert.NoError(t, v.w.Encode(v.in))
		assert.NoError(t, v.w.End(&out))
		assert.Equal(t, out.String(), v.need)
		out.Reset()
	}

}

func Test_WWWForm_Name(t *testing.T) {
	assert.Equal(t, NewWWWFormEncode(setting.Setting{}).Name(), "www-form")
}

type CreateUserMetadataReqBody struct {
	Avatarurl string `www-form:"avatarurl"`
	Nickname  string `www-form:"nickname"`
}

func Test_WWWForm_Struct(t *testing.T) {

	data := CreateUserMetadataReqBody{Avatarurl: "www.hh.com", Nickname: "good"}
	enc := NewWWWFormEncode(setting.Setting{})
	err := enc.Encode(&data)
	assert.NoError(t, err)

	var out bytes.Buffer
	assert.NoError(t, enc.End(&out))
	pos := bytes.Index(out.Bytes(), []byte("avatarurl"))
	pos1 := bytes.Index(out.Bytes(), []byte("nickname"))
	assert.NotEqual(t, pos, -1)
	assert.NotEqual(t, pos1, -1)
}
