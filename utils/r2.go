package utils

import (
	"bytes"
	"fmt"
	"go-upload/config"
	"go-upload/models"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

func UploadToR2(fileBuffer []byte, folderPath string, ext string, brandingOptions *BrandingOptions) (*models.UploadResult, error) {
	if brandingOptions != nil && brandingOptions.Position != "" && brandingOptions.BrandLogo != "" && brandingOptions.Width > 0 && brandingOptions.Height > 0 {
		allowedExtensions := []string{"png", "jpg", "jpeg"}
		extLower := strings.ToLower(ext)
		isAllowed := false
		for _, allowed := range allowedExtensions {
			if extLower == allowed {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			fmt.Println("Applying branding before upload...")
			brandedBuffer, err := BrandImage(fileBuffer, *brandingOptions, ext)
			if err != nil {
				return nil, fmt.Errorf("failed to apply branding: %w", err)
			}
			fileBuffer = brandedBuffer
		} else {
			fmt.Printf("Branding skipped: File type .%s not supported. Only PNG, JPG, JPEG are supported.\n", ext)
		}
	}

	return uploadToCloudflareR2(fileBuffer, folderPath, ext)
}

func uploadToCloudflareR2(fileBuffer []byte, folderPath string, ext string) (*models.UploadResult, error) {
	cfg := config.AppConfig.R2

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	svc := s3.New(sess)

	filename := fmt.Sprintf("%s.%s", uuid.New().String(), ext)
	objectKey := filepath.Join(folderPath, filename)

	contentType := getContentType(ext)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(cfg.BucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(fileBuffer),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to R2: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s", cfg.CustomDomain, objectKey)

	return &models.UploadResult{
		URL:      publicURL,
		PublicID: objectKey,
	}, nil
}

func DeleteFromR2(key string) error {
	cfg := config.AppConfig.R2

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Endpoint:    aws.String(cfg.Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(cfg.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from R2: %w", err)
	}

	return nil
}

func getContentType(ext string) string {
	contentTypes := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"pdf":  "application/pdf",
		"mp4":  "video/mp4",
		"mov":  "video/quicktime",
		"zip":  "application/zip",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}
