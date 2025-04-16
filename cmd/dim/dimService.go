package dim

import (
	//"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	//"fmt"
	. "github.com/discosat/storage-system/internal/measurement"
	. "github.com/discosat/storage-system/internal/measurementMetadata"
	. "github.com/discosat/storage-system/internal/measurementRequest"
	. "github.com/discosat/storage-system/internal/mission"
	. "github.com/discosat/storage-system/internal/observation"
	. "github.com/discosat/storage-system/internal/request"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	//"os"
	"os/exec"
	//"path/filepath"
)

type DimServiceInterface interface {
	handleUploadImage(c *gin.Context)
}

type DimService struct {
	requestRepository             RequestRepository
	missionRepository             MissionRepository
	observationRepository         ObservationRepository
	measurementRequestRepository  MeasurementRequestRepository
	measurementRepository         MeasurementRepository
	measurementMetadataRepository MeasurementMetadataRepository
}

func NewDimService(reRepo RequestRepository, miRepo MissionRepository, oRepo ObservationRepository, mrRepo MeasurementRequestRepository, meRepo MeasurementRepository, mdRepo MeasurementMetadataRepository) *DimService {
	return &DimService{
		requestRepository:             reRepo,
		missionRepository:             miRepo,
		observationRepository:         oRepo,
		measurementRequestRepository:  mrRepo,
		measurementRepository:         meRepo,
		measurementMetadataRepository: mdRepo,
	}
}

func (d DimService) handleUploadImage(c *gin.Context) {
	//Binding POST data
	requestId := c.Request.FormValue("requestId")

	log.Printf("Querying for request with id %v", requestId)
	file, err := c.FormFile("file")
	if err != nil {
		ErrorAbortMessage(c, http.StatusBadRequest, err)
		return
	}

	// Getting request and mission data
	request, err := d.requestRepository.GetById(requestId)
	if err != nil {
		ErrorAbortMessage(c, http.StatusNotFound, err)
		return
	}

	mission, err := d.missionRepository.GetById(request.MissionId)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	// Checking if the observation that relates to the request has already been saved, or creates one if not
	observation, err := d.observationRepository.GetByRequest(request.Id)
	var qErr error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			observation, qErr = d.observationRepository.CreateObservation(request.Id, request.UserId)
		} else {
			ErrorAbortMessage(c, http.StatusInternalServerError, err)
			return
		}
	}

	if qErr != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	// TODO Do error handling
	// Gets related measurment request
	measurementRequestId := extractMetadata(file)
	// Opening file

	// Checks if measurement request exists
	measurementRequest, err := d.measurementRequestRepository.GetById(measurementRequestId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorAbortMessage(c, http.StatusNotFound, err)
			return
		}
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	// Saves image
	measurementId, err := d.measurementRepository.CreateMeasurement(file, mission.Bucket, request.Name, observation.Id, measurementRequest.Id)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		//log.Fatalf("handleImageUpload: %v", err)
		return
	}

	_, err = d.measurementMetadataRepository.CreateMeasurementMetadata(measurementId, 10.4058633, 55.3821913)
	if err != nil {
		log.Fatalf("Geom: %v", err)
	}

	// TODO Genovervej lige sletning ved fejl
	//err = tx.Commit()
	//if err != nil {
	//	ObjectStore.ds.DeleteImage(key, mission.Bucket)
	//	ErrorAbortMessage(c, http.StatusBadRequest, err)
	//	return
	//}

	c.JSON(http.StatusCreated, gin.H{"measurement": measurementId})
	return
}

//func handleUploadBatch(c *gin.Context) {
//
//	bucketName := c.Request.FormValue("bucketName")
//
//	exists, err := ObjectStore.ds.BucketExists(bucketName)
//	if err != nil {
//		ErrorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	if !exists {
//		e := fmt.Errorf("No bucket with name '%v' exists. Specify an already existing bucket", bucketName)
//		ErrorAbortMessage(c, http.StatusBadRequest, e)
//		return
//	}
//
//	file, _, err := c.Request.FormFile("batch")
//	if err != nil {
//		ErrorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	tmpFile, _ := os.CreateTemp("", "temp*.zip")
//	defer os.Remove(tmpFile.Name())
//
//	_, err = io.Copy(tmpFile, file)
//	if err != nil {
//		ErrorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	archive, err := zip.OpenReader(tmpFile.Name())
//	if err != nil {
//		ErrorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	for _, iFile := range archive.File {
//		if iFile.FileInfo().IsDir() {
//			continue
//		}
//		oFile, _ := iFile.Open()
//		_, err := ObjectStore.ds.SaveFile(iFile, oFile, bucketName)
//		if err != nil {
//			log.Printf("UploadBatch: Cannot upload thie file %v, error is %v", filepath.Base(iFile.Name), err)
//			break
//		}
//		//var measurementId string
//		//err = tx.QueryRow("INSERT INTO measurements (ref) VALUES ($1) RETURNING id", ref).Scan(&measurementId)
//		//log.Printf("MEASUREMENT ID: %v", measurementId)
//	}
//
//	log.Println("done")
//
//}
//
//func handleGetMissions(c *gin.Context) ([]Mission, error) {
//	var missions []Mission
//
//	rows, err := db.Query("SELECT * FROM mission")
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		var mission Mission
//		if err := rows.Scan(&mission.Id, &mission.Name, &mission.Bucket); err != nil {
//			return missions, err
//		}
//		missions = append(missions, mission)
//	}
//
//	if err = rows.Err(); err != nil {
//		return missions, err
//	}
//	return missions, nil
//}
//
//func handleGetRequests(c *gin.Context) ([]Request, error) {
//
//	missionId := c.Query("missionId")
//
//	var requests []Request
//	rows, err := db.Query("SELECT * FROM request WHERE mission_id = $1", missionId)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		var request Request
//		if err := rows.Scan(&request.Id, &request.Name, &request.UserId, &request.MissionId); err != nil {
//			return requests, err
//		}
//		requests = append(requests, request)
//	}
//
//	if err = rows.Err(); err != nil {
//		return requests, err
//	}
//
//	return requests, nil
//}
//
//func handleGetRequestsNoObservation(c *gin.Context) ([]Request, error) {
//	var requests []Request
//	rows, err := db.Query("SELECT r.id, r.name, r.user_id, r.mission_id FROM request r LEFT JOIN public.observation o on r.id = o.request_id WHERE o.id IS NULL")
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		var request Request
//		if err := rows.Scan(&request.Id, &request.Name, &request.UserId, &request.MissionId); err != nil {
//			return requests, err
//		}
//		requests = append(requests, request)
//	}
//
//	if err = rows.Err(); err != nil {
//		return requests, err
//	}
//
//	return requests, nil
//
//}

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
