package gout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewWithOptTopMethod(t *testing.T) {
	var total int32

	router := setupMethod(&total)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	err := NewWithOpt().GET(ts.URL + "/someGet").Do()
	assert.NoError(t, err)

	err = NewWithOpt().POST(ts.URL + "/somePost").Do()
	assert.NoError(t, err)

	err = NewWithOpt().PUT(ts.URL + "/somePut").Do()
	assert.NoError(t, err)

	err = NewWithOpt().DELETE(ts.URL + "/someDelete").Do()
	assert.NoError(t, err)

	err = NewWithOpt().PATCH(ts.URL + "/somePatch").Do()
	assert.NoError(t, err)

	err = NewWithOpt().HEAD(ts.URL + "/someHead").Do()
	assert.NoError(t, err)

	err = NewWithOpt().OPTIONS(ts.URL + "/someOptions").Do()
	assert.NoError(t, err)

	assert.Equal(t, int(total), 7)
}
