package dim

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/discosat/storage-system/internal/Commands"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
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
	postgisContainer, err := postgres.Run(s.ctx,
		"postgis/postgis:17-3.4-alpine",
		postgres.WithInitScripts(filepath.Join(".", "create-database.sql"),
			filepath.Join(".", "testData", "test-data.sql")),
		postgres.WithInitScripts("./create-database.sql"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		s.T().Fatal(err)
	}
	s.pgContainer = postgisContainer

	minioContainer, err := minio.Run(s.ctx, "minio/minio:RELEASE.2025-04-22T22-12-26Z")
	if err != nil {
		s.T().Fatal(err)
	}
	s.minioContainer = minioContainer
	minioUrl, _ := s.minioContainer.ConnectionString(s.ctx)
	os.Setenv("MINIO_ENDPOINT", minioUrl)
	os.Setenv("MINIO_ACCESS_KEY_ID", s.minioContainer.Username)
	os.Setenv("MINIO_SECRET_ACCESS_KEY", s.minioContainer.Password)
	os.Setenv("MINIO_USE_SSL", "false")

	connString, err := postgisContainer.ConnectionString(s.ctx)
	db, err := sqlx.Open("pgx", connString)

	router := ConfigureRouter(
		NewDimController(
			NewDimService(
				observation.NewPsqlObservationRepository(db, objectStore.NewMinioStore()),
				observationRequest.NewPsqlObservationRequestRepository(db),
			),
		),
	)

	s.dimRouter = router
}

func (s *DimControllerIntegrationTestSuite) TestPingIntegration() {
	t := s.T()
	expected := `{"message":"pong"}`

	// WHEN
	request, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)
	response, _ := io.ReadAll(w.Body)

	// THEN
	assert.Equal(t, expected, string(response))
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

func TestDimControllerTestSuite(t *testing.T) {
	suite.Run(t, new(DimControllerIntegrationTestSuite))
}
