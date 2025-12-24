package models

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	FileUrl string      `json:"fileUrl,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type UploadResult struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
}

type FileUploadResponse struct {
	FileURL string `json:"fileUrl"`
}
