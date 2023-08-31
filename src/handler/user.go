package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	segmentation_service "segmentationService"
	"segmentationService/src/handler/response"
	"strconv"
	"time"
)

// @Summary User Creation
// @Tags user
// @Description Create a user and return user_id
// @ID create-user
// @Accept  json
// @Produce  json
// @Param input body segmentation_service.User true "username"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Failure default {object} response.errorResponse
// @Router /api/user/create [post]
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

// @Summary User Deletion
// @Tags user
// @Description Delete user from database
// @ID delete-user
// @Accept  json
// @Produce  json
// @Param input body segmentation_service.User true "username"
// @Success 200 {object} response.StatusResponse
// @Failure 400,404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Failure default {object} response.errorResponse
// @Router /api/user/delete [delete]
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

// @Summary Show User Segments
// @Tags user
// @Description Show segments by user_id
// @ID show-segments
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} response.StatusResponse
// @Failure 400,404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Failure default {object} response.errorResponse
// @Router /api/user/showSegments/{id} [get]
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

// @Summary Show User History
// @Tags user
// @Description Returns link to csv file of users history
// @ID user-history
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param input body segmentation_service.ShowHistory true "year-month"
// @Success 200 {string} string "URL to download the CSV file"
// @Failure 400,404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Failure default {object} response.errorResponse
// @Router /api/user/historyLink/{id} [get]
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
