package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	segmentation_service "segmentationService"
	"segmentationService/src/handler/response"
)

// @Summary Segmentation
// @Tags segmentation
// @description Add&Delete user to/from segment
// @ID segment-membership
// @Accept  json
// @Produce  json
// @Param input body segmentation_service.Segmentation true "Segmentation"
// @Success 200 {object} response.StatusResponse
// @Failure 400,404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Failure default {object} response.errorResponse
// @Router /api/segmentation/ [post]
func (h *Handler) segmentMembership(c *gin.Context) {
	var input segmentation_service.Segmentation

	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.SegmentChecker(input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.service.SegmentMembership(input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusResponse{
		Status: "Segmentation was successful",
	})
}
