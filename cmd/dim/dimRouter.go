package dim

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func Start(dimController *DimController) {

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/file", dimController.UploadImage)
	router.GET("/flightPlan", dimController.GetFlightPlan)
	router.POST("/flightPlan", dimController.CreateFlightPlan)
	//router.POST("/batch", UploadBatch)
	//router.GET("/missions", GetMissions)
	//router.GET("/requests", GetRequests)
	//router.GET("/requestsNoObservation", GetRequestsNoObservation)
	router.GET("/test", dimController.Test)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
