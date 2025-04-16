package measurementMetadata

type MeasurementMetadataRepository interface {
	GetById(id int) (MeasurementMetadata, error)
	CreateMeasurementMetadata(measurementId int, long float64, lat float64) (int, error)
}
