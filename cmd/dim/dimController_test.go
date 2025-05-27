package dim

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/discosat/storage-system/cmd/test"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"testing"
)

type DimControllerIntegrationTestSuite struct {
	suite.Suite
	pgContainer    *postgres.PostgresContainer
	minioContainer *minio.MinioContainer
	dimRouter      *gin.Engine
	ctx            context.Context
}

func (s *DimControllerIntegrationTestSuite) SetupSuite() {
	t := s.T()
	s.ctx = context.Background()
	postgisContainer, err := test.SetupPostgresContainer(s.ctx)
	if err != nil {
		s.T().Fatalf("Postgis container was not setup correctly: %v", err)
	}
	s.pgContainer = postgisContainer
	connString, err := postgisContainer.ConnectionString(s.ctx)
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("pgx", connString)
	if err != nil {
		t.Fatal(err)
	}

	minioContainer, err := test.SetupMinioContainer(s.ctx)
	if err != nil {
		s.T().Fatalf("Minio container was not setup correctly: %v", err)
	}
	s.minioContainer = minioContainer
	minioStore := objectStore.NewMinioStore()
	err = minioStore.CreateBucket("testbucket")
	if err != nil {
		s.T().Fatal(err)
	}
	router := ConfigureRouter(
		NewDimController(
			NewDimService(
				observation.NewPsqlObservationRepository(db, minioStore),
				observationRequest.NewPsqlObservationRequestRepository(db),
			),
		),
	)

	s.dimRouter = router
}

func TestDimControllerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(DimControllerIntegrationTestSuite))
}

func (s *DimControllerIntegrationTestSuite) TestGetFlightPlanIntegration() {
	request, _ := http.NewRequest("GET", "/flightPlan?id=1", nil)
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	var jsonMap map[string]observationRequest.FlightPlanAggregate
	err := json.Unmarshal(response, &jsonMap)
	if err != nil {
		log.Fatal("Could not bind flightPlan")
	}
	flightPlan := jsonMap["flightPlan"]

	// ----- THEN -----
	s.Equal(200, w.Code)
	s.NotNil(flightPlan)
	s.Equal(flightPlan.Name, "flight plan 1")
	s.NotNil(flightPlan.ObservationRequests)
	s.Equal(flightPlan.ObservationRequests[0].OType, "image")
}

func (s *DimControllerIntegrationTestSuite) TestCreateFlightPlanIntegration() {

	// ----- GIVEN -----
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fpPart, _ := writer.CreateFormField("flightPlan")
	// New flight Plan
	flightPlan := observationRequest.FlightPlanAggregate{
		Name:      "Integration Test Flight Plan",
		UserId:    1,
		MissionId: 1,
		ObservationRequests: []observationRequest.ObservationRequestDTO{
			{Id: 123, OType: "image"},
			{Id: 124, OType: "other"},
		},
	}
	fpJson, _ := json.Marshal(flightPlan)
	fpPart.Write(fpJson)

	// With two observation requests
	// observation request 1
	writer.Close()

	// ----- WHEN -----
	request, _ := http.NewRequest("POST", "/flightPlan", body)
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	//request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	// ----- THEN -----
	s.Equal(201, w.Code)
	s.Regexp(regexp.MustCompile(`{"flightPlanId":[0-9]+}`), string(response))
	//c.JSON(http.StatusCreated, gin.H)
}

func (s *DimControllerIntegrationTestSuite) TestCreateFlightPlanNoObservationRequestsIntegration() {

	// ----- GIVEN -----
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fpPart, _ := writer.CreateFormField("flightPlan")
	// New flight Plan with no observationRequests
	flightPlan := observationRequest.FlightPlanAggregate{
		Name:                "Integration Test Flight Plan No Observation Requests",
		UserId:              1,
		MissionId:           1,
		ObservationRequests: nil,
	}
	fpJson, _ := json.Marshal(flightPlan)
	fpPart.Write(fpJson)
	writer.Close()

	// ----- WHEN -----
	request, _ := http.NewRequest("POST", "/flightPlan", body)
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	//request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	// ----- THEN -----
	s.Equal(201, w.Code)
	s.Regexp(regexp.MustCompile(`{"flightPlanId":[0-9]+}`), string(response))
	//c.JSON(http.StatusCreated, gin.H)
}

func (s *DimControllerIntegrationTestSuite) TestCreateFlightPlanIntegrationErrorObservationRequest() {

	// ----- GIVEN -----
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fpPart, _ := writer.CreateFormField("flightPlan")
	// New flight Plan
	flightPlan := observationRequest.FlightPlanAggregate{
		Name:      "Integration Test Flight Plan With Error",
		UserId:    1,
		MissionId: 1,
		ObservationRequests: []observationRequest.ObservationRequestDTO{
			{Id: 40, OType: "this should give an error"},
			{Id: 41, OType: "other"},
		},
	}
	fpJson, _ := json.Marshal(flightPlan)
	fpPart.Write(fpJson)
	writer.Close()

	// ----- WHEN -----
	request, _ := http.NewRequest("POST", "/flightPlan", body)
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	//request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	// ----- THEN -----
	s.Equal(400, w.Code)
	s.Regexp(regexp.MustCompile(`{"error":"Observation request is formatted wrong"}`), string(response))
	//c.JSON(http.StatusCreated, gin.H)
}

func (s *DimControllerIntegrationTestSuite) TestUpdateFlightPlanIntegration() {
	// -----GIVEN-----
	//flight plan 2
	t := s.T()
	request, _ := http.NewRequest("GET", "/flightPlan?id=4", nil)
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	//var jsonMap map[string]observationRequest.FlightPlanAggregate
	//err := json.Unmarshal(response, &jsonMap)
	var flightPlan observationRequest.FlightPlanAggregate
	err := test.BindFlightPlanJson(response, &flightPlan)
	if err != nil {
		log.Fatalf("Could not bind flightPlan: %v", err)
	}

	// ----- EXPECT -----
	s.Equal(flightPlan.Name, "flight plan update test")
	s.Equal(flightPlan.ObservationRequests[0].OType, "image")

	//----- WHEN -----
	//altered
	newFpName := "Nyt navn 2"
	newOrType := "number"
	flightPlan.Name = newFpName
	flightPlan.ObservationRequests[0].OType = newOrType

	//----- AND -----
	//requested
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	FpPart, err := writer.CreateFormField("flightPlan")
	if err != nil {
		t.Fatal(err)
	}
	// New flight Plan
	FpJson, err := json.Marshal(flightPlan)
	if err != nil {
		t.Fatal(err)
	}
	_, err = FpPart.Write(FpJson)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	request, err = http.NewRequest("PUT", "/flightPlan", body)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	notNeeded := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(notNeeded, request)
	// ----- THEN -----
	//when Retrieved again
	request, _ = http.NewRequest("GET", "/flightPlan?id=4", nil)
	w = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ = io.ReadAll(w.Body)

	//var jsonMap map[string]observationRequest.FlightPlanAggregate
	var updatedFlightPlan observationRequest.FlightPlanAggregate
	err = test.BindFlightPlanJson(response, &updatedFlightPlan)
	if err != nil {
		t.Fatal(err)
	}

	s.Equal(http.StatusOK, w.Code)
	s.Equal(updatedFlightPlan.Name, newFpName)
	s.Equal(updatedFlightPlan.ObservationRequests[0].OType, newOrType)

}

func (s *DimControllerIntegrationTestSuite) TestUpdateFlightPlanAddObservationRequestIntegration() {
	// -----GIVEN-----
	//flight plan 2
	t := s.T()
	request, _ := http.NewRequest("GET", "/flightPlan?id=4", nil)
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	//var jsonMap map[string]observationRequest.FlightPlanAggregate
	//err := json.Unmarshal(response, &jsonMap)
	var flightPlan observationRequest.FlightPlanAggregate
	err := test.BindFlightPlanJson(response, &flightPlan)
	if err != nil {
		log.Fatalf("Could not bind flightPlan: %v", err)
	}
	numOriginalObservationRequests := len(flightPlan.ObservationRequests)

	// ----- EXPECT -----
	s.Equal(http.StatusOK, w.Code)

	//----- WHEN -----
	//Observation request is added
	flightPlan.ObservationRequests =
		append(flightPlan.ObservationRequests, observationRequest.ObservationRequestDTO{Id: 456, OType: "image"})

	//----- AND -----
	// Updated
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	FpPart, err := writer.CreateFormField("flightPlan")
	if err != nil {
		t.Fatal(err)
	}
	// New flight Plan
	FpJson, err := json.Marshal(flightPlan)
	if err != nil {
		t.Fatal(err)
	}
	_, err = FpPart.Write(FpJson)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	request, err = http.NewRequest("PUT", "/flightPlan", body)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	w = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)

	// EXPECT
	// Success
	s.Equal(http.StatusOK, w.Code)

	// ----- THEN -----
	//when Retrieved again
	request, _ = http.NewRequest("GET", "/flightPlan?id=4", nil)
	w = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ = io.ReadAll(w.Body)

	var updatedFlightPlan observationRequest.FlightPlanAggregate
	err = test.BindFlightPlanJson(response, &updatedFlightPlan)
	if err != nil {
		t.Fatal(err)
	}

	// EXPECT
	// number of observation requests is one larger
	s.Equal(http.StatusOK, w.Code)
	s.Equal(numOriginalObservationRequests+1, len(updatedFlightPlan.ObservationRequests))
}

func (s *DimControllerIntegrationTestSuite) TestUpdateFlightPlanRemoveObservationRequestIntegration() {
	// -----GIVEN-----
	//flight plan 2
	t := s.T()
	request, _ := http.NewRequest("GET", "/flightPlan?id=5", nil)
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	//var jsonMap map[string]observationRequest.FlightPlanAggregate
	//err := json.Unmarshal(response, &jsonMap)
	var flightPlan observationRequest.FlightPlanAggregate
	err := test.BindFlightPlanJson(response, &flightPlan)
	if err != nil {
		log.Fatalf("Could not bind flightPlan: %v", err)
	}
	numOriginalObservationRequests := len(flightPlan.ObservationRequests)

	// ----- EXPECT -----
	s.Equal(http.StatusOK, w.Code)
	s.Equal(flightPlan.Name, "flight plan update delete test")

	//----- WHEN -----
	//Observation request is deleted
	flightPlan.ObservationRequests = slices.Delete(flightPlan.ObservationRequests, 0, 1)

	//----- AND -----
	// Updated
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	FpPart, err := writer.CreateFormField("flightPlan")
	if err != nil {
		t.Fatal(err)
	}
	// New flight Plan
	FpJson, err := json.Marshal(flightPlan)
	if err != nil {
		t.Fatal(err)
	}
	_, err = FpPart.Write(FpJson)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	request, err = http.NewRequest("PUT", "/flightPlan", body)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	w = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)

	// EXPECT
	// Success
	s.Equal(http.StatusOK, w.Code)

	// ----- THEN -----
	//when Retrieved again
	request, _ = http.NewRequest("GET", "/flightPlan?id=5", nil)
	w = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ = io.ReadAll(w.Body)

	var updatedFlightPlan observationRequest.FlightPlanAggregate
	err = test.BindFlightPlanJson(response, &updatedFlightPlan)
	if err != nil {
		t.Fatal(err)
	}

	// EXPECT
	// number of observation requests is one larger
	s.Equal(http.StatusOK, w.Code)
	s.Equal(numOriginalObservationRequests-1, len(updatedFlightPlan.ObservationRequests))
}

func (s *DimControllerIntegrationTestSuite) TestUpdateFlightPlanLockedErrorIntegration() {
	// -----GIVEN-----
	//flight plan 2
	t := s.T()
	request, _ := http.NewRequest("GET", "/flightPlan?id=2", nil)
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	//var jsonMap map[string]observationRequest.FlightPlanAggregate
	//err := json.Unmarshal(response, &jsonMap)
	var preUpdateFlightPlan observationRequest.FlightPlanAggregate
	err := test.BindFlightPlanJson(response, &preUpdateFlightPlan)
	if err != nil {
		log.Fatalf("Could not bind flightPlan: %v", err)
	}

	// ----- EXPECT -----
	s.Equal(http.StatusOK, w.Code)
	s.Equal(preUpdateFlightPlan.Id, 2)
	s.Equal(preUpdateFlightPlan.Locked, true)

	//----- WHEN -----
	//Flight plan is changed
	updatingFlightPlan := preUpdateFlightPlan
	updatingFlightPlan.Name = "This should give an error, as flight plan is locked"

	//----- AND -----
	// Updated
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	FpPart, err := writer.CreateFormField("flightPlan")
	if err != nil {
		t.Fatal(err)
	}
	// New flight Plan
	FpJson, err := json.Marshal(updatingFlightPlan)
	if err != nil {
		t.Fatal(err)
	}
	_, err = FpPart.Write(FpJson)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	updateRequest, err := http.NewRequest("PUT", "/flightPlan", body)
	if err != nil {
		t.Fatal(err)
	}
	updateRequest.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	updateWriter := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(updateWriter, updateRequest)
	updateResponse, _ := io.ReadAll(updateWriter.Body)

	// THEN
	// Error
	s.Equal(http.StatusBadRequest, updateWriter.Code)
	s.Equal(string(updateResponse), `{"error":"Flight plan with id: 2 is locked"}`)

	// AND
	// The id fetched again brings the same object

	newFetchRequest, _ := http.NewRequest("GET", "/flightPlan?id=2", nil)
	newFetchW := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(newFetchW, newFetchRequest)
	response, _ = io.ReadAll(newFetchW.Body)

	var postUpdateFlightPlan observationRequest.FlightPlanAggregate
	err = test.BindFlightPlanJson(response, &postUpdateFlightPlan)
	if err != nil {
		t.Fatal(err)
	}

	s.Equal(preUpdateFlightPlan, postUpdateFlightPlan)

}

func (s *DimControllerIntegrationTestSuite) TestDeleteFlightPlan() {
	t := s.T()
	// GIVEN
	// ObservationRequest with id 2 exists
	request := httptest.NewRequest("GET", "/flightPlan?id=3", nil)
	responseWriter := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(responseWriter, request)
	response, _ := io.ReadAll(responseWriter.Body)
	var flightPlan observationRequest.FlightPlanAggregate
	err := test.BindFlightPlanJson(response, &flightPlan)
	if err != nil {
		t.Fatal(err)
	}

	// EXPECT
	// Looks as expected
	s.Equal(flightPlan.Id, 3)
	s.Equal(len(flightPlan.ObservationRequests), 2)

	// WHEN
	// Flight plan with id 2 is deleted
	request = httptest.NewRequest("DELETE", "/flightPlan?id=3", nil)
	responseWriter = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(responseWriter, request)
	response, _ = io.ReadAll(responseWriter.Body)
	//EXPECT
	// success
	s.Equal(responseWriter.Code, 200)
	s.Equal(string(response), `{"message":"flight plan: 3 has been deleted"}`)

	// THEN
	// No flight plan with id 2 exists
	request = httptest.NewRequest("GET", "/flightPlan?id=3", nil)
	responseWriter = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(responseWriter, request)
	response, _ = io.ReadAll(responseWriter.Body)
	//var emptyFlightPlan observationRequest.FlightPlanAggregate
	//expectedError := test.BindFlightPlanJson(response, &emptyFlightPlan)

	// EXPECT
	// Looks as expected
	//s.NotNil(expectedError)
	s.Equal(responseWriter.Code, 404)
	s.Equal(string(response), `{"error":"no flight plan with id: 3"}`)
}

func (s *DimControllerIntegrationTestSuite) TestDeleteFlightPlanWithObservationsError() {
	t := s.T()
	// GIVEN
	// ObservationRequest with id 2 exists, that has had an observation uploaded to it
	request := httptest.NewRequest("GET", "/flightPlan?id=2", nil)
	responseWriter := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(responseWriter, request)
	response, _ := io.ReadAll(responseWriter.Body)
	var flightPlan observationRequest.FlightPlanAggregate
	err := test.BindFlightPlanJson(response, &flightPlan)
	if err != nil {
		t.Fatal(err)
	}

	// EXPECT
	// Looks as expected
	s.Equal(flightPlan.Id, 2)
	s.Equal(len(flightPlan.ObservationRequests), 2)

	// WHEN
	// Flight plan with id 2 is attempted deleted
	request = httptest.NewRequest("DELETE", "/flightPlan?id=2", nil)
	responseWriter = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(responseWriter, request)
	_, _ = io.ReadAll(responseWriter.Body)
	//EXPECT
	// BadRequest
	s.Equal(responseWriter.Code, 400)
}

// TODO husk også hentning af image_series
// TODO Test sletning af en observation request (via opdater der har en observation. Skal slå fejl (Id 6)

func (s *DimControllerIntegrationTestSuite) TestUploadBatch() {
	t := s.T()

	// GIVEN
	// a batch (.zip) of images
	//reader, err := zip.OpenReader(filepath.Join(".", "testData", "batch.zip"))
	zip, err := os.Open(filepath.Join(".", "testData", "batch.zip"))
	if err != nil {
		t.Fatal(err)
	}
	defer zip.Close()
	// ----- GIVEN -----
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	batchPart, _ := writer.CreateFormFile("batch", filepath.Base(zip.Name()))
	io.Copy(batchPart, zip)
	writer.Close()

	// ----- WHEN -----
	// Batch is uploaded
	request, _ := http.NewRequest("POST", "/batch", body)
	request.Header.Set("Content-Type", writer.FormDataContentType()+"; boundary="+writer.Boundary())
	//request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	// THEN
	s.Equal(http.StatusCreated, w.Code)
	s.Equal(`{"ObservationIds":[2,3,4,5,6,7,8]}`, string(response))

}
