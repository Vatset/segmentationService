package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	segmentation_service "segmentationService"
	"segmentationService/src/handler/response"
	"strconv"
	"time"
)

func (h *Handler) createUser(c *gin.Context) {
	var input segmentation_service.User

	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid username field")
		return
	}
	id, err := h.service.User.CreateUser(input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) deleteUser(c *gin.Context) {
	var input segmentation_service.User

	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid username field")
		return
	}
	err := h.service.User.DeleteUser(input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusResponse{
		Status: "user was successful deleted",
	})
}

func (h *Handler) showUserSegments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	segments, err := h.service.ShowUserSegments(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"User Segments": segments,
	})
}

func (h *Handler) History(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input segmentation_service.ShowHistory

	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	layout := "2006-01"
	parsedDate, err := time.Parse(layout, input.Timestamp)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Error parsing date")
		return
	}
	startDate := parsedDate.Format("2006-01-02")
	endDate := parsedDate.AddDate(0, 1, 0).Format("2006-01-02")

	link, err := h.service.CreateLinkToCSV(id, input.Timestamp, startDate, endDate)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Link": link,
	})
}
