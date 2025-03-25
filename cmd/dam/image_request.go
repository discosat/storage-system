package dam

import (
	"time"
)

type ImageRequest struct {
	ImgID     *string    `form:"img_id" binding:"omitempty"`
	ObsID     *string    `form:"obs_id" binding:"omitempty"`
	StartTime *time.Time `form:"start_time" format:"UnixDate" binding:"omitempty"`
	EndTime   *time.Time `form:"end_time" format:"UnixDate" binding:"omitempty"`
	LatFrom   *float64   `form:"lat_from" binding:"omitempty"`
	LatTo     *float64   `form:"lat_to" binding:"omitempty"`
	LonFrom   *float64   `form:"lon_from" binding:"omitempty"`
	LonTo     *float64   `form:"lon_to" binding:"omitempty"`
	CamType   string     `form:"cam_type" binding:"required"`
	Date      *time.Time `form:"date" binding:"omitempty"`
}

type ImageReqInterface interface {
	GetImageReq() ImageRequest
}
