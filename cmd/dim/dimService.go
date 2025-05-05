package dim

import (
	//"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	//"fmt"
	. "github.com/discosat/storage-system/internal/observation"
	. "github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"log"
	"mime/multipart"
	//"os"
	"os/exec"
	//"path/filepath"
)

type DimServiceInterface interface {
	handleUploadImage(c *gin.Context)
	test(c *gin.Context) (ObservationRequestAggregate, error)
}

type DimService struct {
	observationRequestRepository ObservationRequestRepository
	observationRepository        ObservationRepository
}

func NewDimService(oRepo ObservationRepository, orRepo ObservationRequestRepository) *DimService {
	return &DimService{
		observationRequestRepository: orRepo,
		observationRepository:        oRepo,
	}
}

func (d DimService) test(c *gin.Context) (ObservationRequestAggregate, error) {
	var qId = c.Query("orId")
	//orId, err := strconv.ParseInt(qId, 10, 0)
	orId, err := strconv.Atoi(qId)
	if err != nil {
		log.Fatalf("Not an int")
	}
	log.Println(orId)
	or, err := d.observationRequestRepository.GetObservationRequest(orId)
	if err != nil {
		log.Fatalf("Get observation request went wrong: %v", err)
	}
	log.Println(or)
	return or, err
}

func (d DimService) handleUploadImage(c *gin.Context) {
	//Binding POST data
	file, err := c.FormFile("file")
	if err != nil {
		ErrorAbortMessage(c, http.StatusBadRequest, err)
		return
	}

	//// Getting flightPlan and mission data
	//log.Printf("Querying for FlightPlan with id %v", flightPlanId)
	//flightPlan, err := d.flightPlanRepository.GetById(flightPlanId)
	//if err != nil {
	//	ErrorAbortMessage(c, http.StatusNotFound, err)
	//	return
	//}
	//
	//log.Printf("Querying for Mission with id %v", flightPlan.MissionId)
	//mission, err := d.missionRepository.GetById(flightPlan.MissionId)
	//if err != nil {
	//	ErrorAbortMessage(c, http.StatusInternalServerError, err)
	//	return
	//}

	// TODO Do error handling
	// Gets related measurment flightPlan
	ObservationRequestId := extractMetadata(file)

	log.Printf("Querying for ObservationRequest with id %v", ObservationRequestId)
	observationRequestAggr, err := d.observationRequestRepository.GetObservationRequest(ObservationRequestId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorAbortMessage(c, http.StatusNotFound, err)
			return
		}
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	// Saves image
	observationId, err := d.observationRepository.CreateObservation(file, observationRequestAggr.Mission.Bucket, observationRequestAggr.FlightPlan.Name, observationRequestAggr.ObservationRequest.Id)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		//log.Fatalf("handleImageUpload: %v", err)
		return
	}
	//
	//_, err = d.observationMetadataRepository.CreateObservationMetadata(observationId, 10.4058633, 55.3821913)
	//if err != nil {
	//	log.Fatalf("Geom: %v", err)
	//}

	// TODO Genovervej lige sletning ved fejl
	//err = tx.Commit()
	//if err != nil {
	//	ObjectStore.ds.DeleteImage(key, mission.Bucket)
	//	ErrorAbortMessage(c, http.StatusBadRequest, err)
	//	return
	//}

	c.JSON(http.StatusCreated, gin.H{"observation": observationId})
	return
}

//func handleUploadBatch(c *gin.Context) {
//
//	bucketName := c.FlightPlan.FormValue("bucketName")
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
//	file, _, err := c.FlightPlan.FormFile("batch")
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
//func handleGetRequests(c *gin.Context) ([]FlightPlan, error) {
//
//	missionId := c.Query("missionId")
//
//	var requests []FlightPlan
//	rows, err := db.Query("SELECT * FROM flightPlan WHERE mission_id = $1", missionId)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		var flightPlan FlightPlan
//		if err := rows.Scan(&flightPlan.Id, &flightPlan.Name, &flightPlan.UserId, &flightPlan.MissionId); err != nil {
//			return requests, err
//		}
//		requests = append(requests, flightPlan)
//	}
//
//	if err = rows.Err(); err != nil {
//		return requests, err
//	}
//
//	return requests, nil
//}
//
//func handleGetRequestsNoObservation(c *gin.Context) ([]FlightPlan, error) {
//	var requests []FlightPlan
//	rows, err := db.Query("SELECT r.id, r.name, r.user_id, r.mission_id FROM flightPlan r LEFT JOIN public.observation o on r.id = o.request_id WHERE o.id IS NULL")
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		var flightPlan FlightPlan
//		if err := rows.Scan(&flightPlan.Id, &flightPlan.Name, &flightPlan.UserId, &flightPlan.MissionId); err != nil {
//			return requests, err
//		}
//		requests = append(requests, flightPlan)
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
	relatedObservationRequest := l["measurementRequest"]
	return relatedObservationRequest
}
