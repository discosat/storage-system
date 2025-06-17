package dim

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/discosat/storage-system/internal/Commands"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/observationRequest"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"log"
	"log/slog"
	"os/exec"
)

type DimServiceInterface interface {
	handleUploadImage(file *io.ReadCloser, fileName string, fileSize int64) (int, error)
	handleUploadBatch(archive *zip.ReadCloser) ([]int, error)
	handleGetFlightPlan(id int) (observationRequest.FlightPlanAggregate, error)
	handleCreateFlightPlan(flightPlan observationRequest.FlightPlanAggregate) (int, error)
	handleUpdateFlightPlan(flightPlan observationRequest.FlightPlanAggregate) (int, error)
	handleDeleteFlightPlan(id int) (bool, error)
}

type DimService struct {
	observationRequestRepository observationRequest.ObservationRequestRepository
	observationRepository        observation.ObservationRepository
}

func NewDimService(oRepo observation.ObservationRepository, orRepo observationRequest.ObservationRequestRepository) *DimService {
	return &DimService{
		observationRequestRepository: orRepo,
		observationRepository:        oRepo,
	}
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
	observationCommand := Commands.CreateObservationCommand{File: fileReader, FileName: fileName, FileSize: fileSize,
		Bucket:               observationRequestAggr.Bucket,
		FlightPlanName:       observationRequestAggr.FlightPlanName,
		ObservationRequestId: observationRequestAggr.ObservationRequest.Id,
		ObservationType:      observationRequestAggr.ObservationRequest.OType}
	// Saves image
	observationId, err := d.observationRepository.CreateObservation(observationCommand, &metadata)
	if err != nil {
		return -1, err
	}

	return observationId, nil
}

func (d DimService) handleUploadBatch(archive *zip.ReadCloser) ([]int, error) {
	uploadedIds := make([]int, 0)
	for _, iFile := range archive.File {
		if iFile.FileInfo().IsDir() {
			continue
		}
		oFile, err := iFile.Open()
		if err != nil {
			return nil, err
		}
		//oFile, err = iFile.Open()
		id, err := d.handleUploadImage(&oFile, iFile.Name, iFile.FileInfo().Size())
		uploadedIds = append(uploadedIds, id)
		if err != nil {
			log.Printf("UploadBatch: Cannot upload thie file %v, error is %v", iFile.Name, err)
			return nil, err
		}
		err = oFile.Close()
		if err != nil {
			return nil, err
		}
		//log.Println(result)
	}

	log.Println("batch is uploaded")

	return uploadedIds, nil
}

func (d DimService) handleGetFlightPlan(id int) (observationRequest.FlightPlanAggregate, error) {
	return d.observationRequestRepository.GetFlightPlanById(id)
}

func (d DimService) handleCreateFlightPlan(flightPlan observationRequest.FlightPlanAggregate) (int, error) {
	return d.observationRequestRepository.CreateFlightPlan(flightPlan)
}

func (d DimService) handleUpdateFlightPlan(flightPlan observationRequest.FlightPlanAggregate) (int, error) {
	return d.observationRequestRepository.UpdateFlightPlan(flightPlan)

}

func (d DimService) handleDeleteFlightPlan(id int) (bool, error) {
	return d.observationRequestRepository.DeleteFlightPlan(id)
}

func extractMetadata(raw []byte) (int, observation.ObservationMetadata, error) {

	// call exifTool
	cmd := exec.Command("exiftool", "-json", "-")
	cmd.Stdin = bytes.NewReader(raw)
	result, err := cmd.Output()
	if err != nil {
		//log.Fatalf("err: %v", err)
		return -1, observation.ObservationMetadata{}, err
	}

	metadata := observation.ObservationMetadata{
		Size:          12345678,
		Height:        1080,
		Width:         1920,
		Channels:      2,
		Timestamp:     123456789,
		BitsPixels:    6,
		ImageOffset:   24,
		Camera:        "",
		GnssLongitude: 0,
		GnssLatitude:  0,
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
	var l map[string]any
	err = json.Unmarshal([]byte(s), &l)
	if err != nil {
		return -1, observation.ObservationMetadata{}, fmt.Errorf("extractMetadata: %v", err)
	}

	relatedObservationRequest := int(l["measurementRequest"].(float64))
	if relatedObservationRequest == 0 {
		return -1, observation.ObservationMetadata{}, fmt.Errorf("extractMetadata: %v", err)
	}
	lon := l["lon"]
	if lon == 0 {
		return -1, observation.ObservationMetadata{}, fmt.Errorf("extractMetadata: %v", err)
	} else {
		metadata.GnssLongitude = lon.(float64)
	}
	lat := l["lat"]
	if lat == 0 {
		return -1, observation.ObservationMetadata{}, fmt.Errorf("extractMetadata: %v", err)
	} else {
		metadata.GnssLatitude = int(lat.(float64))
	}
	cam := l["cam"]
	if cam == "" {
		return -1, observation.ObservationMetadata{}, fmt.Errorf("extractMetadata: %v", err)
	} else {
		metadata.Camera = cam.(string)
	}
	return relatedObservationRequest, metadata, nil
}
