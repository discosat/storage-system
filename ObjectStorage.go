package main

import (
	"archive/zip"
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"path/filepath"
)

type OClient struct {
	s3c *minio.Client
}

func (o OClient) saveItem(S3Client *minio.Client, iFile *zip.File, oFile io.ReadCloser) (minio.UploadInfo, error) {
	status, err := S3Client.PutObject(context.Background(), "disco2data", filepath.Base(iFile.Name), oFile, iFile.FileInfo().Size(), minio.PutObjectOptions{})
	//log.Println(OClient.Test())
	return status, err
}
