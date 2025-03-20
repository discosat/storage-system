package qom

type DiscoQO struct{}

func (q *DiscoQO) Optimize(query string) (string, error) {
	return query, nil
}
