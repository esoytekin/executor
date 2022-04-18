package executor

import (
	"hash/fnv"
	"sync"
)

var mutex sync.Mutex

type MemoizableFutureTask struct {
	taskItem Task
	cache    *sync.Map
}

func NewMemoizableFutureTask(t Task, cache *sync.Map) Task {
	return &MemoizableFutureTask{
		taskItem: t,
		cache:    cache,
	}
}

// Exec godoc
func (m *MemoizableFutureTask) Exec() interface{} {

	hashVal := getHash(m.Hash(hashStr))

	result, ok := m.cache.Load(hashVal)

	if !ok {

		mutex.Lock()
		result, ok = m.cache.Load(hashVal)

		if !ok {
			ft := NewFutureTask(m.taskItem)
			m.cache.Store(hashVal, ft)
			result = ft
		}

		mutex.Unlock()

	}

	return result.(*FutureTask).Get()
}

// Hash godoc
// is used for memoization
// generate it using input parameters
func (m *MemoizableFutureTask) Hash(hashStr func(string) int) []int {
	return m.taskItem.Hash(hashStr)
}

func hashStr(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func getHash(is []int) int {
	ourhashes := 5

	for _, h := range is {
		ourhashes = 21*ourhashes + h
	}
	return ourhashes
}
