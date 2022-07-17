package executor

import (
	"hash/fnv"
	"sync"
)

var mutex sync.Mutex

type MemoizableFutureTask[V any] struct {
	taskItem Task[V]
	cache    *sync.Map
}

func NewMemoizableFutureTask[V any](t Task[V], cache *sync.Map) Task[V] {
	return &MemoizableFutureTask[V]{
		taskItem: t,
		cache:    cache,
	}
}

// Exec godoc
func (m *MemoizableFutureTask[V]) Exec() V {

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

	return result.(*FutureTask[V]).Get()
}

// Hash godoc
// is used for memoization
// generate it using input parameters
func (m *MemoizableFutureTask[V]) Hash(hashStr func(string) int) []int {
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
