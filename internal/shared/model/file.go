<<<<<<< HEAD
package core
=======
package model
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b

import (
	"time"
)

// FileUploadResponse represents the response for file upload
type FileUploadResponse struct {
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	URL          string    `json:"url"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

// FileInfo represents file information
type FileInfo struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ContentType  string    `json:"content_type"`
}

// FileUploadRequest represents the request for file upload
type FileUploadRequest struct {
	Folder string `json:"folder" validate:"required"`
}

// OrderPhotoUploadRequest represents the request for order photo upload
type OrderPhotoUploadRequest struct {
<<<<<<< HEAD
	OrderID    string `json:"order_id" validate:"required"`
	PhotoType  string `json:"photo_type" validate:"required,oneof=pickup service delivery"`
=======
	OrderID   string `json:"order_id" validate:"required"`
	PhotoType string `json:"photo_type" validate:"required,oneof=pickup service delivery"`
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
}

// UserAvatarUploadRequest represents the request for user avatar upload
type UserAvatarUploadRequest struct {
	UserID string `json:"user_id" validate:"required"`
}
