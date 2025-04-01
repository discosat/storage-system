package dim

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"os"
)

type Mission struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Bucket string `json:"bucket"`
}

type Request struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	MissionId int    `json:"mission_id"`
	UserId    int    `json:"user_id"`
}

type MeasurementRequest struct {
	Id    int    `json:"id"`
	MType string `json:"m_type"`
	RId   int    `json:"r_id"`
}

type Observation struct {
	Id        int `json:"id"`
	RequestId int `json:"request_id"`
	UserId    int `json:"user_id"`
}

var ObjectStore SimpleStore
var db *sql.DB

func Start() {
	ObjectStore = NewSimpleStore(NewMinioStore())
	db = ConfigDatabase()
	defer db.Close()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/file", UploadImage)
	router.POST("/batch", UploadBatch)
	router.GET("/missions", GetMissions)
	router.GET("/requests", GetRequests)
	router.GET("/requestsNoObservation", GetRequestsNoObservation)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func ConfigDatabase() *sql.DB {
	db, err := sql.Open("pgx", fmt.Sprint("postgres://", os.Getenv("PGUSER"), ":", os.Getenv("PGPASSWORD"), "@", os.Getenv("PGHOST"), "/", os.Getenv("PGDATABASE")))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	return db
}

func ErrorAbortMessage(c *gin.Context, statusCode int, err error) {
	log.Println(err)
	c.JSON(statusCode, gin.H{"error": fmt.Sprint(err)})
}
