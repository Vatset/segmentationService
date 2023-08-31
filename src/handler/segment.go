package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	segmentation_service "segmentationService"
	"segmentationService/src/handler/response"
)

// @Summary Segment Creation
// @Tags segment
// @description Creation of segments
// @description (Before creating a segment, it undergoes the slug validation)
// @description -----------------------------------------------------------------
// @description In addition to the segment name, the user can specify the number of percentages. What percentage of all users should a new segment be automatically added to
// @ID create-segment
// @Accept  json
// @Produce  json
// @Param input body segmentation_service.SegmentPattern true "segment"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Failure default {object} response.errorResponse
// @Router /api/segment/create [post]
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

// @Summary Segment Update
// @Tags segment
// @description Updates segment's name in the database
// @description (Before name changing, new name undergoes the slug validation)
// @ID update-segment
// @Accept  json
// @Produce  json
// @Param input body segmentation_service.UpdateSegment true "last_name , new_name"
// @Success 200 {object} response.StatusResponse
// @Failure 400,404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Failure default {object} response.errorResponse
// @Router /api/segment/update [put]
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

// @Summary Segment Deletion
// @Tags segment
// @Description Delete segment in the database
// @ID delete-segment
// @Accept  json
// @Produce  json
// @Param input body segmentation_service.SegmentPattern true "segment"
// @Success 200 {object} response.StatusResponse
// @Failure 400,404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Failure default {object} response.errorResponse
// @Router /api/segment/delete [delete]
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
