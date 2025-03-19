package dim

import (
	"archive/zip"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var ObjectStore SimpleStore
var db *sql.DB

func Start() {
	ObjectStore = NewSimpleStore(NewMinioStore())
	db = ConfigDatabase()
	defer db.Close()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/batch", UploadBatch)
	router.POST("/file", UploadImage)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func ConfigDatabase() *sql.DB {
	db, err := sql.Open("pgx", fmt.Sprint("postgres://", os.Getenv("PGUSER"), ":", os.Getenv("PGPASSWORD"), "@", os.Getenv("PGHOST"), "/", os.Getenv("PGDATABASE")))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	return db
}

func UploadBatch(c *gin.Context) {

	bucketName := c.Request.FormValue("bucketName")

	exists, err := ObjectStore.ds.BucketExists(bucketName)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	if !exists {
		e := fmt.Errorf("no bucket with name '%v' exists. Specify an already existing bucket", bucketName)
		ErrorAbortMessage(c, http.StatusBadRequest, e)
		return
	}

	file, _, err := c.Request.FormFile("batch")
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	tmpFile, _ := os.CreateTemp("", "temp*.zip")
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	archive, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	//var success string
	for _, iFile := range archive.File {
		if iFile.FileInfo().IsDir() {
			continue
		}
		oFile, _ := iFile.Open()
		_, err := ObjectStore.ds.SaveFile(iFile, oFile, bucketName)
		if err != nil {
			log.Printf("UploadBatch: Cannot upload thie file %v, error is %v", filepath.Base(iFile.Name), err)
			break
		}
		//var measurementId string
		//err = db.QueryRow("INSERT INTO measurements (ref) VALUES ($1) RETURNING id", ref).Scan(&measurementId)
		//log.Printf("MEASUREMENT ID: %v", measurementId)
	}

	log.Println("done")

}

type Mission struct {
	id     int
	name   string
	bucket string
}

type Request struct {
	id   int
	name string
	mId  int
	uId  int
}

type MeasurementRequest struct {
	Id    int    `json:"id"`
	MType string `json:"m_type"`
	RId   int    `json:"r_id"`
}

type Observation struct {
	Id        int
	RequestId int
	UserId    int
}

func UploadImage(c *gin.Context) {

	//Binding POST data
	requestId := c.Request.FormValue("requestId")
	log.Printf("Querying for request with id %v", requestId)
	file, err := c.FormFile("file")
	if err != nil {
		ErrorAbortMessage(c, http.StatusBadRequest, err)
		return
	}
	log.Println(file.Filename)

	// Getting request and mission data
	var mission Mission
	var request Request
	row := db.QueryRow("SELECT * FROM request r INNER JOIN mission m on m.id = r.mission_id where r.id = $1", requestId)
	if err := row.Scan(&request.id, &request.name, &request.uId, &request.mId, &mission.id, &mission.name, &mission.bucket); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorAbortMessage(c, http.StatusNotFound, err)
			return
		}
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	// Checking if the observation that relates to the request has already been saved, or creates one if not
	var observation Observation
	row = db.QueryRow("SELECT * FROM observation WHERE request_id = $1", request.id)
	err = row.Scan(&observation.Id, &observation.RequestId, &observation.UserId)

	var qErr error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			qErr = db.QueryRow("INSERT INTO observation(request_id, user_id) VALUES ($1, $2) RETURNING id, request_id, user_id", request.id, request.uId).
				Scan(&observation.Id, &observation.RequestId, &observation.UserId)
		} else {
			ErrorAbortMessage(c, http.StatusInternalServerError, err)
			return
		}
	}

	if qErr != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	//var measurementRequests []MeasurementRequest
	//rows, err := db.Query("SELECT * FROM measurement_request WHERE request_id = $1", request.id)
	//if err != nil {
	//	ErrorAbortMessage(c, http.StatusNotFound, err)
	//	return
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var mr MeasurementRequest
	//	if err := rows.Scan(&mr.Id, &mr.RId, &mr.MType); err != nil {
	//		ErrorAbortMessage(c, http.StatusInternalServerError, err)
	//		return
	//	}
	//	measurementRequests = append(measurementRequests, mr)
	//}
	//if err = rows.Err(); err != nil {
	//	ErrorAbortMessage(c, http.StatusInternalServerError, err)
	//}

	// Gets the measurement request to relate the measurement to
	var measurementRequest MeasurementRequest
	row = db.QueryRow("SELECT * FROM measurement_request WHERE request_id = $1 AND type = $2", requestId, "Narrow-Image")
	if err := row.Scan(&measurementRequest.Id, &measurementRequest.RId, &measurementRequest.MType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorAbortMessage(c, http.StatusNotFound, err)
			return
		}
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	// Inserts measurement into db/object store
	var measurementId int

	oFile, err := file.Open()
	if err != nil {
		ErrorAbortMessage(c, http.StatusBadRequest, err)
		return
	}
	key, err := ObjectStore.ds.SaveImage(file, oFile, mission.bucket, request.name)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}
	err = db.QueryRow("INSERT INTO measurement(object_reference, observation_id, measurement_request_id) VALUES ($1, $2, $3) RETURNING id", key, observation.Id, measurementRequest.Id).Scan(&measurementId)

	log.Printf("Object key: %v", key)

	c.JSON(http.StatusOK, gin.H{"measurement": measurementId})

	return

}

func ErrorAbortMessage(c *gin.Context, statusCode int, err error) {
	log.Println(err)
	c.JSON(statusCode, gin.H{"error": fmt.Sprint(err)})
}
