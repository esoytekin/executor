# Parallel Task Executor for GO

Executes tasks in parallel with `MemoizableFuture` support.

## Install

`$ go get -u github.com/esoytekin/executor`

## Introduction

Executes tasks in parallel with `MemoizableFuture` support.
tasks should implement `executor.Task` interface

```go
type Task interface {
	// Exec godoc
	Exec() interface{}

	// Hash godoc
	// is used for memoization
	// generate it using input parameters
	Hash(hashStr func(string) int) []int
}

```

## executor.NewFutureTask

accepts `executor.Task` as parameter and returns new futureTask

```go
package main


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

func main(){

    taskItem := &TaskImpl{1}
    ft := executor.NewFutureTask(taskItem)
    result := ft.Get() // blocks execute operation
    fmt.Println("result", result)
}
```

## Usage Example

```go
package main

import (
	"fmt"
	"testing"
	"time"

	"gotest.tools/assert"
    	"github.com/esoytekin/executor"
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

func TestExecutor(t *testing.T) {

	expected := []int{1, 4, 4, 9, 16, 25, 36, 49, 64, 81, 4}

	e := executor.NewTaskExecutor(100)

	e.Progress(func(p int) {
		fmt.Println("progress", p)
	})

	results := e.ExecuteTask(&TaskImpl{1}, &TaskImpl{2}, &TaskImpl{2}, &TaskImpl{3}, &TaskImpl{4}, &TaskImpl{5}, &TaskImpl{6}, &TaskImpl{7}, &TaskImpl{8}, &TaskImpl{9}, &TaskImpl{2})

	assert.Equal(t, len(expected), len(results))

	for x := range results {
		assert.Equal(t, expected[x], results[x].(int))
	}

}

```
