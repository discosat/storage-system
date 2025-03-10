package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var oClient OClient

func main() {
	ctx := context.Background()
	oClient = OClient{
		s3c: configueMinIO(),
	}

	bucketName := "disco2data"
	bucketRegion := "eu-north-0"

	err := oClient.s3c.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: bucketRegion})

	if err != nil {
		exists, existsErr := oClient.s3c.BucketExists(ctx, bucketName)
		if existsErr == nil && exists {
			log.Printf("We already own %s", bucketName)
		} else {
			log.Fatalf("main: %v", err)
		}
	}
	log.Printf("Successfulle created %s \n", bucketName)

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

func configueMinIO() *minio.Client {
	endpoint := "app-disco-minio.cloud.sdu.dk"
	accessKeyID := ""
	secretAccessKey := ""
	useSSL := true

	var err error
	minioC, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalf("configureMinIO: %v", err)
	}

	return minioC
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

	file, _, _ := c.Request.FormFile("batch")
	//fileName := filepath.Base(header.Filename)
	//log.Println(fileName)
	tmpFile, _ := os.CreateTemp("", "temp*.zip")
	defer os.Remove(tmpFile.Name())

	_, err := io.Copy(tmpFile, file)
	if err != nil {
		log.Fatalf("uploadBatch: %v \n", err)
	}
	log.Printf("uploadBatch: %+v \n", tmpFile)
	archive, err := zip.OpenReader(tmpFile.Name())
	for _, iFile := range archive.File {
		if iFile.FileInfo().IsDir() {
			continue
		}
		oFile, _ := iFile.Open()
		log.Printf("uploadBatch: File name: %v", iFile.Name)
		log.Printf("uploadBatch: filePath.base: %v", filepath.Base(iFile.Name))
		//status, errr := saveItem(oClient.s3c, iFile, oFile)
		status, errr := oClient.saveItem(oClient.s3c, iFile, oFile)
		//status, errr := oClient.s3c.PutObject(context.Background(), "disco2data", filepath.Base(iFile.Name), oFile, iFile.FileInfo().Size(), minio.PutObjectOptions{})
		//status, err := oClient.s3c.FPutObject(context.Background(), "disco2data", filepath.Base(iFile.Name), iFile.Name, minio.PutObjectOptions{ContentType: "image/png"})
		if errr != nil {
			log.Printf("uploadBatch: Cannot upload thie file %v, error is %v", filepath.Base(iFile.Name), err)
			break
		}
		log.Println(status.Key)

	}
	//_, err := io.Copy(tmpFile, file)
	//open, _ := file.Open()
	//io.Reader(open)

}
