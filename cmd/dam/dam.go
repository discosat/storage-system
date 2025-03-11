package dam

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ImageReq struct {
	ImgID     string `form:"img_id" binding:"omitempty"`
	ObsID     string `form:"obs_id" binding:"omitempty"`
	StartTime string `form:"start_time" binding:"omitempty"`
	EndTime   string `form:"end_time" binding:"omitempty"`
	LatFrom   string `form:"lat_from" binding:"omitempty"`
	LatTo     string `form:"lat_to" binding:"omitempty"`
	LonFrom   string `form:"lon_from" binding:"omitempty"`
	LonTo     string `form:"lon_to" binding:"omitempty"`
	CamType   string `form:"cam_type" binding:"required"`
	Date      string `form:"date" binding:"omitempty"`
}

func Start() {
	g := gin.Default()

	g.GET("/search-images", func(context *gin.Context) {
		var req ImageReq
		if err := context.ShouldBindQuery(&req); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"message":  "Success",
			"recieved": req,
		})
	})
	g.Run(":8081")
}
