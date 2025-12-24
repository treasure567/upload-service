package routes

import (
	"go-upload/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	uploadController := controllers.GetUploadController()

	router.GET("/health", uploadController.HealthCheck)
	router.POST("/upload", uploadController.UploadFile)
	router.POST("/delete", uploadController.DeleteFile)

	// api := router.Group("/api")
	// {
	// 	upload := api.Group("/upload")
	// 	{
	// 		upload.POST("", uploadController.UploadFile)
	// 		upload.POST("/delete", uploadController.DeleteFile)
	// 	}
	// }
}
