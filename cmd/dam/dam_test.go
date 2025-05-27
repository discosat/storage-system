package dam

import (
	"context"
	"github.com/discosat/storage-system/test"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/suite"
	minioTestCont "github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type DamControllerIntegrationTestSuite struct {
	suite.Suite
	pgContainer    *postgres.PostgresContainer
	minioContainer *minioTestCont.MinioContainer
	damRouter      *gin.Engine
	ctx            context.Context
}

func (s *DamControllerIntegrationTestSuite) SetupSuite() {
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
	os.Setenv("POSTGRES_CONN", connString)

	minioContainer, err := test.SetupMinioContainer(s.ctx)
	if err != nil {
		s.T().Fatalf("Minio container was not setup correctly: %v", err)
	}
	s.minioContainer = minioContainer

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")

	minioTestSetupClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		t.Fatal("Could not create minio test setup client")
	}

	testBucket := "testbucket"
	err = minioTestSetupClient.MakeBucket(s.ctx, testBucket, minio.MakeBucketOptions{Region: "eu-north-0"})
	if err != nil {
		t.Fatal("Could not create bucket")
	}
	test.PopulateMinioForDam(minioTestSetupClient, testBucket, s.ctx)

	damRouter, err := ConfigureRouter()
	if err != nil {
		t.Fatal("Cant create router")
	}
	s.damRouter = damRouter

}

func TestDimControllerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(DamControllerIntegrationTestSuite))
}

func (s *DamControllerIntegrationTestSuite) TestFetchById() {
	t := s.T()
	request, err := http.NewRequest("GET", "/images?camera=W", nil)
	if err != nil {
		t.Fatal("request fail")
	}
	responseW := httptest.NewRecorder()
	s.damRouter.ServeHTTP(responseW, request)
	response, _ := io.ReadAll(responseW.Body)
	//log.Println(response)
	os.WriteFile(filepath.Join("..", "..", "testData", "bundle.zip"), response, 0644)
	//resp, err := http.DefaultClient.Get("http://0.0.0.0:8081/images?camera=w")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//log.Println(resp)
}
