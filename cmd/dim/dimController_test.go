package dim

import (
	"context"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
	postgisContainer, err := postgres.Run(s.ctx,
		"postgis/postgis:17-3.4-alpine",
		postgres.WithInitScripts("./create-database.sql"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
	)

	if err != nil {
		log.Fatal(err)
	}
	s.pgContainer = postgisContainer

	minioContainer, err := minio.Run(s.ctx, "minio/minio:RELEASE.2025-04-22T22-12-26Z")
	if err != nil {
		log.Fatal(err)
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

func (s *DimControllerIntegrationTestSuite) TestCreateFlightPlanIntegration() {
	t := s.T()
	expected := `{"message":"pong"}`

	request, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	s.dimRouter.ServeHTTP(w, request)

	response, _ := io.ReadAll(w.Body)
	assert.Equal(t, expected, string(response))
	println(response)

}

func TestDimControllerTestSuite(t *testing.T) {
	suite.Run(t, new(DimControllerIntegrationTestSuite))
}
