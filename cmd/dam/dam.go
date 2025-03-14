package dam

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"time"
)

type ImageReq struct {
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

func Start() {
	g := gin.Default()

	g.GET("/search-images", ReqReturn)

	err := g.Run(":8081")
	if err != nil {
		log.Fatal("Failed to start server")
	}
}

func FilterEmptyFields(req ImageReq) map[string]interface{} {
	filteredResult := make(map[string]interface{})
	val := reflect.ValueOf(req)
	typ := reflect.TypeOf(req)

	for i := 0; i < val.NumField(); i++ {
		formVal := val.Field(i)
		formType := typ.Field(i).Tag.Get("form")

		if formType == "" {
			formType = typ.Field(i).Name
		}

		if formVal.Kind() == reflect.Ptr && !formVal.IsNil() {
			filteredResult[formType] = formVal.Elem().Interface()
		} else if formVal.Kind() != reflect.Ptr && !formVal.IsZero() {
			filteredResult[formType] = formVal.Interface()
		}
	}
	return filteredResult
}

func ReqReturn(c *gin.Context) {
	var req ImageReq

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	auth := AuthService()
	log.Println(auth)

	filteredReq := FilterEmptyFields(req)
	log.Println(filteredReq)

	c.JSON(http.StatusOK, gin.H{
		"message":          "Success",
		"recieved-request": filteredReq,
	})
}
