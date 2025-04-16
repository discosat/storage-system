package dim

import (
	"github.com/gin-gonic/gin"
)

type DimController struct {
	dimService DimServiceInterface
}

func NewDimController(dimService DimServiceInterface) *DimController {
	return &DimController{dimService: dimService}
}

func (d DimController) UploadImage(c *gin.Context) {
	d.dimService.handleUploadImage(c)
	return
}

//func UploadBatch(c *gin.Context) {
//	//handleUploadBatch(c)
//	return
//}
//
//func GetMissions(c *gin.Context) {
//	missions, err := handleGetMissions(c)
//	if err != nil {
//		ErrorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"missions": missions})
//}
//
//func GetRequests(c *gin.Context) {
//	requests, err := handleGetRequests(c)
//	if err != nil {
//		ErrorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"requests": requests})
//}
//
//func GetRequestsNoObservation(c *gin.Context) {
//	requests, err := handleGetRequestsNoObservation(c)
//	if err != nil {
//		ErrorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"requests": requests})
//}
