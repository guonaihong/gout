package core

import (
	"errors"
)

type ReadCloseFail struct{}

// 供测试使用
func (r *ReadCloseFail) Read(p []byte) (n int, err error) {
	return 0, errors.New("must fail")
}
func (r *ReadCloseFail) Close() error {
	return nil
}
