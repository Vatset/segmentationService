package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	segmentation_service "segmentationService"
	"segmentationService/src/handler/response"
)

func (h *Handler) createSegment(c *gin.Context) {
	var segmentCreation segmentation_service.SegmentPattern
	var segmentation segmentation_service.Segmentation

	if err := c.BindJSON(&segmentCreation); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := segmentation_service.SlugValidation(segmentCreation.Segment)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Segment's name doesn't match slug pattern")
		return
	}

	id, err := h.service.Segment.CreateSegment(segmentCreation)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if segmentCreation.Percent != 0 {
		h.service.Segmentation.AutoSegmentation(segmentCreation.Percent, segmentCreation.Segment, segmentation)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) updateSegment(c *gin.Context) {
	var input segmentation_service.UpdateSegment

	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := segmentation_service.SlugValidation(input.NewName); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "New name of segment does not match slug pattern")
		return
	}

	err := h.service.UpdateSegment(input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusResponse{
		Status: "segment was successful updated",
	})
}

func (h *Handler) deleteSegment(c *gin.Context) {
	var segmentDeletion segmentation_service.SegmentPattern
	var segmentationDeletion segmentation_service.Segmentation
	if err := c.BindJSON(&segmentDeletion); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	segmentationDeletion.DeletingSegments = []string{segmentDeletion.Segment}
	if err := h.service.SegmentChecker(segmentationDeletion); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.service.DeleteSegment(segmentDeletion, segmentationDeletion)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusResponse{
		Status: "segment was successful deleted",
	})
}
