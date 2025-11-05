package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"service/internal/shared/config"
	"service/internal/shared/model"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// FileService handles file upload and management
type FileService struct {
	minioClient *minio.Client
	bucketName  string
}

// NewFileService creates a new file service
func NewFileService() (*FileService, error) {
	// Parse endpoint URL to extract host and port
	endpoint := config.Config.S3Endpoint
	if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.TrimPrefix(endpoint, "http://")
	} else if strings.HasPrefix(endpoint, "https://") {
		endpoint = strings.TrimPrefix(endpoint, "https://")
	}

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Config.S3AccessKey, config.Config.S3SecretKey, ""),
		Secure: false, // Set to true for HTTPS
	})
	if err != nil {
		return nil, err
	}

	// Ensure bucket exists
	bucketName := config.Config.S3BucketName
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &FileService{
		minioClient: minioClient,
		bucketName:  bucketName,
	}, nil
}

// UploadFile uploads a file to S3-compatible storage
func (s *FileService) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*model.FileUploadResponse, error) {
	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", folder, uuid.New().String(), ext)
	objectName := fmt.Sprintf("%s/%s", folder, filename)

	// Upload file
	_, err = s.minioClient.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return nil, err
	}

	// Generate file URL
	fileURL := fmt.Sprintf("%s/%s/%s", config.Config.S3Endpoint, s.bucketName, objectName)

	response := &model.FileUploadResponse{
		Filename:     filename,
		OriginalName: file.Filename,
		Size:         file.Size,
		ContentType:  file.Header.Get("Content-Type"),
		URL:          fileURL,
		UploadedAt:   time.Now(),
	}

	return response, nil
}

// UploadOrderPhoto uploads a photo for an order
func (s *FileService) UploadOrderPhoto(ctx context.Context, file *multipart.FileHeader, orderID uuid.UUID, photoType string) (*model.FileUploadResponse, error) {
	// Validate photo type
	validTypes := []string{"pickup", "service", "delivery"}
	if !contains(validTypes, photoType) {
		return nil, fmt.Errorf("invalid photo type: %s", photoType)
	}

	// Validate file type
	if !isImageFile(file.Filename) {
		return nil, fmt.Errorf("file must be an image")
	}

	// Upload file
	folder := fmt.Sprintf("orders/%s/%s", orderID.String(), photoType)
	return s.UploadFile(ctx, file, folder)
}

// UploadUserAvatar uploads a user avatar
func (s *FileService) UploadUserAvatar(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID) (*model.FileUploadResponse, error) {
	// Validate file type
	if !isImageFile(file.Filename) {
		return nil, fmt.Errorf("file must be an image")
	}

	// Upload file
	folder := fmt.Sprintf("users/%s/avatar", userID.String())
	return s.UploadFile(ctx, file, folder)
}

// DeleteFile deletes a file from storage
func (s *FileService) DeleteFile(ctx context.Context, objectName string) error {
	return s.minioClient.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
}

// GetFileURL generates a presigned URL for file access
func (s *FileService) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := s.minioClient.PresignedGetObject(ctx, s.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// ListFiles lists files in a folder
func (s *FileService) ListFiles(ctx context.Context, folder string) ([]model.FileInfo, error) {
	objectCh := s.minioClient.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix:    folder + "/",
		Recursive: true,
	})

	var files []model.FileInfo
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		files = append(files, model.FileInfo{
			Name:         object.Key,
			Size:         object.Size,
			LastModified: object.LastModified,
			ContentType:  object.ContentType,
		})
	}

	return files, nil
}

// isImageFile checks if file is an image
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	return contains(imageExts, ext)
}

// contains checks if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
