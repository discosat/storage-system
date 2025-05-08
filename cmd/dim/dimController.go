package dim

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type DimController struct {
	dimService DimServiceInterface
}

func NewDimController(dimService DimServiceInterface) *DimController {
	return &DimController{dimService: dimService}
}

func (d DimController) UploadImage(c *gin.Context) {

	//Binding POST data
	file, err := c.FormFile("file")
	if err != nil {
		errorAbortMessage(c, http.StatusBadRequest, err)
		return
	}

	observationId, err := d.dimService.handleUploadImage(file)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorAbortMessage(c, http.StatusNotFound, err)
			return
		}
		errorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"observation": observationId})
	return
}

func (d DimController) GetFlightPlan(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		errorAbortMessage(c, http.StatusBadRequest, fmt.Errorf("please enter an id"))
		return
	}

	fpId, err := strconv.Atoi(id)
	if err != nil {
		errorAbortMessage(c, http.StatusBadRequest, fmt.Errorf("id is not a number: %v", id))
		return
	}

	flightPLan, err := d.dimService.handleGetFlightPlan(fpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorAbortMessage(c, http.StatusNotFound, fmt.Errorf("no flight plan with id: %v", fpId))
			return
		}
		errorAbortMessage(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"flightPlan": flightPLan})
	return
}

func (d DimController) CreateFlightPlan(c *gin.Context) {
	mId, _ := strconv.Atoi(c.PostForm("missionId"))
	uId, _ := strconv.Atoi(c.PostForm("userId"))
	name := c.PostForm("name")
	rList := c.PostFormArray("requestList")
	var orList []observationRequest.ObservationRequest
	for _, r := range rList {
		var or observationRequest.ObservationRequest
		json.Unmarshal([]byte(r), &or)
		orList = append(orList, or)
	}

	log.Printf("mId: %v, uId: %v, name: %v", mId, uId, name)
	log.Printf("List: %v", rList)
	fpId, err := d.dimService.handleCreateFlightPlan(mId, uId, name, orList)
	if err != nil {
		errorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"flightPlanId": fpId})
}

func (d DimController) Test(c *gin.Context) {
	or, err := d.dimService.test(c)
	if err != nil {
		log.Fatalf("Pis og papir")
	}
	c.JSON(http.StatusCreated, or)
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
//		errorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"missions": missions})
//}
//
//func GetRequests(c *gin.Context) {
//	requests, err := handleGetRequests(c)
//	if err != nil {
//		errorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"requests": requests})
//}
//
//func GetRequestsNoObservation(c *gin.Context) {
//	requests, err := handleGetRequestsNoObservation(c)
//	if err != nil {
//		errorAbortMessage(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"requests": requests})
//}

func errorAbortMessage(c *gin.Context, statusCode int, err error) {
	log.Println(err)
	c.JSON(statusCode, gin.H{"error": fmt.Sprint(err)})
}
