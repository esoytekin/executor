package executor

// Task godoc
type Task interface {
	// Exec godoc
	Exec() interface{}

	// Hash godoc
	// is used for memoization
	// generate it using input parameters
	Hash(hashStr func(string) int) []int
}
