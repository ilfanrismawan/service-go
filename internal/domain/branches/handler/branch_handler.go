package handler

import (
	"net/http"
	"service/internal/domain/branches/service"
	"service/internal/shared/model"
	"service/internal/shared/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BranchHandler handles branch endpoints
type BranchHandler struct {
	branchService *service.BranchService
}

// NewBranchHandler creates a new branch handler
func NewBranchHandler() *BranchHandler {
	return &BranchHandler{
		branchService: service.NewBranchService(),
	}
}

// CreateBranch godoc
// @Summary Create a new branch
// @Description Create a new branch/outlet
// @Tags branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.BranchRequest true "Branch data"
// @Success 201 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /branches [post]
func (h *BranchHandler) CreateBranch(c *gin.Context) {
	var req model.BranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	branch, err := h.branchService.CreateBranch(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"branch_creation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(branch, "Branch created successfully"))
}

// GetBranch godoc
// @Summary Get branch by ID
// @Description Get branch details by ID
// @Tags branches
// @Accept json
// @Produce json
// @Param id path string true "Branch ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /branches/{id} [get]
func (h *BranchHandler) GetBranch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_id",
			"Invalid branch ID format",
			nil,
		))
		return
	}

	branch, err := h.branchService.GetBranch(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == model.ErrBranchNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.CreateErrorResponse(
			"branch_not_found",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(branch, "Branch retrieved successfully"))
}

// UpdateBranch godoc
// @Summary Update branch
// @Description Update branch information
// @Tags branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID"
// @Param request body model.BranchRequest true "Branch update data"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /branches/{id} [put]
func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_id",
			"Invalid branch ID format",
			nil,
		))
		return
	}

	var req model.BranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"validation_error",
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"validation_error",
			"Validation failed",
			err.Error(),
		))
		return
	}

	branch, err := h.branchService.UpdateBranch(c.Request.Context(), id, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == model.ErrBranchNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.CreateErrorResponse(
			"branch_update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(branch, "Branch updated successfully"))
}

// DeleteBranch godoc
// @Summary Delete branch
// @Description Soft delete a branch
// @Tags branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /branches/{id} [delete]
func (h *BranchHandler) DeleteBranch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_id",
			"Invalid branch ID format",
			nil,
		))
		return
	}

	err = h.branchService.DeleteBranch(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == model.ErrBranchNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.CreateErrorResponse(
			"branch_deletion_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil, "Branch deleted successfully"))
}

// ListBranches godoc
// @Summary List branches
// @Description Get list of branches with pagination and filters
// @Tags branches
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param city query string false "Filter by city"
// @Param province query string false "Filter by province"
// @Success 200 {object} model.PaginatedResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /branches [get]
func (h *BranchHandler) ListBranches(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	city := c.Query("city")
	province := c.Query("province")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Set filter pointers
	var cityPtr, provincePtr *string
	if city != "" {
		cityPtr = &city
	}
	if province != "" {
		provincePtr = &province
	}

	result, err := h.branchService.ListBranches(c.Request.Context(), page, limit, cityPtr, provincePtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"branch_list_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetNearbyBranches godoc
// @Summary Get nearby branches
// @Description Get branches within a certain radius from given coordinates
// @Tags branches
// @Accept json
// @Produce json
// @Param latitude query number true "Latitude"
// @Param longitude query number true "Longitude"
// @Param radius query number false "Radius in kilometers" default(10)
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /branches/nearby [get]
func (h *BranchHandler) GetNearbyBranches(c *gin.Context) {
	// Parse query parameters
	latStr := c.Query("latitude")
	lonStr := c.Query("longitude")
	radiusStr := c.DefaultQuery("radius", "10")

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_latitude",
			"Invalid latitude value",
			nil,
		))
		return
	}

	longitude, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_longitude",
			"Invalid longitude value",
			nil,
		))
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radius <= 0 {
		radius = 10
	}

	branches, err := h.branchService.GetNearbyBranches(c.Request.Context(), latitude, longitude, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"nearby_branches_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(branches, "Nearby branches retrieved successfully"))
}

// GetActiveBranches godoc
// @Summary Get active branches
// @Description Get all active branches
// @Tags branches
// @Accept json
// @Produce json
// @Success 200 {object} model.APIResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /branches/active [get]
func (h *BranchHandler) GetActiveBranches(c *gin.Context) {
	branches, err := h.branchService.GetActiveBranches(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"active_branches_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(branches, "Active branches retrieved successfully"))
}

// GetBranches godoc
// @Summary Get all branches
// @Description Get all branches with pagination
// @Tags branches
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} model.APIResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /branches [get]
func (h *BranchHandler) GetBranches(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	branches, total, err := h.branchService.GetBranches(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"branches_failed",
			err.Error(),
			nil,
		))
		return
	}

	// Create paginated response
	pagination := model.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	c.JSON(http.StatusOK, model.PaginatedSuccessResponse(branches, pagination, "Branches retrieved successfully"))
}

// GetNearestBranches godoc
// @Summary Get nearest branches
// @Description Get nearest branches based on location
// @Tags branches
// @Accept json
// @Produce json
// @Param lat query float64 true "Latitude"
// @Param lon query float64 true "Longitude"
// @Param radius query float64 false "Search radius in km" default(10)
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /branches/nearest [get]
func (h *BranchHandler) GetNearestBranches(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	radiusStr := c.DefaultQuery("radius", "10")

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"missing_coordinates",
			"Latitude and longitude are required",
			nil,
		))
		return
	}

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_latitude",
			"Invalid latitude value",
			nil,
		))
		return
	}

	longitude, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(
			"invalid_longitude",
			"Invalid longitude value",
			nil,
		))
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radius <= 0 {
		radius = 10
	}

	branches, err := h.branchService.GetNearbyBranches(c.Request.Context(), latitude, longitude, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(
			"nearest_branches_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(branches, "Nearest branches retrieved successfully"))
}
