package executor

import (
	"sync"
	"time"
)

const defaultThreadLimitCount = 6

type taskChanItem struct {
	idx  int
	task Task
}

// TaskExecutor godoc
type TaskExecutor struct {
	taskChan       chan taskChanItem
	d              chan bool
	fin            chan bool
	taskLen        int
	results        []interface{}
	threadLimit    chan struct{}
	locker         sync.Mutex
	operationCount int
	ProgressC      chan int
	cache          sync.Map
}

func NewSingleThreadTaskExecutor() *TaskExecutor {
	return NewTaskExecutor(1)
}

func NewTaskExecutor(threadLimit int) *TaskExecutor {
	taskChan := make(chan taskChanItem, 0)
	done := make(chan bool)
	m := make(chan struct{}, threadLimit)
	progressC := make(chan int, 100)
	return &TaskExecutor{
		taskChan:    taskChan,
		d:           done,
		fin:         make(chan bool),
		threadLimit: m,
		locker:      sync.Mutex{},
		ProgressC:   progressC,
	}
}

func (t *TaskExecutor) init(tasks []Task) {
	t.taskLen = len(tasks)
	t.results = make([]interface{}, t.taskLen)

}

func (t *TaskExecutor) ExecuteTask(tasks ...Task) []interface{} {

	t.init(tasks)

	go t.produce(tasks)

	go t.consumeTasks()

	go t.progressBar()

	t.wait()

	return t.results

}

func (t *TaskExecutor) wait() {
	<-t.fin
}

func (t *TaskExecutor) done() {

	t.d <- true
}

func (t *TaskExecutor) produce(tasks []Task) {

	defer close(t.taskChan)

	for idx, x := range tasks {
		memoizableT := NewMemoizableFutureTask(x, &t.cache)
		t.taskChan <- taskChanItem{idx, memoizableT}
	}
}

func (t *TaskExecutor) progressBar() {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			progress := t.operationCount * 100 / t.taskLen
			t.ProgressC <- progress

		case <-t.d:
			t.ProgressC <- 100
			ticker.Stop()
			t.fin <- true
			return

		}
	}
}

func (t *TaskExecutor) operationComplete() {

	t.locker.Lock()

	t.operationCount++

	if t.operationCount == t.taskLen {
		t.done()
	}

	t.locker.Unlock()

}

func (t *TaskExecutor) consumeTasks() {

	for x := range t.taskChan {
		t.threadLimit <- struct{}{}

		go func(i taskChanItem) {

			t.results[i.idx] = i.task.Exec()

			t.operationComplete()

			<-t.threadLimit

		}(x)
	}

}

func (t *TaskExecutor) Progress(progressF func(x int)) {
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
