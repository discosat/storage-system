package interfaces

type ImageMetadata struct {
	ID              int
	CreatedAt       string
	UpdatedAt       string
	ObservationID   int
	Size            int
	Height          int
	Width           int
	Channels        int
	Timestamp       int64
	BitsPixels      int
	ImageOffset     int
	Camera          string
	Location        string
	GnssDate        int64
	GnssTime        int64
	GnssSpeed       float64
	GnssAltitude    float64
	GnssCourse      float64
	BucketName      string
	ObjectReference string
}

type ImageMinIOData struct {
	ObjectReference string
	BucketName      string
}

type RetrievedImages struct {
	ObjectReference string
	BucketName      string
	Image           []byte
}
