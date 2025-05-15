package interfaces

type ImageRequest struct {
	ImgID     *string  `form:"img_id" binding:"omitempty"`
	ObsID     *string  `form:"obs_id" binding:"omitempty"`
	StartTime *int64   `form:"start_time" binding:"omitempty"`
	EndTime   *int64   `form:"end_time" binding:"omitempty"`
	LatFrom   *float64 `form:"lat_from" binding:"omitempty"`
	LatTo     *float64 `form:"lat_to" binding:"omitempty"`
	LonFrom   *float64 `form:"lon_from" binding:"omitempty"`
	LonTo     *float64 `form:"lon_to" binding:"omitempty"`
	Camera    *string  `form:"camera" binding:"required"`
	Date      *int64   `form:"date" binding:"omitempty"`
}
