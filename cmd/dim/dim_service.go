package dim

import (
	"archive/zip"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func handleUploadImage(c *gin.Context) {

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
	if err := row.Scan(&request.Id, &request.Name, &request.UserId, &request.MissionId, &mission.Id, &mission.Name, &mission.Bucket); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorAbortMessage(c, http.StatusNotFound, err)
			return
		}
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	// Checking if the observation that relates to the request has already been saved, or creates one if not
	var observation Observation
	row = db.QueryRow("SELECT * FROM observation WHERE request_id = $1", request.Id)
	err = row.Scan(&observation.Id, &observation.RequestId, &observation.UserId)

	var qErr error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			qErr = db.QueryRow("INSERT INTO observation(request_id, user_id) VALUES ($1, $2) RETURNING id, request_id, user_id", request.Id, request.UserId).
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
	key, err := ObjectStore.ds.SaveImage(file, oFile, mission.Bucket, request.Name)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}
	err = db.QueryRow("INSERT INTO measurement(object_reference, observation_id, measurement_request_id) VALUES ($1, $2, $3) RETURNING id", key, observation.Id, measurementRequest.Id).Scan(&measurementId)
	log.Printf("Object key: %v", key)
	c.JSON(http.StatusOK, gin.H{"measurement": measurementId})
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
		//err = db.QueryRow("INSERT INTO measurements (ref) VALUES ($1) RETURNING id", ref).Scan(&measurementId)
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
