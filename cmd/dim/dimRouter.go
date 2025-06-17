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

	router.GET("/flightPlan", dimController.GetFlightPlan)
	router.POST("/flightPlan", dimController.CreateFlightPlan)
	router.PUT("/flightPlan", dimController.UpdateFlightPlan)
	router.DELETE("/flightPlan", dimController.DeleteFlightPlan)
	router.POST("/batch", dimController.UploadBatch)
	//router.GET("/missions", GetMissions)
	//router.GET("/requests", GetRequests)
	//router.GET("/requestsNoObservation", GetRequestsNoObservation)
	router.Use(CORSMiddleware())

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
