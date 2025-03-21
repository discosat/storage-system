package qom

import (
	"github.com/discosat/storage-system/cmd/interfaces"
	"log"
)

type DiscoQO struct{}

func (q *DiscoQO) Optimize(query string) error {
	log.Println("Logging Query in QOM: ", query)
	return nil
}

var _ interfaces.QueryOptimizer = (*DiscoQO)(nil)
