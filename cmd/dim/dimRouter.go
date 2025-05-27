package dim

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ConfigureRouter(dimController *DimController) *gin.Engine {

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/flightPlan", dimController.GetFlightPlan)
	router.POST("/flightPlan", dimController.CreateFlightPlan)
	router.PUT("/flightPlan", dimController.UpdateFlightPlan)
	router.DELETE("/flightPlan", dimController.DeleteFlightPlan)
	router.POST("/batch", dimController.UploadBatch)
	//router.GET("/missions", GetMissions)
	//router.GET("/requests", GetRequests)
	//router.GET("/requestsNoObservation", GetRequestsNoObservation)

	return router
}
