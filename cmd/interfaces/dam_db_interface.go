package interfaces

type MetadataFetcher interface {
	FetchMetadata(query string, args []interface{}) ([]ImageMetadata, error)
}

type ImageFetcher interface {
	FetchImages(images []ImageMinIOData) ([]RetrievedImages, error)
}
