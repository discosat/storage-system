package dam

import (
	"github.com/discosat/storage-system/cmd/interfaces"
	"log"
)

type QueryParser struct {
	optimizer interfaces.QueryOptimizer
}

func NewQueryParser(optimizer interfaces.QueryOptimizer) *QueryParser {
	return &QueryParser{optimizer: optimizer}
}

func (qp *QueryParser) ParseQuery(query map[string]interface{}) error {
	stringQuery := ""

	log.Println("DAM query logged in query_parser: ", query)

	for key, val := range query {
		stringQuery += key + ":" + val.(string) + ", "
	}

	return qp.optimizer.Optimize(stringQuery)
}
