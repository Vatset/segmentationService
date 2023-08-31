package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "segmentationService/docs"
	"segmentationService/src/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		segment := api.Group("/segment")
		{
			segment.POST("/create", h.createSegment)
			segment.PUT("/update", h.updateSegment)
			segment.DELETE("/delete", h.deleteSegment)

		}
		user := api.Group("/user")
		{
			user.POST("/create", h.createUser)
			user.DELETE("/delete", h.deleteUser)
			user.GET("/showSegments/:id", h.showUserSegments)
			user.GET("/historyLink/:id", h.History)

		}
		segmentation := api.Group("/segmentation")
		{
			segmentation.POST("/", h.segmentMembership)
		}
		router.Static("/api/user/history", "./history")
	}
	return router
}
