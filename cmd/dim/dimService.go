package dim

import (
	"archive/zip"
	"bytes"
	. "github.com/discosat/storage-system/internal/Commands"

	"encoding/json"
	"path/filepath"
	"strconv"

	//"fmt"
	. "github.com/discosat/storage-system/internal/observation"
	. "github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"log"
	//"os"
	"os/exec"
	//"path/filepath"
)

type DimServiceInterface interface {
	handleUploadImage(file *io.ReadCloser, fileName string, fileSize int64) (int, error)
	handleUploadBatch(archive *zip.ReadCloser) error
	handleGetFlightPlan(id int) (FlightPlan, error)
	handleCreateFlightPlan(flightPlan FlightPlanCommand, requestList []ObservationRequestCommand) (int, error)
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

func (d DimService) handleUploadImage(file *io.ReadCloser, fileName string, fileSize int64) (int, error) {

	// TODO Do error handling

	// Need to open again after extracting metadata
	raw, _ := io.ReadAll(*file)
	ObservationRequestId := extractMetadata(raw)

	log.Printf("Querying for ObservationRequest with id %v", ObservationRequestId)
	observationRequestAggr, err := d.observationRequestRepository.GetObservationRequest(ObservationRequestId)
	if err != nil {
		return -1, err
	}

	fileReader := bytes.NewReader(raw)
	observationCommand := ObservationCommand{File: fileReader, FileName: fileName, FileSize: fileSize, Bucket: observationRequestAggr.Mission.Bucket, FlightPlanName: observationRequestAggr.FlightPlan.Name, ObservationRequestId: observationRequestAggr.ObservationRequest.Id}
	// Saves image
	observationId, err := d.observationRepository.CreateObservation(observationCommand)
	if err != nil {
		return -1, err
	}

	return observationId, nil
}

func (d DimService) handleGetFlightPlan(id int) (FlightPlan, error) {
	return d.observationRequestRepository.GetFlightPlantById(id)
}

func (d DimService) handleCreateFlightPlan(flightPlan FlightPlanCommand, requestList []ObservationRequestCommand) (int, error) {
	return d.observationRequestRepository.CreateFlightPlan(flightPlan, requestList)
}

func (d DimService) handleUploadBatch(archive *zip.ReadCloser) error {
	for _, iFile := range archive.File {
		if iFile.FileInfo().IsDir() {
			continue
		}
		oFile, err := iFile.Open()
		if err != nil {
			return err
		}
		//oFile, err = iFile.Open()
		_, err = d.handleUploadImage(&oFile, iFile.FileInfo().Name(), iFile.FileInfo().Size())
		if err != nil {
			log.Printf("UploadBatch: Cannot upload thie file %v, error is %v", filepath.Base(iFile.Name), err)
			break
		}
		err = oFile.Close()
		if err != nil {
			return err
		}
		//log.Println(result)
	}

	log.Println("done")

	return nil
}

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

func extractMetadata(raw []byte) int {

	// call exifTool
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
