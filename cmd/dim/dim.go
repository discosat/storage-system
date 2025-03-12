package dim

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var ObjectStore SimpleStore

func Start() {
	ObjectStore = NewSimpleStore(newMinioStore())
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/file", uploadFile)
	router.POST("/files", uploadFiles)
	router.POST("/batch", UploadBatch)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func UploadBatch(c *gin.Context) {

	bucketName := c.Request.FormValue("bucketName")
	file, _, _ := c.Request.FormFile("batch")

	tmpFile, _ := os.CreateTemp("", "temp*.zip")
	defer os.Remove(tmpFile.Name())

	_, err := io.Copy(tmpFile, file)
	if err != nil {
		log.Fatalf("UploadBatch: %v \n", err)
	}

	archive, err := zip.OpenReader(tmpFile.Name())
	success, err := ObjectStore.ds.SaveBatch(archive, bucketName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(success)

}
