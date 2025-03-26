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

	log.Println(splitString)

	return nil
}

var _ interfaces.QueryOptimizer = (*DiscoQO)(nil)
