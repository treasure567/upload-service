package controllers

import (
	"io"
	"net/http"
	"strconv"

	"go-upload/models"
	"go-upload/services"
	"go-upload/utils"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	uploadService *services.UploadService
}

var uploadControllerInstance *UploadController

func GetUploadController() *UploadController {
	if uploadControllerInstance == nil {
		uploadControllerInstance = &UploadController{
			uploadService: services.GetUploadService(),
		}
	}
	return uploadControllerInstance
}

func (ctrl *UploadController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "No file provided",
			Error:   err.Error(),
		})
		return
	}

	path := c.PostForm("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Path is required",
		})
		return
	}

	fileHandle, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to open file",
			Error:   err.Error(),
		})
		return
	}
	defer fileHandle.Close()

	fileBuffer, err := io.ReadAll(fileHandle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to read file",
			Error:   err.Error(),
		})
		return
	}

	position := c.PostForm("position")
	brandLogo := c.PostForm("brandLogo")
	widthStr := c.PostForm("width")
	heightStr := c.PostForm("height")

	var brandingOptions *utils.BrandingOptions

	if position != "" && brandLogo != "" && widthStr != "" && heightStr != "" {
		width, err := strconv.Atoi(widthStr)
		if err == nil {
			height, err := strconv.Atoi(heightStr)
			if err == nil {
				brandingOptions = &utils.BrandingOptions{
					Position:  utils.BrandPosition(position),
					BrandLogo: brandLogo,
					Width:     width,
					Height:    height,
				}
			}
		}
	}

	result, err := ctrl.uploadService.UploadFile(fileBuffer, file.Filename, path, brandingOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to upload file",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "File uploaded successfully",
		FileUrl: result.URL,
	})
}

func (ctrl *UploadController) DeleteFile(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "URL is required",
			Error:   err.Error(),
		})
		return
	}

	err := ctrl.uploadService.DeleteFile(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete file",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "File deleted successfully",
	})
}

func (ctrl *UploadController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Server is healthy",
	})
}
