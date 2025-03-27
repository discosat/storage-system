package dam

import (
	"github.com/discosat/storage-system/cmd/interfaces"
	"log"
)

type QueryParser struct {
	qom interfaces.QueryOptimizer
}

/*func NewQueryParser(optimizer interfaces.QueryOptimizer) *QueryParser {
	return &QueryParser{optimizer: optimizer}
}

func (qp *QueryParser) PushQuery(query map[string]interface{}) error {
	stringQuery := ""

	log.Println("DAM query logged in query_parser: ", query)

	for key, val := range query {

		if val == nil {
			continue
		}

		switch v := val.(type) {
		case string:
			stringQuery += key + ":" + v + ","
		case float64:
			stringQuery += key + ":" + fmt.Sprintf("%f", v) + ","
		case int64:
			stringQuery += key + ":" + fmt.Sprintf("%d", v) + ","
		case *string:
			if v != nil {
				stringQuery += key + ":" + fmt.Sprintf(*v) + ","
			}
		case *float64:
			if v != nil {
				stringQuery += key + ":" + fmt.Sprintf("%f", *v) + ","
			}
		case *int64:
			if v != nil {
				stringQuery += key + ":" + fmt.Sprintf("%d", *v) + ","
			}
		default:
			log.Fatal("Unsupported type in query", v)
		}
	}

	if len(stringQuery) > 0 {
		stringQuery = stringQuery[:len(stringQuery)-2]
	}

	return qp.optimizer.Optimize(stringQuery)
}*/

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
