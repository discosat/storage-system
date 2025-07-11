package dim

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
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

	flightPlan, err := d.dimService.handleGetFlightPlan(fpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorAbortMessage(c, http.StatusNotFound, fmt.Errorf("no flight plan with id: %v", fpId))
			return
		}
		errorAbortMessage(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"flightPlan": flightPlan})
}

func (d DimController) CreateFlightPlan(c *gin.Context) {

	var flightPlan observationRequest.FlightPlanAggregate
	err := json.Unmarshal([]byte(c.PostForm("flightPlan")), &flightPlan)
	if err != nil {
		errorAbortMessage(c, http.StatusBadRequest, err)
		return
	}

	slog.Info("CreateFlightPlan: Request is sucessfully bound, persisting")

	fpId, err := d.dimService.handleCreateFlightPlan(flightPlan)
	if err != nil {
		if err.(*observationRequest.ObservationRequestError).Code() == observationRequest.ObservationRequestParseError {
			slog.Error(fmt.Sprintf("One or more observation requests are formatted wrong: %v", flightPlan.ObservationRequests))
			errorAbortMessage(c, http.StatusBadRequest, err)
			return
		}
		slog.Error(fmt.Sprintf("Could not create flight plan: %v, wiht observation requests: %v, %v", flightPlan, flightPlan.ObservationRequests, err))
		errorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"flightPlanId": fpId})
}

func (d DimController) UpdateFlightPlan(c *gin.Context) {
	// TODO Check permissions!!!!!!!!
	var flightPlan observationRequest.FlightPlanAggregate
	err := json.Unmarshal([]byte(c.PostForm("flightPlan")), &flightPlan)
	if err != nil {
		errorAbortMessage(c, http.StatusBadRequest, err)
		return
	}
	id, err := d.dimService.handleUpdateFlightPlan(flightPlan)
	if err != nil {
		if lockedErr, ok := err.(*observationRequest.ObservationRequestError); ok && lockedErr.Code() == observationRequest.FlightPlanIsLocked {
			errorAbortMessage(c, http.StatusBadRequest, err)
			return
		}
		errorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"flightPlanId": id})
}

func (d DimController) DeleteFlightPlan(c *gin.Context) {
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

	_, err = d.dimService.handleDeleteFlightPlan(fpId)
	if err != nil {
		if er, ok := err.(*observationRequest.ObservationRequestError); ok && er.Code() == observationRequest.FlightPlanNotFound {
			errorAbortMessage(c, http.StatusBadRequest, err)
			return
		}
		if pgEr, ok := err.(*pgconn.PgError); ok && pgEr.Code == "23503" {
			errorAbortMessage(c, http.StatusBadRequest, pgEr)
			return
		}
		errorAbortMessage(c, http.StatusInternalServerError, fmt.Errorf("error in deleting flight plan wtih id: %v: %v", fpId, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("flight plan: %v has been deleted", fpId)})
}

func (d DimController) UploadBatch(c *gin.Context) {
	file, err := c.FormFile("batch")
	if err != nil {
		errorAbortMessage(c, http.StatusBadRequest, err)
	}

	oFile, err := file.Open()
	if err != nil {
		errorAbortMessage(c, http.StatusBadRequest, err)
	}

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

	ids, err := d.dimService.handleUploadBatch(reader)
	if err != nil {
		errorAbortMessage(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ObservationIds": ids})
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
	slog.Error(fmt.Sprint(err.Error()))
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(statusCode, gin.H{"error": fmt.Sprint(err.Error())})
}
