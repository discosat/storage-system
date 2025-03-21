package dam

import (
	"fmt"
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
		switch v := val.(type) {
		case string:
			stringQuery += key + ":" + v + ", "
		case float64:
			stringQuery += key + ":" + fmt.Sprintf("%f", v) + ", "
		case int64:
			stringQuery += key + ":" + fmt.Sprintf("%d", v) + ", "
		default:
			log.Fatal("Unsupported type in query", v)
		}
	}

	return qp.optimizer.Optimize(stringQuery)
}
