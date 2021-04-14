package encode

import (
	"errors"
	"io"
	"reflect"

	"github.com/guonaihong/gout/core"
	"google.golang.org/protobuf/proto"
)

var ErrNotProtobuf = errors.New("Not protobuf data")

type ProtoBufEncode struct {
	obj interface{}
}

func NewProtoBufEncode(obj interface{}) *ProtoBufEncode {
	return &ProtoBufEncode{obj: obj}
}

func (p *ProtoBufEncode) Encode(w io.Writer) (err error) {
	if v, ok := core.GetBytes(p.obj); ok {
		//TODO找一个检测protobuf数据格式的函数
		_, err = w.Write(v)
		return err
	}

	var m proto.Message
	var ok bool

	for i := 0; i < 1; i++ {
		m, ok = p.obj.(proto.Message)
		if !ok {
			objVal := reflect.ValueOf(p.obj)
			if !objVal.CanAddr() {
				return ErrNotProtobuf
			}

			objVal.Addr().Interface()

			continue
		}
		break
	}

	all, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	_, err = w.Write(all)
	return err
}

func (p *ProtoBufEncode) Name() string {
	return "protobuf"
}
