package main

import (
	"archive/zip"
)

type IDataStore interface {
	SaveBatch(zipArchive *zip.ReadCloser, bucketName string) (string, error)
}

type SimpleStore struct {
	ds IDataStore
}

func NewSimpleStore(ds IDataStore) SimpleStore {
	return SimpleStore{
		ds: ds,
	}
}
