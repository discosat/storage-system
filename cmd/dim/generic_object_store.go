package dim

import (
	"archive/zip"
	"io"
	"mime/multipart"
)

type IDataStore interface {
	SaveFile(fileInfo *zip.File, openFile io.ReadCloser, bucketName string) (string, error)
	SaveImage(fileHeader *multipart.FileHeader, openFile io.ReadCloser, bucketName string, observationName string) (string, error)
	BucketExists(bucketName string) (bool, error)
	DeleteImage(imgRef string, bucketName string) (bool, error)
}

type SimpleStore struct {
	ds IDataStore
}

func NewSimpleStore(ds IDataStore) SimpleStore {
	return SimpleStore{
		ds: ds,
	}
}
