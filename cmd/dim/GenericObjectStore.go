package dim

import (
	"archive/zip"
	"io"
)

type IDataStore interface {
	SaveFile(fileInfo *zip.File, openFile io.ReadCloser, bucketName string) (string, error)
	BucketExists(bucketName string) (bool, error)
}

type SimpleStore struct {
	ds IDataStore
}

func NewSimpleStore(ds IDataStore) SimpleStore {
	return SimpleStore{
		ds: ds,
	}
}
