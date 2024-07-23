package gout

import "testing"

func Test_Issue_389(t *testing.T) {
	_ = NewWithOpt(WithProxy("http://127.0.0.1:7897"), WithInsecureSkipVerify())

	// client.GET("https://api.ipify.org?format=json").Debug(true).Do()
}
