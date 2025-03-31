package disco_qom

import (
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
	"strings"
)

type DiscoQO struct{}

func (q *DiscoQO) Optimize(query interfaces.ImageRequest) error {
	//sqlQuery, args := StringToSQLTranslator(query)
	StringToSQLTranslator(query)

	return nil
}

func StringToSQLTranslator(query interfaces.ImageRequest) (string, []interface{}) {
	var conditions []string
	/*var args []interface{}
	argIndex := 1*/

	//Creates an SQL query based on the query string
	//Is under risk of SQL injection - Fix using placeholders
	if query.ImgID != nil {
		fmt.Println("Logging image ID after stripping it from image request", *query.ImgID)
		conditions = append(conditions, fmt.Sprintf("img_id = %s", *query.ImgID))
	}

	if query.CamType != nil {
		fmt.Println("Logging cam type after stripping it from image request", *query.CamType)
		conditions = append(conditions, fmt.Sprintf("cam_type = %s", *query.CamType))
	}

	sqlQuery := "SELECT * FROM images"
	if len(conditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	fmt.Println("Logging sqlQuery: ", sqlQuery)

	return "", nil
}

//var _ interfaces.QueryOptimizer = (*DiscoQO)(nil)
