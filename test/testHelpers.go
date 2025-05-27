package test

import (
	"context"
	"encoding/json"
	"github.com/discosat/storage-system/internal/observationRequest"
	"github.com/minio/minio-go/v7"
	"github.com/testcontainers/testcontainers-go"
	minioTestCont "github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
	"time"
)

func SetupPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	postgisContainer, err := postgres.Run(ctx,
		"postgis/postgis:17-3.4-alpine",
		postgres.WithInitScripts(filepath.Join("..", "..", "sql", "create-database.sql"),
			filepath.Join("..", "..", "testData", "test-data.sql")),
		//postgres.WithInitScripts("./create-database.sql"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	return postgisContainer, err
}
func SetupMinioContainer(ctx context.Context) (*minioTestCont.MinioContainer, error) {
	minioContainer, err := minioTestCont.Run(ctx, "minio/minio:RELEASE.2025-04-22T22-12-26Z")
	minioUrl, _ := minioContainer.ConnectionString(ctx)
	os.Setenv("MINIO_ENDPOINT", minioUrl)
	os.Setenv("MINIO_ACCESS_KEY_ID", minioContainer.Username)
	os.Setenv("MINIO_SECRET_ACCESS_KEY", minioContainer.Password)
	os.Setenv("MINIO_USE_SSL", "false")
	return minioContainer, err
}

func BindFlightPlanJson(source []byte, sink *observationRequest.FlightPlanAggregate) error {

	var jsonMap map[string]observationRequest.FlightPlanAggregate
	err := json.Unmarshal(source, &jsonMap)
	if err != nil {
		return err
	}
	*sink = jsonMap["flightPlan"]
	return nil
}

func PopulateMinioForDam(client *minio.Client, bucketName string, ctx context.Context) {

	file, err := os.Open(filepath.Join("..", "..", "testData", "Ooo.jpg"))
	if err != nil {
		log.Fatal("File error")
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal("Fileinfo error")
	}

	_, err = client.PutObject(ctx, bucketName, fileInfo.Name(), file, fileInfo.Size(), minio.PutObjectOptions{})
	if err != nil {
		log.Fatal("Could not create test-setup minio client")
	}

	file, err = os.Open(filepath.Join("..", "..", "testData", "forgor.jpg"))
	if err != nil {
		log.Fatal("File error")
	}
	fileInfo, err = file.Stat()
	if err != nil {
		log.Fatal("Fileinfo error")
	}

	_, err = client.PutObject(ctx, bucketName, fileInfo.Name(), file, fileInfo.Size(), minio.PutObjectOptions{})
	if err != nil {
		log.Fatal("Could not create test-setup minio client")
	}

	file, err = os.Open(filepath.Join("..", "..", "testData", "Mary.jpg"))
	if err != nil {
		log.Fatal("File error")
	}
	fileInfo, err = file.Stat()
	if err != nil {
		log.Fatal("Fileinfo error")
	}

	_, err = client.PutObject(ctx, bucketName, fileInfo.Name(), file, fileInfo.Size(), minio.PutObjectOptions{})
	if err != nil {
		log.Fatal("Could not create test-setup minio client")
	}

	file, err = os.Open(filepath.Join("..", "..", "testData", "goblinMode.jpg"))
	if err != nil {
		log.Fatal("File error")
	}
	fileInfo, err = file.Stat()
	if err != nil {
		log.Fatal("Fileinfo error")
	}

	_, err = client.PutObject(ctx, bucketName, fileInfo.Name(), file, fileInfo.Size(), minio.PutObjectOptions{})
	if err != nil {
		log.Fatal("Could not create test-setup minio client")
	}
}
