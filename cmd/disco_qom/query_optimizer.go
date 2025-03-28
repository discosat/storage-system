package disco_qom

import (
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
)

type DiscoQO struct{}

func (q *DiscoQO) Optimize(query interfaces.ImageRequest) error {
	//sqlQuery, args := StringToSQLTranslator(query)
	StringToSQLTranslator(query)

	return nil
}

func StringToSQLTranslator(query interfaces.ImageRequest) (string, []interface{}) {
	/*var conditions []string
	var args []interface{}
	argIndex := 1*/

	if query.ImgID != nil {
		fmt.Println("Logging image ID after stripping it from image request", *query.ImgID)
	}
	return "", nil
}

//var _ interfaces.QueryOptimizer = (*DiscoQO)(nil)
