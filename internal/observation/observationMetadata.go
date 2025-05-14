package observation

import "time"

type ObservationMetadata struct {
	Id            int       `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObservationId int       `json:"observation_id"`
	Size          int       `json:"size"`
	Height        int       `json:"height"`
	Width         int       `json:"width"`
	Channels      int       `json:"channels"`
	Timestamp     int64     `json:"timestamp"`
	BitsPixels    int       `json:"bits_pixels"`
	ImageOffset   int       `json:"image_offset"`
	Camera        string    `json:"camera"`
	GnssLongitude float64   `json:"gnss_longitude"`
	GnssLatitude  int       `json:"gnss_latitude"`
	GnssDate      int       `json:"gnss_date"`
	GnssTime      int       `json:"gnss_time"`
	GnssSpeed     float64   `json:"gnss_speed"`
	GnssAltitude  float64   `json:"gnss_altitude"`
	GnssCourse    float64   `json:"gnss_course"`
}
