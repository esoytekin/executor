package executor

// FutureTask godoc
type FutureTask[V any] struct {
	result V
	signal chan struct{}
}

// Get godoc
func (f *FutureTask[V]) Get() V {

	<-f.signal
	return f.result

}

// NewFutureTask godoc
func NewFutureTask[V any](t Task[V]) *FutureTask[V] {
	f := new(FutureTask[V])

	f.signal = make(chan struct{})

	go func() {
		defer close(f.signal)
		f.result = t.Exec()
	}()

	return f
}
