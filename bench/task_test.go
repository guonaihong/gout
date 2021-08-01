package bench

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testTask struct {
	number int32
}

func (t *testTask) Init() {
}

func (t *testTask) Process(work chan struct{}) {
	for range work {
		atomic.AddInt32(&t.number, 1)
	}
}

func (t *testTask) Cancel() {
}

func (t *testTask) WaitAll() {
}

// 测试发送次数
func Test_Bench_Task_number(t *testing.T) {
	task := Task{
		Number: 1000,
	}

	testTask := &testTask{}
	task.Run(testTask)

	assert.Equal(t, int32(task.Number), int32(testTask.number))
}

// 测试发送duration
func Test_Bench_Task_duration(t *testing.T) {
	task := Task{
		Duration: time.Millisecond * 300,
	}

	s := time.Now()
	testTask := &testTask{}
	task.Run(testTask)

	e := time.Since(s)

	assert.LessOrEqual(t, int64(e), int64(task.Duration+100*time.Millisecond))
	assert.GreaterOrEqual(t, int64(e), int64(200*time.Millisecond))
}

// 测试发送频率
func Test_Bench_Task_rate(t *testing.T) {
	task := Task{
		Number: 300,
		Rate:   1000,
	}

	s := time.Now()
	testTask := &testTask{}
	task.Run(testTask)

	e := time.Now().Sub(s)

	assert.LessOrEqual(t, int64(e), int64(400*time.Millisecond))
	assert.GreaterOrEqual(t, int64(e), int64(200*time.Millisecond))
}

// TODO:测试卡住的情况
func Test_Bench_Task_block(t *testing.T) {
}
