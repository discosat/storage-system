package dim

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func ConfigureRouter(dimController *DimController) *gin.Engine {

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	//router.POST("/file", dimController.UploadImage)
	router.GET("/flightPlan", dimController.GetFlightPlan)
	router.POST("/flightPlan", dimController.CreateFlightPlan)
	router.PUT("/flightPlan", dimController.UpdateFlightPlan)
	router.POST("/batch", dimController.UploadBatch)
	//router.GET("/missions", GetMissions)
	//router.GET("/requests", GetRequests)
	//router.GET("/requestsNoObservation", GetRequestsNoObservation)

	return router
}
