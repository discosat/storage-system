package test

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"path/filepath"
	"time"
)

func SetupPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	postgisContainer, err := postgres.Run(ctx,
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
	return postgisContainer, err
}
func SetupMinioContainer(ctx context.Context) (*minio.MinioContainer, error) {
	minioContainer, err := minio.Run(ctx, "minio/minio:RELEASE.2025-04-22T22-12-26Z")
	minioUrl, _ := minioContainer.ConnectionString(ctx)
	os.Setenv("MINIO_ENDPOINT", minioUrl)
	os.Setenv("MINIO_ACCESS_KEY_ID", minioContainer.Username)
	os.Setenv("MINIO_SECRET_ACCESS_KEY", minioContainer.Password)
	os.Setenv("MINIO_USE_SSL", "false")
	return minioContainer, err
}
