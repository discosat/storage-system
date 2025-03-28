package disco_qom

import (
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
)

type DiscoQO struct{}

func (q *DiscoQO) Optimize(query interfaces.ImageRequest) error {

	return nil
}

func StringToSQLTranslator(queryArray []string) {
	for v := range queryArray {
		fmt.Println("Query String: ", queryArray[v])
	}
}

//var _ interfaces.QueryOptimizer = (*DiscoQO)(nil)
