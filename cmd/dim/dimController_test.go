package dim

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/discosat/storage-system/cmd/test"
	"github.com/discosat/storage-system/internal/Commands"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
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
	s.ctx = context.Background()
	postgisContainer, err := test.SetupPostgresContainer(s.ctx)
	if err != nil {
		s.T().Fatalf("Postgis container was not setup correctly: %v", err)
	}
	s.pgContainer = postgisContainer
	connString, err := postgisContainer.ConnectionString(s.ctx)
	db, err := sql.Open("pgx", connString)

	minioContainer, err := test.SetupMinioContainer(s.ctx)
	if err != nil {
		s.T().Fatalf("Minio container was not setup correctly: %v", err)
	}
	s.minioContainer = minioContainer
	minioStore := objectStore.NewMinioStore()

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

func TestDimControllerTestSuite(t *testing.T) {
	suite.Run(t, new(DimControllerIntegrationTestSuite))
}

func (s *DimControllerIntegrationTestSuite) TestFlightPlanIntegration() {
	t := s.T()

	// ----- GIVEN -----
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fpPart, _ := writer.CreateFormField("flightPlan")
	// New flight Plan
	fpCommand := Commands.CreateFlightPlanCommand{
		Name:      "Test Flight Plan",
		UserId:    1,
		MissionId: 1,
	}
	fpJson, _ := json.Marshal(fpCommand)
	fpPart.Write(fpJson)

	// With two observation request
	// observation request 1
	orPart1, _ := writer.CreateFormField("requestList")
	orCommand := Commands.CreateObservationRequestCommand{
		OType: "image",
	}
	orJson, _ := json.Marshal(orCommand)
	orPart1.Write(orJson)

	// observation request 2
	orPart2, _ := writer.CreateFormField("requestList")
	orCommand = Commands.CreateObservationRequestCommand{
		OType: "other",
	}
	orJson, _ = json.Marshal(orCommand)
	orPart2.Write(orJson)
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

func (s *DimControllerIntegrationTestSuite) TestFlightPlanNoObservationRequestsIntegration() {
	t := s.T()

	// ----- GIVEN -----
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fpPart, _ := writer.CreateFormField("flightPlan")
	// New flight Plan with no observationRequests
	fpCommand := Commands.CreateFlightPlanCommand{
		Name:      "Flight Plan no requests",
		UserId:    1,
		MissionId: 1,
	}
	fpJson, _ := json.Marshal(fpCommand)
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

func (s *DimControllerIntegrationTestSuite) TestFlightPlanIntegrationErrorObservationRequest() {
	t := s.T()

	// ----- GIVEN -----
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fpPart, _ := writer.CreateFormField("flightPlan")
	// New flight Plan
	fpCommand := Commands.CreateFlightPlanCommand{
		Name:      "Test Flight Plan error requests",
		UserId:    1,
		MissionId: 1,
	}
	fpJson, _ := json.Marshal(fpCommand)
	fpPart.Write(fpJson)

	// With two observation request
	// observation request 1
	orPart1, _ := writer.CreateFormField("requestList")
	orCommand := Commands.CreateObservationRequestCommand{
		OType: "this should give an error",
	}
	orJson, _ := json.Marshal(orCommand)
	orPart1.Write(orJson)
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
