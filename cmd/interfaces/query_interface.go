package interfaces

/*QueryOptimizer interface which defines a contract that wants to implement the Query Optimizer module.
It (so far) contains a single method "Translate" which takes a query string as input, and returns an error if one occurs*/

type QueryOptimizer interface {
	Translate(query ImageRequest) (string, []interface{}, error)
}
