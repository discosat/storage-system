package qom

import "log"

type DiscoQO struct{}

func (q *DiscoQO) Optimize(query string) error {
	log.Println("Logging Query in QOM: ", query)
	return nil
}
