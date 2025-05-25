package objectStore

import (
	"archive/zip"
	"github.com/discosat/storage-system/internal/Commands"
	"io"
)

type IDataStore interface {
	SaveFile(fileInfo *zip.File, openFile io.ReadCloser, bucketName string) (string, error)
	SaveObservation(observationCommand Commands.ObservationCommand) (string, error)
	BucketExists(bucketName string) (bool, error)
	DeleteObservation(imgRef string, bucketName string) (bool, error)
}
