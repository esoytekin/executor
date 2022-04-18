package executor

// FutureTask godoc
type FutureTask struct {
	result interface{}
	signal chan struct{}
}

// Get godoc
func (f *FutureTask) Get() interface{} {

	<-f.signal
	return f.result

}

// NewFutureTask godoc
func NewFutureTask(t Task) *FutureTask {
	f := new(FutureTask)

	f.signal = make(chan struct{})

	go func() {
		defer close(f.signal)
		f.result = t.Exec()
	}()

	return f
}
