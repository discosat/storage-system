package disco_qom

import (
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
	"strings"
)

type DiscoQO struct{}

func (q *DiscoQO) Translate(query interfaces.ImageRequest) (string, []interface{}, error) {
	sqlQuery, args := StringToSQLTranslator(query)

	return sqlQuery, args, nil
}

func StringToSQLTranslator(query interfaces.ImageRequest) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	//Creates an SQL query based on the query string
	if query.ImgID != nil {
		fmt.Println("Logging image ID after stripping it from image request", *query.ImgID)
		conditions = append(conditions, fmt.Sprintf("img_id = $%d", argIndex))
		args = append(args, *query.ImgID)
		argIndex++
	}

	if query.ObsID != nil {
		fmt.Println("Logging observation ID after stripping it from image request", *query.ObsID)
		conditions = append(conditions, fmt.Sprintf("obs_id = $%d", argIndex))
		args = append(args, *query.ObsID)
		argIndex++
	}

	if query.StartTime != nil {
		fmt.Println("Logging start time after stripping it from image request", *query.StartTime)
		conditions = append(conditions, fmt.Sprintf("time >= $%d", argIndex))
		args = append(args, *query.StartTime)
		argIndex++
	}

	if query.EndTime != nil {
		fmt.Println("Logging end time after stripping it from image request", *query.EndTime)
		conditions = append(conditions, fmt.Sprintf("time <= $%d", argIndex))
		args = append(args, *query.EndTime)
		argIndex++
	}

	if query.LatFrom != nil {
		fmt.Println("Logging latitude from after stripping it from image request", *query.LatFrom)
		conditions = append(conditions, fmt.Sprintf("latitude >= $%d", argIndex))
		args = append(args, *query.LatFrom)
		argIndex++
	}

	if query.LatTo != nil {
		fmt.Println("Logging latitude to after stripping it from image request", *query.LatTo)
		conditions = append(conditions, fmt.Sprintf("latitude <= $%d", argIndex))
		args = append(args, *query.LatTo)
		argIndex++
	}

	if query.LonFrom != nil {
		fmt.Println("Logging longitude from after stripping it from image request", *query.LonFrom)
		conditions = append(conditions, fmt.Sprintf("longitude >= $%d", argIndex))
		args = append(args, *query.LonFrom)
		argIndex++
	}

	if query.LonTo != nil {
		fmt.Println("Logging longitude to after stripping it from image request", *query.LonTo)
		conditions = append(conditions, fmt.Sprintf("longitude <= $%d", argIndex))
		args = append(args, *query.LonTo)
		argIndex++
	}

	if query.Camera != nil {
		fmt.Println("Logging cam type after stripping it from image request", *query.Camera)
		conditions = append(conditions, fmt.Sprintf("camera = $%d", argIndex))
		args = append(args, *query.Camera)
		argIndex++
	}

	if query.Date != nil {
		fmt.Println("Logging date after stripping it from image request", *query.Date)
		conditions = append(conditions, fmt.Sprintf("date = $%d", argIndex))
		args = append(args, *query.Date)
		argIndex++
	}

	sqlQuery := `SELECT observation_metadata.*, observation.bucket_name, observation.object_reference
				 FROM observation_metadata
				 JOIN observation ON observation_metadata.observation_id = observation.id`

	if len(conditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	sqlQuery += " AND observation.bucket_name = 'demodata'"

	return sqlQuery, args
}
