package measurement

import "mime/multipart"

type MeasurementRepository interface {
	GetById(id int) (Measurement, error)
	CreateMeasurement(file *multipart.FileHeader, bucket string, requestName string, observationId int, measurementRequestId int) (int, error)
}
