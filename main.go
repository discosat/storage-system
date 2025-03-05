package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/file", uploadFile)
	router.POST("/files", uploadFiles)
	router.POST("/batch", uploadBatch)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func uploadFile(c *gin.Context) {

	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	if _, err := os.Stat("./" + file.Filename); err == nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File already exists"))
		return
	}
	log.Println("Jeg når her")

	err := c.SaveUploadedFile(file, "./"+file.Filename)

	if err != nil {
		log.Fatalf("UploadBatch: %v", err)
	}

	c.String(http.StatusCreated, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func uploadFiles(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		log.Println(file.Filename)
		c.SaveUploadedFile(file, "./"+file.Filename)
	}
	c.String(http.StatusCreated, fmt.Sprintf("%d files have been uploaded", len(files)))
}

func uploadBatch(c *gin.Context) {

	file, _ := c.FormFile("batch")
	log.Println(file.Filename)
	tmpFile, _ := os.CreateTemp("", "temp.zip")
	defer os.Remove(tmpFile.Name())
	//_, err := io.Copy(tmpFile, file)
	open, _ := file.Open()
	//io.Reader(open)
	
}
