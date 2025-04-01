package disco_qom

import (
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
	"strings"
)

type DiscoQO struct{}

func (q *DiscoQO) Optimize(query interfaces.ImageRequest) error {
	sqlQuery, args := StringToSQLTranslator(query)

	fmt.Println("Logging optimized query: ", sqlQuery)
	fmt.Println("Logging optimized query arguments: ", args)

	return nil
}

func StringToSQLTranslator(query interfaces.ImageRequest) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	//Creates an SQL query based on the query string
	if query.ImgID != nil {
		fmt.Println("Logging image ID after stripping it from image request", *query.ImgID)
		conditions = append(conditions, fmt.Sprintf("img_id = ?", argIndex))
		args = append(args, *query.ImgID)
		argIndex++
	}

	if query.ObsID != nil {
		fmt.Println("Logging observation ID after stripping it from image request", *query.ObsID)
		conditions = append(conditions, fmt.Sprintf("obs_id = ?", argIndex))
		args = append(args, *query.ObsID)
		argIndex++
	}

	if query.StartTime != nil {
		fmt.Println("Logging start time after stripping it from image request", *query.StartTime)
		conditions = append(conditions, fmt.Sprintf("time >= ?", argIndex))
		args = append(args, *query.StartTime)
		argIndex++
	}

	if query.EndTime != nil {
		fmt.Println("Logging end time after stripping it from image request", *query.EndTime)
		conditions = append(conditions, fmt.Sprintf("time <= ?", argIndex))
		args = append(args, *query.EndTime)
		argIndex++
	}

	if query.LatFrom != nil {
		fmt.Println("Logging latitude from after stripping it from image request", *query.LatFrom)
		conditions = append(conditions, fmt.Sprintf("latitude >= ?", argIndex))
		args = append(args, *query.LatFrom)
		argIndex++
	}

	if query.LatTo != nil {
		fmt.Println("Logging latitude to after stripping it from image request", *query.LatTo)
		conditions = append(conditions, fmt.Sprintf("latitude <= ?", argIndex))
		args = append(args, *query.LatTo)
		argIndex++
	}

	if query.LonFrom != nil {
		fmt.Println("Logging longitude from after stripping it from image request", *query.LonFrom)
		conditions = append(conditions, fmt.Sprintf("longitude >= ?", argIndex))
		args = append(args, *query.LonFrom)
		argIndex++
	}

	if query.LonTo != nil {
		fmt.Println("Logging longitude to after stripping it from image request", *query.LonTo)
		conditions = append(conditions, fmt.Sprintf("longitude <= ?", argIndex))
		args = append(args, *query.LonTo)
		argIndex++
	}

	if query.CamType != nil {
		fmt.Println("Logging cam type after stripping it from image request", *query.CamType)
		conditions = append(conditions, fmt.Sprintf("cam_type = ?", argIndex))
		args = append(args, *query.CamType)
		argIndex++
	}

	if query.Date != nil {
		fmt.Println("Logging date after stripping it from image request", *query.Date)
		conditions = append(conditions, fmt.Sprintf("date = ?", argIndex))
		args = append(args, *query.Date)
		argIndex++
	}

	sqlQuery := "SELECT * FROM images"
	if len(conditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	return sqlQuery, args
}
