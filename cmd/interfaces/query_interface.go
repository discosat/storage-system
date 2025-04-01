package interfaces

/*QueryOptimizer interface which defines a contract that wants to implement the Query Optimizer module.
It (so far) contains a single method "Optimize" which takes a query string as input, and returns an error if one occurs*/

type QueryOptimizer interface {
	Optimize(query ImageRequest) error
}
