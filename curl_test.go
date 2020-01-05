package gout

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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

		g := POST("www.qq.com")
		if v.flags&testCurlHeader > 0 {
			g.SetHeader(A{"h1", "hv1", "h2", "hv2"})
		}

		if v.flags&testCurlQuery > 0 {
			g.SetQuery(A{"q1", "qv1", "q2", "qv2"})
		}
		if v.flags&testCurlForm > 0 {
			g.SetForm(
				A{
					"mode", "A",
					"text", "good",
					"voice", FormFile("./testdata/voice.pcm")},
			)
		}

		if v.flags&testJSON > 0 {
			g.SetJSON(
				H{
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