package dim

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var ObjectStore SimpleStore

func Start() {
	ObjectStore = NewSimpleStore(NewMinioStore())
	router := gin.Default()
	//router.Use(ErrorHandler)

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

	exists, err := ObjectStore.ds.BucketExists(bucketName)
	if err != nil {
		ErrortAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	if !exists {
		e := fmt.Errorf("no bucket with name '%v' exists. Specify an already existing bucket", bucketName)
		ErrortAbortMessage(c, http.StatusBadRequest, e)
		return
	}

	file, _, err := c.Request.FormFile("batch")
	if err != nil {
		ErrortAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	tmpFile, _ := os.CreateTemp("", "temp*.zip")
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		ErrortAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	archive, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		ErrortAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	//var success string
	for _, iFile := range archive.File {
		if iFile.FileInfo().IsDir() {
			continue
		}
		oFile, _ := iFile.Open()
		_, err = ObjectStore.ds.SaveFile(iFile, oFile, bucketName)

		if err != nil {
			log.Printf("UploadBatch: Cannot upload thie file %v, error is %v", filepath.Base(iFile.Name), err)
			break
		}
	}

	log.Println("done")

}

func ErrortAbortMessage(c *gin.Context, statusCode int, err error) {
	log.Println(err)
	c.JSON(statusCode, gin.H{"error": fmt.Sprint(err)})
}

//func ErrorHandler(c *gin.Context) {
//	c.Next()
//
//	for _, err := range c.Errors {
//		log.Println(err)
//		c.JSON(400, err)
//
//	}
//
//	//c.JSON(400, "dø lige")
//}
