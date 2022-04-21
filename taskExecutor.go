package executor

import (
	"sync"
	"time"
)

type taskChanItem struct {
	idx  int
	task Task
}

type taskExecutor struct {
	taskChan       chan taskChanItem
	d              chan bool
	fin            chan bool
	results        []interface{}
	threadLimit    chan struct{}
	locker         sync.Mutex
	operationCount int
	ProgressC      chan int
	cache          sync.Map
}

// NewSingleThreadTaskExecutor returns new TaskExecutor instance with single thread limit
func NewSingleThreadTaskExecutor() *taskExecutor {
	return NewTaskExecutor(1)
}

// NewTaskExecutor returns new TaskExecutor instance.
//
// threadLimit limits thread count.
func NewTaskExecutor(threadLimit int) *taskExecutor {
	taskChan := make(chan taskChanItem)
	done := make(chan bool)
	m := make(chan struct{}, threadLimit)
	progressC := make(chan int, 100)
	return &taskExecutor{
		taskChan:    taskChan,
		d:           done,
		fin:         make(chan bool),
		threadLimit: m,
		locker:      sync.Mutex{},
		ProgressC:   progressC,
	}
}

func (t *taskExecutor) init(tasks []Task) {
	taskLen := len(tasks)
	t.results = make([]interface{}, taskLen)

}

// ExecuteTask accepts a list of tasks and executes them in parallel.
func (t *taskExecutor) ExecuteTask(tasks ...Task) []interface{} {

	t.init(tasks)

	go t.produce(tasks)

	go t.consumeTasks()

	go t.progressBar()

	t.wait()

	return t.results

}

func (t *taskExecutor) wait() {
	<-t.fin
}

func (t *taskExecutor) done() {

	t.d <- true
}

func (t *taskExecutor) produce(tasks []Task) {

	defer close(t.taskChan)

	for idx, x := range tasks {
		memoizableT := NewMemoizableFutureTask(x, &t.cache)
		t.taskChan <- taskChanItem{idx, memoizableT}
	}
}

func (t *taskExecutor) progressBar() {
	ticker := time.NewTicker(1 * time.Second)

	taskLen := len(t.results)

	for {
		select {
		case <-ticker.C:
			progress := t.operationCount * 100 / taskLen
			t.ProgressC <- progress

		case <-t.d:
			t.ProgressC <- 100
			ticker.Stop()
			t.fin <- true
			return

		}
	}
}

func (t *taskExecutor) operationComplete() {

	taskLen := len(t.results)

	t.locker.Lock()

	t.operationCount++

	if t.operationCount == taskLen {
		t.done()
	}

	t.locker.Unlock()

}

func (t *taskExecutor) consumeTasks() {

	for x := range t.taskChan {
		t.threadLimit <- struct{}{}

		go func(i taskChanItem) {

			t.results[i.idx] = i.task.Exec()

			t.operationComplete()

			<-t.threadLimit

		}(x)
	}

}

// Progress accepts a function and passes current progress
// to that function.
// Triggered every second.
func (t *taskExecutor) Progress(progressF func(x int)) {
	go func() {

		isClosed := false
		l := sync.Mutex{}

		for x := range t.ProgressC {
			l.Lock()
			if !isClosed {
				progressF(x)
				if x == 100 {
					isClosed = true
					close(t.ProgressC)
				}
			}
			l.Unlock()

		}
	}()

}
