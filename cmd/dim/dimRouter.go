package dim

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
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
	//router.POST("/batch", UploadBatch)
	//router.GET("/missions", GetMissions)
	//router.GET("/requests", GetRequests)
	//router.GET("/requestsNoObservation", GetRequestsNoObservation)
	router.GET("/test", dimController.Test)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func ErrorAbortMessage(c *gin.Context, statusCode int, err error) {
	log.Println(err)
	c.JSON(statusCode, gin.H{"error": fmt.Sprint(err)})
}
