package executor

import (
	"fmt"
	"testing"
	"time"

	"gotest.tools/assert"
)

type TaskImpl struct {
	input int
}

func (t *TaskImpl) Exec() interface{} {
	time.Sleep(3 * time.Second)
	return t.input * t.input
}

func (t *TaskImpl) Hash(f func(string) int) []int {
	return []int{t.input}
}

func TestLimitedThread(t *testing.T) {

	expected := []int{1, 4, 4, 9, 16, 25, 36, 49, 64, 81, 4}

	e := NewTaskExecutor(100)

	e.Progress(func(p int) {
		fmt.Println("progress", p)
	})

	results := e.ExecuteTask(&TaskImpl{1}, &TaskImpl{2}, &TaskImpl{2}, &TaskImpl{3}, &TaskImpl{4}, &TaskImpl{5}, &TaskImpl{6}, &TaskImpl{7}, &TaskImpl{8}, &TaskImpl{9}, &TaskImpl{2})

	// assert.DeepEqual(t, expected, results)

	assert.Equal(t, len(expected), len(results))

	for x := range results {
		assert.Equal(t, expected[x], results[x].(int))
	}

}
