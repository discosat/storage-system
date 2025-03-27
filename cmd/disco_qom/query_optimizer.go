package disco_qom

import (
	"github.com/discosat/storage-system/cmd/interfaces"
	"log"
)

type DiscoQO struct{}

func (q *DiscoQO) Optimize(query interfaces.ImageRequest) error {

	return nil
}

func StringToSQLTranslator(queryArray []string) {
	for v := range queryArray {
		log.Println("Query String: ", queryArray[v])
	}
}

//var _ interfaces.QueryOptimizer = (*DiscoQO)(nil)
