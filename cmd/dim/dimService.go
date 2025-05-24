package dim

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/discosat/storage-system/internal/Commands"
	"log/slog"
	"strconv"

	. "github.com/discosat/storage-system/internal/observation"
	. "github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"log"
	"os/exec"
)

type DimServiceInterface interface {
	handleUploadImage(file *io.ReadCloser, fileName string, fileSize int64) (int, error)
	handleUploadBatch(archive *zip.ReadCloser) error
	handleGetFlightPlan(id int) (FlightPlanAggregate, error)
	handleCreateFlightPlan(flightPlan CreateFlightPlanCommand, requestList []CreateObservationRequestCommand) (int, error)
	handleUpdateFlightPlan(flightPlan FlightPlanAggregate) (int, error)
	handleDeleteFlightPlan(id int) (bool, error)
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

	// Reading bytes, to be able to use them twice
	raw, err := io.ReadAll(*file)
	if err != nil {
		slog.Error("Could not read file")
		return -1, err
	}

	log.Printf("Extracting metadata from observation")
	observationRequestId, metadata, err := extractMetadata(raw)
	if err != nil {
		return -1, err
	}

	log.Println(metadata)

	log.Printf("Querying for ObservationRequest with id %v", observationRequestId)
	observationRequestAggr, err := d.observationRequestRepository.GetObservationRequest(observationRequestId)
	if err != nil {
		return -1, err
	}

	fileReader := bytes.NewReader(raw)
	observationCommand := ObservationCommand{File: fileReader, FileName: fileName, FileSize: fileSize, Bucket: observationRequestAggr.Mission.Bucket, FlightPlanName: observationRequestAggr.FlightPlan.Name, ObservationRequestId: observationRequestAggr.ObservationRequest.Id}
	// Saves image
	observationId, err := d.observationRepository.CreateObservation(observationCommand, &metadata)
	if err != nil {
		return -1, err
	}

	return observationId, nil
}

func (d DimService) handleGetFlightPlan(id int) (FlightPlanAggregate, error) {
	return d.observationRequestRepository.GetFlightPlanById(id)
}

func (d DimService) handleCreateFlightPlan(flightPlan CreateFlightPlanCommand, requestList []CreateObservationRequestCommand) (int, error) {
	return d.observationRequestRepository.CreateFlightPlan(flightPlan, requestList)
}

func (d DimService) handleUpdateFlightPlan(flightPlan FlightPlanAggregate) (int, error) {
	return d.observationRequestRepository.UpdateFlightPlan(flightPlan)

}

func (d DimService) handleDeleteFlightPlan(id int) (bool, error) {
	return d.observationRequestRepository.DeleteFlightPlan(id)
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
		_, err = d.handleUploadImage(&oFile, iFile.Name, iFile.FileInfo().Size())
		if err != nil {
			log.Printf("UploadBatch: Cannot upload thie file %v, error is %v", iFile.Name, err)
			break
		}
		err = oFile.Close()
		if err != nil {
			return err
		}
		//log.Println(result)
	}

	log.Println("batch is uploaded")

	return nil
}

func extractMetadata(raw []byte) (int, ObservationMetadata, error) {

	// call exifTool
	cmd := exec.Command("exiftool", "-json", "-")
	cmd.Stdin = bytes.NewReader(raw)
	result, err := cmd.Output()
	if err != nil {
		//log.Fatalf("err: %v", err)
		return -1, ObservationMetadata{}, err
	}

	metadata := ObservationMetadata{
		Size:          12345678,
		Height:        1080,
		Width:         1920,
		Channels:      2,
		Timestamp:     123456789,
		BitsPixels:    6,
		ImageOffset:   24,
		Camera:        "I",
		GnssLongitude: 10.4058633,
		GnssLatitude:  553821913,
		GnssDate:      123456789,
		GnssTime:      123456789,
		GnssSpeed:     420,
		GnssAltitude:  17000,
		GnssCourse:    2,
	}

	// Unmarshal the EXIF data to a map of properties in the comment tag
	var t []map[string]string
	json.Unmarshal(result, &t)
	s := t[0]["Comment"]
	var l map[string]int
	err = json.Unmarshal([]byte(s), &l)
	if err != nil {
		return -1, ObservationMetadata{}, fmt.Errorf("extractMetadata: %v", err)
	}

	relatedObservationRequest := l["measurementRequest"]
	if relatedObservationRequest == 0 {
		return -1, ObservationMetadata{}, fmt.Errorf("extractMetadata: %v", err)
	}
	return relatedObservationRequest, metadata, nil
}
