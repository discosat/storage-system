package dam

import (
	"github.com/discosat/storage-system/cmd/interfaces"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
)

func Start() {
	g := gin.Default()

	g.GET("/search-images", RequestHandler)

	err := g.Run(":8081")
	if err != nil {
		log.Fatal("Failed to start server")
	}
}

func RequestHandler(c *gin.Context) {
	var req interfaces.ImageRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	//Authenticating the request
	auth := AuthService()
	log.Println(auth)

	//Passing request on to QOM
	//insert function to pass on query here

	//bundling images together
	imageFound := ImageBundler()

	//Handle response
	ResponseHandler(c, imageFound)
}

func ResponseHandler(c *gin.Context, imageFound bool) {
	if imageFound {
		imagePath := filepath.Join("cmd", "mock_images", "Greenland_Glacier_Mock.jpg")
		c.File(imagePath)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "No images matched the query",
	})
}
