package export

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout/core"
	"github.com/guonaihong/gout/dataflow"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	testCurlHeader = 1 << iota
	testCurlForm
	testCurlQuery
	testLong
	testJSON
)

const noPortExists = 12345

type testCurl struct {
	flags int
	need  string
}

// 测试生成curl命令
func Test_Curl(t *testing.T) {

	tests := []testCurl{
		{testCurlHeader, `curl -X POST -H "H1:hv1" -H "H2:hv2" "http://www.qq.com"`},

		{testCurlHeader | testCurlQuery, `curl -X POST -H "H1:hv1" -H "H2:hv2" "http://www.qq.com?q1=qv1&q2=qv2"`},
		{testCurlHeader | testCurlQuery | testCurlForm, `curl -X POST -H "H1:hv1" -H "H2:hv2" -F "mode=A" -F "text=good" -F "voice=@./voice" "http://www.qq.com?q1=qv1&q2=qv2"`},
		{testCurlHeader | testLong, `curl --request POST --header "H1:hv1" --header "H2:hv2" --url "http://www.qq.com"`},
		{testCurlHeader | testCurlQuery | testLong, `curl --request POST --header "H1:hv1" --header "H2:hv2" --url "http://www.qq.com?q1=qv1&q2=qv2"`},
		{testCurlHeader | testCurlQuery | testCurlForm | testLong, `curl --request POST --header "H1:hv1" --header "H2:hv2" --form "mode=A" --form "text=good" --form "voice=@./voice" --url "http://www.qq.com?q1=qv1&q2=qv2"`},

		{testCurlHeader | testJSON, `curl -X POST -H "Content-Type:application/json" -H "H1:hv1" -H "H2:hv2" -d "{\"jk1\":\"jv1\"}" "http://www.qq.com"`},
		{testCurlHeader | testCurlQuery | testJSON, `curl -X POST -H "Content-Type:application/json" -H "H1:hv1" -H "H2:hv2" -d "{\"jk1\":\"jv1\"}" "http://www.qq.com?q1=qv1&q2=qv2"`},

		{testCurlHeader | testLong | testJSON, `curl --request POST --header "Content-Type:application/json" --header "H1:hv1" --header "H2:hv2" --data "{\"jk1\":\"jv1\"}" --url "http://www.qq.com"`},
		{testCurlHeader | testCurlQuery | testLong | testJSON, `curl --request POST --header "Content-Type:application/json" --header "H1:hv1" --header "H2:hv2" --data "{\"jk1\":\"jv1\"}" --url "http://www.qq.com?q1=qv1&q2=qv2"`},
	}

	for _, v := range tests {
		var buf strings.Builder

		g := dataflow.POST("www.qq.com")
		if v.flags&testCurlHeader > 0 {
			g.SetHeader(core.A{"h1", "hv1", "h2", "hv2"})
		}

		if v.flags&testCurlQuery > 0 {
			g.SetQuery(core.A{"q1", "qv1", "q2", "qv2"})
		}
		if v.flags&testCurlForm > 0 {
			g.SetForm(
				core.A{
					"mode", "A",
					"text", "good",
					"voice", core.FormFile("../testdata/voice.pcm")},
			)
		}

		if v.flags&testJSON > 0 {
			g.SetJSON(
				core.H{
					"jk1": "jv1",
				},
			)
		}

		c := g.Export().Curl()
		if v.flags&testLong > 0 {
			c.LongOption()
		}
		c.SetOutput(&buf)

		err := c.Do()
		assert.NoError(t, err)

		os.Remove(fmt.Sprintf("./voice"))
		for i := 0; i < 10; i++ {
			os.Remove(fmt.Sprintf("./voice.%d", i))
		}

		if err != nil {
			return
		}

		fmt.Printf("%s\n%s\n", buf.String(), v.need)
		assert.Equal(t, strings.TrimSpace(buf.String()), v.need)
	}

}

func Test_Curl_GenAndSend(t *testing.T) {
	// test ok
	type testData struct {
		A string
		B string
	}

	yes := false
	router := func(b *bool) *gin.Engine {
		router := gin.Default()

		router.POST("/test.json", func(c *gin.Context) {
			test := testData{}
			c.BindJSON(&test)
			*b = true
			c.JSON(200, gin.H{"1": "1"})
		})

		return router
	}(&yes)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	var out strings.Builder
	err := dataflow.POST(ts.URL + "/test.json").SetJSON(core.H{"a": "a", "b": "b"}).Export().Curl().SetOutput(&out).GenAndSend().Do()
	assert.NoError(t, err)
	assert.Equal(t, yes, true)
	need := fmt.Sprintf(`curl -X POST -H "Content-Type:application/json" -d "{\"a\":\"a\",\"b\":\"b\"}" "%s/test.json"`, ts.URL)

	assert.Equal(t, strings.TrimSpace(out.String()), need)

	// ==================================
	// test fail
	tests := []dataflow.Curl{
		dataflow.POST(fmt.Sprintf(":%d/test.json", noPortExists)).SetJSON(core.H{"a": "a", "b": "b"}).Export().Curl().GenAndSend(),
		dataflow.POST(ts.URL + "/test.json").Debug(true).SetJSON(core.H{"a": "a", "b": "b"}).BindBody(&testData{}).Export().Curl().GenAndSend(),
		dataflow.POST(ts.URL + "/test.json").Debug(true).SetBody(&testData{}).BindBody(&testData{}).Export().Curl().GenAndSend(),
	}

	for _, v := range tests {
		err := v.Do()
		assert.Error(t, err)
	}
}
