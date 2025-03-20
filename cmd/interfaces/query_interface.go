package interfaces

type QueryOptimizer interface {
	Optimize(query string) (string, error)
}
