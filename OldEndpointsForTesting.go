package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

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
