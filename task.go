package executor

// Task godoc
type Task[V any] interface {
	// Exec godoc
	Exec() (V, error)

	// Hash godoc
	// is used for memoization
	// generate it using input parameters
	Hash(hashStr func(string) int) []int
}
