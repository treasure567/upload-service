package routes

import (
	"go-upload/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	uploadController := controllers.GetUploadController()

	router.GET("/health", uploadController.HealthCheck)

	api := router.Group("/api")
	{
		upload := api.Group("/upload")
		{
			upload.POST("/", uploadController.UploadFile)
			upload.DELETE("/", uploadController.DeleteFile)
		}
	}
}
