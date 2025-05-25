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
	"github.com/stretchr/testify/assert"
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
	t := s.T()
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
	assert.Equal(t, 200, w.Code)
	assert.NotNil(t, flightPlan)
	assert.Equal(t, flightPlan.Name, "flight plan 1")
	assert.NotNil(t, flightPlan.ObservationRequests)
	assert.Equal(t, flightPlan.ObservationRequests[0].OType, "image")
}

func (s *DimControllerIntegrationTestSuite) TestCreateFlightPlanIntegration() {
	t := s.T()

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
			{Id: 40, OType: "image"},
			{Id: 41, OType: "other"},
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
	assert.Equal(t, 201, w.Code)
	assert.Regexp(t, regexp.MustCompile(`{"flightPlanId":[0-9]+}`), string(response))
	//c.JSON(http.StatusCreated, gin.H)
}

func (s *DimControllerIntegrationTestSuite) TestCreateFlightPlanNoObservationRequestsIntegration() {
	t := s.T()

	// ----- GIVEN -----
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fpPart, _ := writer.CreateFormField("flightPlan")
	// New flight Plan with no observationRequests
	flightPlan := observationRequest.FlightPlanAggregate{
		Name:                "Integration Test Flight Plan",
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
	assert.Equal(t, 201, w.Code)
	assert.Regexp(t, regexp.MustCompile(`{"flightPlanId":[0-9]+}`), string(response))
	//c.JSON(http.StatusCreated, gin.H)
}

func (s *DimControllerIntegrationTestSuite) TestCreateFlightPlanIntegrationErrorObservationRequest() {
	t := s.T()

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
	assert.Equal(t, 400, w.Code)
	assert.Regexp(t, regexp.MustCompile(`{"error":"Observation request is formatted wrong"}`), string(response))
	//c.JSON(http.StatusCreated, gin.H)
}

func (s *DimControllerIntegrationTestSuite) TestUpdateFlightPlanIntegration() {
	// -----GIVEN-----
	//flight plan 2
	t := s.T()
	request, _ := http.NewRequest("GET", "/flightPlan?id=2", nil)
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	//var jsonMap map[string]observationRequest.FlightPlanAggregate
	//err := json.Unmarshal(response, &jsonMap)
	var flightPlan2 observationRequest.FlightPlanAggregate
	err := test.BindFlightPlanJson(response, &flightPlan2)
	if err != nil {
		log.Fatalf("Could not bind flightPlan: %v", err)
	}

	// ----- EXPECT -----
	assert.Equal(t, flightPlan2.Name, "flight plan 2")
	assert.Equal(t, flightPlan2.ObservationRequests[0].OType, "image")

	//----- WHEN -----
	//altered
	newFpName := "Nyt navn 2"
	newOrType := "number"
	flightPlan2.Name = newFpName
	flightPlan2.ObservationRequests[0].OType = newOrType

	//----- AND -----
	//requested
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	FpPart, err := writer.CreateFormField("flightPlan")
	if err != nil {
		t.Fatal(err)
	}
	// New flight Plan
	FpJson, err := json.Marshal(flightPlan2)
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
	request, _ = http.NewRequest("GET", "/flightPlan?id=2", nil)
	w = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ = io.ReadAll(w.Body)

	//var jsonMap map[string]observationRequest.FlightPlanAggregate
	var updatedFlightPlan observationRequest.FlightPlanAggregate
	err = test.BindFlightPlanJson(response, &updatedFlightPlan)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, updatedFlightPlan.Name, newFpName)
	assert.Equal(t, updatedFlightPlan.ObservationRequests[0].OType, newOrType)

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
	assert.Equal(t, flightPlan.Id, 3)
	assert.Equal(t, len(flightPlan.ObservationRequests), 2)

	// WHEN
	// Flight plan with id 2 is deleted
	request = httptest.NewRequest("DELETE", "/flightPlan?id=3", nil)
	responseWriter = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(responseWriter, request)
	response, _ = io.ReadAll(responseWriter.Body)
	//EXPECT
	// success
	assert.Equal(t, responseWriter.Code, 200)
	assert.Equal(t, string(response), `{"message":"flight plan: 3 has been deleted"}`)

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
	//assert.NotNil(t, expectedError)
	assert.Equal(t, responseWriter.Code, 404)
	assert.Equal(t, string(response), `{"error":"no flight plan with id: 3"}`)
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
	assert.Equal(t, flightPlan.Id, 2)
	assert.Equal(t, len(flightPlan.ObservationRequests), 2)

	// WHEN
	// Flight plan with id 2 is attempted deleted
	request = httptest.NewRequest("DELETE", "/flightPlan?id=2", nil)
	responseWriter = httptest.NewRecorder()
	s.dimRouter.ServeHTTP(responseWriter, request)
	_, _ = io.ReadAll(responseWriter.Body)
	//EXPECT
	// BadRequest
	assert.Equal(t, responseWriter.Code, 400)
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
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"ObservationIds":[2,3,4,5,6,7,8]}`, string(response))

}
