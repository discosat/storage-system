package disco_qom

import (
	"github.com/discosat/storage-system/cmd/interfaces"
	"log"
	"strings"
)

type DiscoQO struct{}

func (q *DiscoQO) Optimize(query string) error {
	log.Println("Logging Query in QOM: ", query)

	splitString := strings.Split(query, ",")

	log.Println("Split string in QOM: ", splitString)

	StringToSQLTranslator(splitString)

	return nil
}

func StringToSQLTranslator(queryArray []string) {
	for i, v := range queryArray {
		log.Println(i, ": ", v)
	}
}

var _ interfaces.QueryOptimizer = (*DiscoQO)(nil)
