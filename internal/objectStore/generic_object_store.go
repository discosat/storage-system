package objectStore

import (
	"archive/zip"
	"github.com/discosat/storage-system/internal/Commands"
	"io"
)

type IDataStore interface {
	SaveFile(fileInfo *zip.File, openFile io.ReadCloser, bucketName string) (string, error)
	SaveImage(observationCommand Commands.ObservationCommand) (string, error)
	BucketExists(bucketName string) (bool, error)
	DeleteImage(imgRef string, bucketName string) (bool, error)
}
