package bench

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func Test_Bench_newTemplate(t *testing.T) {
	tmpl := newTemplate()

	r := report{
		Failed:        3,
		Percentage55:  time.Second,
		Percentage66:  2 * time.Second,
		Percentage75:  3 * time.Second,
		Percentage80:  4 * time.Second,
		Percentage90:  5 * time.Second,
		Percentage95:  6 * time.Second,
		Percentage98:  7 * time.Second,
		Percentage99:  8 * time.Second,
		Percentage100: 9 * time.Second,
		StatusCodes: map[int]int{
			200: 100,
			500: 3,
		},
	}

	err := tmpl.Execute(os.Stdout, r)
	assert.NoError(t, err)
}
