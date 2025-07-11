package dam

import (
	"fmt"
	"github.com/discosat/storage-system/cmd/db"
	"github.com/discosat/storage-system/cmd/disco_qom"
	"github.com/discosat/storage-system/cmd/interfaces"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

func ConfigureRouter() (*gin.Engine, error) {
	g := gin.Default()

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Disposition"}, // So frontend can access download filename
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	pgClient, err := db.NewPostgresClient(os.Getenv("POSTGRES_CONN"))
	if err != nil {
		log.Fatal("Failed to init Postgres:", err)
	}

	minioClient, err := db.NewMinIOClient()
	if err != nil {
		log.Fatal("Failed to init MinIO:", err)
	}

	// Create handler
	handler := &DataAccessHandler{
		MetadataFetcher: pgClient,
		ImageFetcher:    minioClient,
	}

	g.GET("/images", handler.RequestHandler)

	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return g, nil
}

type DataAccessHandler struct {
	MetadataFetcher interfaces.MetadataFetcher
	ImageFetcher    interfaces.ImageFetcher
}

func (h *DataAccessHandler) RequestHandler(c *gin.Context) {
	var req interfaces.ImageRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	//Authenticating the request
	//auth := AuthService()
	//fmt.Println(auth)

	//Passing request on to QOM
	discoQO := &disco_qom.DiscoQO{}
	queryPusher := newQueryPusher(discoQO)

	//Optimized query + arguments gets returned
	sqlQuery, args, passingErr := queryPusher.PushQuery(req)
	if passingErr != nil {
		log.Fatal("Failed to pass query to QOM", passingErr)
	}

	//Calling db with SQL query string and arguments
	imageMetadata, PostgresErr := h.MetadataFetcher.FetchMetadata(sqlQuery, args)
	if PostgresErr != nil {
		log.Fatal("Failed to call PostgreSQL DB with SQL query", PostgresErr)
	}
	fmt.Println("Logging optimized query in dam.go: ", sqlQuery)
	fmt.Println("Logging optimized query arguments in dam.go: ", args)
	fmt.Println("Logging output of postgres service in dam.go: ", imageMetadata)

	var imageMinIOData []interfaces.ImageMinIOData
	for _, metadata := range imageMetadata {
		imageMinIOData = append(imageMinIOData, interfaces.ImageMinIOData{
			ObjectReference: metadata.ObjectReference,
			BucketName:      metadata.BucketName,
		})
	}

	//Calling Minio service with image IDs
	retrievedImages, minIOErr := h.ImageFetcher.FetchImages(imageMinIOData)
	if minIOErr != nil {
		log.Fatal("Failed to call MinIO service", minIOErr)
	}

	//bundling images together
	zipBundle, bundleErr := ImageBundler(retrievedImages, imageMetadata)
	if bundleErr != nil {
		log.Fatal("Failed to bundle images", bundleErr)
	}

	//Handle response
	ResponseHandler(c, zipBundle)
}

func ResponseHandler(c *gin.Context, zipBundle []byte) {
	if len(zipBundle) > 0 {
		filename := fmt.Sprintf("image_bundle_%d.zip", time.Now().Unix())
		fmt.Println("Logging filename of zip: ", filename)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(http.StatusOK, "application/zip", zipBundle)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "No images matched the query",
	})
}
