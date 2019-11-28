package bench

import (
	"net/http"
)

func cloneRequest(r *http.Request) (*http.Request, error) {
	var err error

	r0 := &http.Request{}
	*r0 = *r

	r0.Header = make(http.Header, len(r.Header))

	for k, h := range r.Header {
		r.Header[k] = append([]string(nil), h...)
	}

	r0.Body, err = r.GetBody()
	return r0, err
}
