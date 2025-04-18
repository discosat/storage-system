package dam

import (
	"github.com/discosat/storage-system/cmd/interfaces"
	"log"
)

type QueryPusher struct {
	qom interfaces.QueryOptimizer
}

func newQueryPusher(qom interfaces.QueryOptimizer) *QueryPusher {
	return &QueryPusher{qom: qom}
}

func (p *QueryPusher) PushQuery(query interfaces.ImageRequest) error {
	err := p.qom.Optimize(query)

	if err != nil {
		log.Println("Error passing on query", err)
		return err
	}

	return nil
}
