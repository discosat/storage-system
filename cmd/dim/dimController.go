package dim

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discosat/storage-system/internal/Commands"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
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
	oFile, err := file.Open()
	//defer oFile.Close()
	if err != nil {
		errorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	readCloser, ok := oFile.(io.ReadCloser)
	if !ok {
		// Handle the error, file does not implement io.ReadCloser
		return
	}

	observationId, err := d.dimService.handleUploadImage(&readCloser, file.Filename, file.Size)
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

	var flightPlan Commands.FlightPlanCommand
	err := json.Unmarshal([]byte(c.PostForm("flightPlan")), &flightPlan)
	if err != nil {
		errorAbortMessage(c, http.StatusBadRequest, err)
		return
	}

	rList := c.PostFormArray("requestList")
	var orList []Commands.ObservationRequestCommand
	for _, r := range rList {
		var or Commands.ObservationRequestCommand
		err = json.Unmarshal([]byte(r), &or)
		if err != nil {
			slog.Warn(fmt.Sprintf("Could not bind request to ObservationRequuest: %v", err))
			errorAbortMessage(c, http.StatusBadRequest, err)
			return
		}
		orList = append(orList, or)
	}

	slog.Info(fmt.Sprintf("CreateFlightPlan: Request is sucessfully bound, persisting"))

	fpId, err := d.dimService.handleCreateFlightPlan(flightPlan, orList)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not create flight plan: %v, wiht observation requests: %v, %v", flightPlan, orList, err))
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

func (d DimController) UploadBatch(c *gin.Context) {
	file, err := c.FormFile("batch")
	oFile, err := file.Open()

	tmpFile, _ := os.CreateTemp("", "temp*.zip")
	defer os.Remove(tmpFile.Name())
	_, err = io.Copy(tmpFile, oFile)
	if err != nil {
		errorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	reader, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		errorAbortMessage(c, http.StatusBadRequest, err)
		return
	}

	err = d.dimService.handleUploadBatch(reader)
	if err != nil {
		errorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, nil)
	return
}

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
	slog.Error(fmt.Sprint(err))
	c.JSON(statusCode, gin.H{"error": fmt.Sprint(err)})
}
