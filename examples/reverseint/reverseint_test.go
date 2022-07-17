package reverseint

import (
	"fmt"
	"testing"

	"github.com/esoytekin/executor"
	"gotest.tools/assert"
)

type ReverseTask struct {
	input int32
}

// Hash godoc
// is used for memoization
// generate it using input parameters
func (r *ReverseTask) Hash(f func(string) int) []int {
	return []int{int(r.input)}
}

func (r *ReverseTask) Exec() int32 {
	return Reverse(r.input)
}

func TestReverseInt(t *testing.T) {

	tests := []struct {
		input    int32
		expected int32
	}{
		{
			input:    14,
			expected: 41,
		},
		{
			input:    123,
			expected: 321,
		},
		{
			input:    -123,
			expected: -321,
		},
		{
			input:    120,
			expected: 21,
		},
	}

	exec := executor.NewTaskExecutor[int32](100)

	var tasks []executor.Task[int32]

	var expected []int32

	for _, ti := range tests {

		task := &ReverseTask{ti.input}

		tasks = append(tasks, task)

		expected = append(expected, ti.expected)

	}

	exec.Progress(func(x int) {
		fmt.Println("progress", x)
	})

	results := exec.ExecuteTask(tasks...)

	for x := range expected {
		assert.Equal(t, expected[x], results[x])
	}

}
