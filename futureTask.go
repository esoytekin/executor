package executor

// FutureTask godoc
type FutureTask[V any] struct {
	result V
	err    error
	signal chan struct{}
}

// Get godoc
func (f *FutureTask[V]) Get() (V, error) {

	<-f.signal
	return f.result, f.err

}

// NewFutureTask godoc
func NewFutureTask[V any](t Task[V]) *FutureTask[V] {
	f := new(FutureTask[V])

	f.signal = make(chan struct{})

	go func() {
		defer close(f.signal)
		f.result, f.err = t.Exec()
	}()

	return f
}
