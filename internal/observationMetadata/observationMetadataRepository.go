package observationMetadata

type ObservationMetadataRepository interface {
	GetById(id int) (ObservationMetadata, error)
	CreateObservationMetadata(observationId int, long float64, lat float64) (int, error)
}
