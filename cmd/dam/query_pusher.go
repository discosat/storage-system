package dam

import (
	"github.com/discosat/storage-system/cmd/interfaces"
	"log"
)

type QueryParser struct {
	qom interfaces.QueryOptimizer
}

func newQueryPusher(qom interfaces.QueryOptimizer) *QueryParser {
	return &QueryParser{qom: qom}
}

func (p *QueryParser) PushQuery(query interfaces.ImageRequest) error {
	err := p.qom.Optimize(query)

	if err != nil {
		log.Println("Error passing on query", err)
		return err
	}

	return nil
}
