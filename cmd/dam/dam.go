package dam

import (
	"fmt"
	"github.com/discosat/storage-system/cmd/disco_qom"
	"github.com/discosat/storage-system/cmd/interfaces"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Start() {
	g := gin.Default()

	g.GET("/search-images", RequestHandler)

	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	err := g.Run(":8081")
	if err != nil {
		log.Fatal("Failed to start server")
	}
}

func RequestHandler(c *gin.Context) {
	var req interfaces.ImageRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	//Authenticating the request
	auth := AuthService()
	fmt.Println(auth)

	//Passing request on to QOM
	discoQO := &disco_qom.DiscoQO{}
	queryPusher := newQueryPusher(discoQO)

	//Optimized query + arguments gets returned
	sqlQuery, args, passingErr := queryPusher.PushQuery(req)
	if passingErr != nil {
		log.Fatal("Failed to pass query to QOM", passingErr)
	}

	//Calling db with SQL query string and arguments
	imageMetadata, PostgresErr := PostgresService(sqlQuery, args)
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
	retrievedImages, err := MinIOService(imageMinIOData)
	if err != nil {
		log.Fatal("Failed to call MinIO service", err)
	}

	//bundling images together
	//imageFound := ImageBundler()

	//Handle response
	ResponseHandler(c, retrievedImages)
}

func ResponseHandler(c *gin.Context, retrievedImages []interfaces.RetrievedImages) {
	if len(retrievedImages) > 0 {
		imageData := retrievedImages[0].Image

		c.Data(http.StatusOK, "image/jpeg", imageData)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "No images matched the query",
	})
}
