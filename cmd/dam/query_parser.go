package dam

import (
	"github.com/discosat/storage-system/cmd/interfaces"
)

type QueryParser struct {
	optimizer interfaces.QueryOptimizer
}

func NewQueryParser(optimizer interfaces.QueryOptimizer) *QueryParser {
	return &QueryParser{optimizer: optimizer}
}
func (qp *QueryParser) ParseQuery(query string) error {
	return qp.optimizer.Optimize(query)
}
