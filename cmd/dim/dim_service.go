package dim

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/discosat/storage-system/internal/mission"
	. "github.com/discosat/storage-system/internal/request"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var requestRepo = PsqlRequestRepository{RequestRepository: PsqlRequestRepository{}}
var missionRepo = PsqlMissionRepository{MissionRepository: PsqlMissionRepository{}}

func handleUploadImage(c *gin.Context) {

	//Binding POST data
	requestId := c.Request.FormValue("requestId")

	log.Printf("Querying for request with id %v", requestId)
	file, err := c.FormFile("file")
	if err != nil {
		ErrorAbortMessage(c, http.StatusBadRequest, err)
		return
	}

	//Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Transaction: %v", err)
	}

	// Getting request and mission data
	request, err := requestRepo.GetById(db, requestId)
	if err != nil {
		ErrorAbortMessage(c, http.StatusNotFound, err)
		return
	}

	mission, err := missionRepo.GetById(db, request.MissionId)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	// Checking if the observation that relates to the request has already been saved, or creates one if not
	var observation Observation
	row := db.QueryRow("SELECT * FROM observation WHERE request_id = $1", request.Id)
	err = row.Scan(&observation.Id, &observation.RequestId, &observation.UserId)

	var qErr error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			qErr = tx.QueryRow("INSERT INTO observation(request_id, user_id) VALUES ($1, $2) RETURNING id, request_id, user_id", request.Id, request.UserId).
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
	// Gets related measurment request
	measurementRequest := extractMetadata(file)

	// Opening file
	oFile, err := file.Open()
	if err != nil {
		ErrorAbortMessage(c, http.StatusBadRequest, err)
		return
	}

	// Checks if measurement request exists
	row = db.QueryRow("SELECT * FROM measurement_request WHERE id = $1", measurementRequest)
	if err := row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorAbortMessage(c, http.StatusNotFound, err)
			return
		}
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	var measurementId int
	// Saves image in object store
	key, err := ObjectStore.ds.SaveImage(file, oFile, mission.Bucket, request.Name)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}
	// Saves reference to object in SQL DB
	err = tx.QueryRow("INSERT INTO measurement(object_reference, observation_id, measurement_request_id) VALUES ($1, $2, $3) RETURNING id", key, observation.Id, measurementRequest).Scan(&measurementId)
	if err != nil {
		log.Fatalf("handleImageUpload: %v", err)
	}

	var metaId int
	err = tx.QueryRow("INSERT INTO measurement_metadata(measurement_id, location) VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326)) RETURNING id", measurementId, 10.4058633, 55.3821913).Scan(&metaId)
	if err != nil {
		log.Fatalf("Geom: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		ObjectStore.ds.DeleteImage(key, mission.Bucket)
		ErrorAbortMessage(c, http.StatusBadRequest, err)
		return
	}

	log.Printf("Object key: %v", key)

	c.JSON(http.StatusCreated, gin.H{"measurement": measurementId})
	return
}

func handleUploadBatch(c *gin.Context) {

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
		//err = tx.QueryRow("INSERT INTO measurements (ref) VALUES ($1) RETURNING id", ref).Scan(&measurementId)
		//log.Printf("MEASUREMENT ID: %v", measurementId)
	}

	log.Println("done")

}

func handleGetMissions(c *gin.Context) ([]Mission, error) {
	var missions []Mission

	rows, err := db.Query("SELECT * FROM mission")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var mission Mission
		if err := rows.Scan(&mission.Id, &mission.Name, &mission.Bucket); err != nil {
			return missions, err
		}
		missions = append(missions, mission)
	}

	if err = rows.Err(); err != nil {
		return missions, err
	}
	return missions, nil
}

func handleGetRequests(c *gin.Context) ([]Request, error) {

	missionId := c.Query("missionId")

	var requests []Request
	rows, err := db.Query("SELECT * FROM request WHERE mission_id = $1", missionId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var request Request
		if err := rows.Scan(&request.Id, &request.Name, &request.UserId, &request.MissionId); err != nil {
			return requests, err
		}
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return requests, err
	}

	return requests, nil
}

func handleGetRequestsNoObservation(c *gin.Context) ([]Request, error) {
	var requests []Request
	rows, err := db.Query("SELECT r.id, r.name, r.user_id, r.mission_id FROM request r LEFT JOIN public.observation o on r.id = o.request_id WHERE o.id IS NULL")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var request Request
		if err := rows.Scan(&request.Id, &request.Name, &request.UserId, &request.MissionId); err != nil {
			return requests, err
		}
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return requests, err
	}

	return requests, nil

}

func extractMetadata(file *multipart.FileHeader) int {
	oFile, err := file.Open()
	if err != nil {
		log.Fatalf("extractMetadata: %v", err)
	}
	defer oFile.Close()

	// call exifTool
	raw, err := io.ReadAll(oFile)
	cmd := exec.Command("exiftool", "-json", "-")
	cmd.Stdin = bytes.NewReader(raw)
	result, err := cmd.Output()
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	// Unmarshal the EXIF data to a map of properties in the comment tag
	var t []map[string]string
	json.Unmarshal(result, &t)
	s := t[0]["Comment"]
	var l map[string]int
	err = json.Unmarshal([]byte(s), &l)
	if err != nil {
		log.Fatalf("err 3: %v", err)
	}
	relatedMeasurementRequest := l["measurementRequest"]
	return relatedMeasurementRequest
}
