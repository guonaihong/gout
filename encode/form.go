package encode

type FormEncode struct {
	body bytes.Buffer
}

func NewFormEncode() *FormEncode {
	return FormEncode{}
}

func (h *FormEncode) Add(key string, v reflect.Value, sf reflect.StructField) error {
}
