package dim

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadBatch(c *gin.Context) {
	handleUploadBatch(c)
}

func UploadImage(c *gin.Context) {
	handleUploadImage(c)
}

func GetMissions(c *gin.Context) {
	missions, err := handleGetMissions(c)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"missions": missions})
}

func GetRequests(c *gin.Context) {
	requests, err := handleGetRequests(c)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

func GetRequestsNoObservation(c *gin.Context) {
	requests, err := handleGetRequestsNoObservation(c)
	if err != nil {
		ErrorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}
