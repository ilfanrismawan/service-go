package delivery

import (
	"log"
	"net/http"
	"service/internal/modules/media/service"
	"service/internal/shared/model"
	"service/internal/shared/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FileHandler handles file upload endpoints
type FileHandler struct {
	fileService *service.FileService
}

// NewFileHandler creates a new file handler
func NewFileHandler() *FileHandler {
	fileService, err := service.NewFileService()
	if err != nil {
		// Log error but return a handler with nil service so the app can continue to run
		// and other endpoints (including Swagger) can be tested. File endpoints will
		// return 503 Service Unavailable when used.
		log.Printf("Failed to initialize file service: %v", err)
		return &FileHandler{fileService: nil}
	}
	return &FileHandler{
		fileService: fileService,
	}
}

// UploadFile godoc
// @Summary Upload file
// @Description Upload a file to storage
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param folder formData string true "Folder name"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /files/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	if h == nil || h.fileService == nil {
		c.JSON(http.StatusServiceUnavailable, model.CreateErrorResponse(
			"file_service_unavailable",
			"File service is not available",
			nil,
		))
		return
	}
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"file_required",
			"File is required",
			nil,
		))
		return
	}

	// Get folder from form
	folder := c.PostForm("folder")
	if folder == "" {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"folder_required",
			"Folder is required",
			nil,
		))
		return
	}

	// Upload file
	response, err := h.fileService.UploadFile(c.Request.Context(), file, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"file_upload_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(response, "File uploaded successfully"))
}

// UploadOrderPhoto godoc
// @Summary Upload order photo
// @Description Upload a photo for an order (pickup, service, or delivery)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Photo file"
// @Param order_id formData string true "Order ID"
// @Param photo_type formData string true "Photo type (pickup, service, delivery)"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /files/orders/photo [post]
func (h *FileHandler) UploadOrderPhoto(c *gin.Context) {
	if h == nil || h.fileService == nil {
		c.JSON(http.StatusServiceUnavailable, model.CreateErrorResponse(
			"file_service_unavailable",
			"File service is not available",
			nil,
		))
		return
	}
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"file_required",
			"File is required",
			nil,
		))
		return
	}

	// Get order ID from form
	orderIDStr := c.PostForm("order_id")
	if orderIDStr == "" {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"order_id_required",
			"Order ID is required",
			nil,
		))
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_order_id",
			"Invalid order ID format",
			nil,
		))
		return
	}

	// Get photo type from form
	photoType := c.PostForm("photo_type")
	if photoType == "" {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"photo_type_required",
			"Photo type is required",
			nil,
		))
		return
	}

	// Upload order photo
	response, err := h.fileService.UploadOrderPhoto(c.Request.Context(), file, orderID, photoType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"photo_upload_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(response, "Order photo uploaded successfully"))
}

// UploadUserAvatar godoc
// @Summary Upload user avatar
// @Description Upload a user avatar
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Avatar file"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /files/users/avatar [post]
func (h *FileHandler) UploadUserAvatar(c *gin.Context) {
	if h == nil || h.fileService == nil {
		c.JSON(http.StatusServiceUnavailable, model.CreateErrorResponse(
			"file_service_unavailable",
			"File service is not available",
			nil,
		))
		return
	}
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.CreateErrorResponse(
			"unauthorized",
			"User ID not found in context",
			nil,
		))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"internal_error",
			"Invalid user ID type",
			nil,
		))
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"file_required",
			"File is required",
			nil,
		))
		return
	}

	// Upload user avatar
	response, err := h.fileService.UploadUserAvatar(c.Request.Context(), file, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"avatar_upload_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(response, "User avatar uploaded successfully"))
}

// GetFileURL godoc
// @Summary Get file URL
// @Description Get a presigned URL for file access
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param object_name query string true "Object name"
// @Param expiry query int false "Expiry in minutes" default(60)
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /files/url [get]
func (h *FileHandler) GetFileURL(c *gin.Context) {
	if h == nil || h.fileService == nil {
		c.JSON(http.StatusServiceUnavailable, model.CreateErrorResponse(
			"file_service_unavailable",
			"File service is not available",
			nil,
		))
		return
	}
	objectName := c.Query("object_name")
	if objectName == "" {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"object_name_required",
			"Object name is required",
			nil,
		))
		return
	}

	// Parse expiry
	expiryStr := c.DefaultQuery("expiry", "60")
	expiryMinutes, err := utils.ParseInt(expiryStr)
	if err != nil || expiryMinutes < 1 || expiryMinutes > 1440 {
		expiryMinutes = 60
	}

	expiry := time.Duration(expiryMinutes) * time.Minute

	// Get file URL
	url, err := h.fileService.GetFileURL(c.Request.Context(), objectName, expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"file_url_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"url": url}, "File URL generated successfully"))
}

// ListFiles godoc
// @Summary List files
// @Description List files in a folder
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param folder query string true "Folder name"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /files/list [get]
func (h *FileHandler) ListFiles(c *gin.Context) {
	if h == nil || h.fileService == nil {
		c.JSON(http.StatusServiceUnavailable, model.CreateErrorResponse(
			"file_service_unavailable",
			"File service is not available",
			nil,
		))
		return
	}
	folder := c.Query("folder")
	if folder == "" {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"folder_required",
			"Folder is required",
			nil,
		))
		return
	}

	// List files
	files, err := h.fileService.ListFiles(c.Request.Context(), folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"list_files_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(files, "Files listed successfully"))
}

// DeleteFile godoc
// @Summary Delete file
// @Description Delete a file from storage
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param object_name query string true "Object name"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /files/delete [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	if h == nil || h.fileService == nil {
		c.JSON(http.StatusServiceUnavailable, model.CreateErrorResponse(
			"file_service_unavailable",
			"File service is not available",
			nil,
		))
		return
	}
	objectName := c.Query("object_name")
	if objectName == "" {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"object_name_required",
			"Object name is required",
			nil,
		))
		return
	}

	// Delete file
	err := h.fileService.DeleteFile(c.Request.Context(), objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"file_delete_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil, "File deleted successfully"))
}
