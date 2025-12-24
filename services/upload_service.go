package services

import (
	"fmt"
	"go-upload/models"
	"go-upload/utils"
	"path/filepath"
	"strings"
)

type UploadService struct{}

var uploadServiceInstance *UploadService

func GetUploadService() *UploadService {
	if uploadServiceInstance == nil {
		uploadServiceInstance = &UploadService{}
	}
	return uploadServiceInstance
}

func (s *UploadService) UploadFile(fileBuffer []byte, filename string, folderPath string) (*models.UploadResult, error) {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	if ext == "" {
		ext = "bin"
	}

	result, err := utils.UploadToR2(fileBuffer, folderPath, ext)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return result, nil
}

func (s *UploadService) DeleteFile(key string) error {
	err := utils.DeleteFromR2(key)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
