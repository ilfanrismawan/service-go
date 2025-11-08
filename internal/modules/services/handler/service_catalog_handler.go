package handler

import (
	"net/http"
	"service/internal/modules/services/repository"
	"service/internal/modules/services/service"
	"service/internal/shared/middleware"
	"service/internal/shared/model"
	"service/internal/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ServiceCatalogHandler handles service catalog HTTP requests
type ServiceCatalogHandler struct {
	categoryService *service.ServiceCategoryService
	catalogService  *service.ServiceCatalogService
	providerService *service.ServiceProviderService
}

// NewServiceCatalogHandler creates a new service catalog handler
func NewServiceCatalogHandler() *ServiceCatalogHandler {
	return &ServiceCatalogHandler{
		categoryService: service.NewServiceCategoryService(),
		catalogService:  service.NewServiceCatalogService(),
		providerService: service.NewServiceProviderService(),
	}
}

// RegisterRoutes registers service catalog routes
// This method registers public routes only. Protected routes should be registered separately in router.go
func (h *ServiceCatalogHandler) RegisterRoutes(r *gin.RouterGroup) {
	services := r.Group("/services")
	{
		// Public routes
		services.GET("/categories", h.GetAllCategories)
		services.GET("/categories/:id", h.GetCategory)
		services.GET("/catalogs", h.GetAllCatalogs)
		services.GET("/catalogs/:id", h.GetCatalog)
		services.GET("/catalogs/category/:categoryId", h.GetCatalogsByCategory)
		services.GET("/providers", h.GetAllProviders)
		services.GET("/providers/:id", h.GetProvider)
	}
}

// RegisterProtectedRoutes registers protected service catalog routes
func (h *ServiceCatalogHandler) RegisterProtectedRoutes(r *gin.RouterGroup) {
	services := r.Group("/services")
	{
		// Category management (Admin only)
		services.POST("/categories", h.CreateCategory)
		services.PUT("/categories/:id", h.UpdateCategory)
		services.DELETE("/categories/:id", h.DeleteCategory)

		// Catalog management (Admin only)
		services.POST("/catalogs", h.CreateCatalog)
		services.PUT("/catalogs/:id", h.UpdateCatalog)
		services.DELETE("/catalogs/:id", h.DeleteCatalog)

		// Provider management
		services.POST("/providers", h.CreateProvider)
		services.PUT("/providers/:id", h.UpdateProvider)
		services.DELETE("/providers/:id", h.DeleteProvider)
		services.POST("/providers/:id/services/:serviceId", h.AddServiceToProvider)
		services.DELETE("/providers/:id/services/:serviceId", h.RemoveServiceFromProvider)
	}
}

// CreateCategory creates a new service category
// @Summary Create service category
// @Description Create a new service category (Admin only)
// @Tags Services
// @Accept json
// @Produce json
// @Param request body model.ServiceCategoryRequest true "Category request"
// @Success 201 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Router /services/categories [post]
func (h *ServiceCatalogHandler) CreateCategory(c *gin.Context) {
	var req model.ServiceCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(category, "Category created successfully"))
}

// GetCategory retrieves a category by ID
// @Summary Get service category
// @Description Get a service category by ID
// @Tags Services
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} model.APIResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /services/categories/{id} [get]
func (h *ServiceCatalogHandler) GetCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	category, err := h.categoryService.GetCategory(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(category, "Category retrieved successfully"))
}

// GetAllCategories retrieves all categories
// @Summary Get all service categories
// @Description Get all service categories
// @Tags Services
// @Produce json
// @Param include_inactive query bool false "Include inactive categories"
// @Success 200 {object} model.APIResponse
// @Router /services/categories [get]
func (h *ServiceCatalogHandler) GetAllCategories(c *gin.Context) {
	includeInactive := c.Query("include_inactive") == "true"

	categories, err := h.categoryService.GetAllCategories(c.Request.Context(), includeInactive)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(categories, "Categories retrieved successfully"))
}

// UpdateCategory updates a category
// @Summary Update service category
// @Description Update a service category (Admin only)
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body model.ServiceCategoryRequest true "Category request"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /services/categories/{id} [put]
func (h *ServiceCatalogHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	var req model.ServiceCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(category, "Category updated successfully"))
}

// DeleteCategory deletes a category
// @Summary Delete service category
// @Description Delete a service category (Admin only)
// @Tags Services
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} model.APIResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /services/categories/{id} [delete]
func (h *ServiceCatalogHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil, "Category deleted successfully"))
}

// CreateCatalog creates a new service catalog
// @Summary Create service catalog
// @Description Create a new service catalog (Admin only)
// @Tags Services
// @Accept json
// @Produce json
// @Param request body model.ServiceCatalogRequest true "Catalog request"
// @Success 201 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /services/catalogs [post]
func (h *ServiceCatalogHandler) CreateCatalog(c *gin.Context) {
	var req model.ServiceCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	catalog, err := h.catalogService.CreateCatalog(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(catalog, "Catalog created successfully"))
}

// GetCatalog retrieves a catalog by ID
// @Summary Get service catalog
// @Description Get a service catalog by ID
// @Tags Services
// @Produce json
// @Param id path string true "Catalog ID"
// @Success 200 {object} model.APIResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /services/catalogs/{id} [get]
func (h *ServiceCatalogHandler) GetCatalog(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	catalog, err := h.catalogService.GetCatalog(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(catalog, "Catalog retrieved successfully"))
}

// GetAllCatalogs retrieves all catalogs
// @Summary Get all service catalogs
// @Description Get all service catalogs with optional filters
// @Tags Services
// @Produce json
// @Param category_id query string false "Filter by category ID"
// @Param is_active query bool false "Filter by active status"
// @Param requires_appointment query bool false "Filter by requires appointment"
// @Param search query string false "Search by name or description"
// @Success 200 {object} model.APIResponse
// @Router /services/catalogs [get]
func (h *ServiceCatalogHandler) GetAllCatalogs(c *gin.Context) {
	filters := &repository.ServiceCatalogFilters{}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			filters.CategoryID = &categoryID
		}
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filters.IsActive = &isActive
	}

	if requiresAppointmentStr := c.Query("requires_appointment"); requiresAppointmentStr != "" {
		requiresAppointment := requiresAppointmentStr == "true"
		filters.RequiresAppointment = &requiresAppointment
	}

	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	catalogs, err := h.catalogService.GetAllCatalogs(c.Request.Context(), filters)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(catalogs, "Catalogs retrieved successfully"))
}

// GetCatalogsByCategory retrieves all catalogs in a category
// @Summary Get catalogs by category
// @Description Get all service catalogs in a specific category
// @Tags Services
// @Produce json
// @Param categoryId path string true "Category ID"
// @Param include_inactive query bool false "Include inactive catalogs"
// @Success 200 {object} model.APIResponse
// @Router /services/catalogs/category/{categoryId} [get]
func (h *ServiceCatalogHandler) GetCatalogsByCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	includeInactive := c.Query("include_inactive") == "true"

	catalogs, err := h.catalogService.GetCatalogsByCategory(c.Request.Context(), categoryID, includeInactive)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(catalogs, "Catalogs retrieved successfully"))
}

// UpdateCatalog updates a catalog
// @Summary Update service catalog
// @Description Update a service catalog (Admin only)
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "Catalog ID"
// @Param request body model.ServiceCatalogRequest true "Catalog request"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /services/catalogs/{id} [put]
func (h *ServiceCatalogHandler) UpdateCatalog(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	var req model.ServiceCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	catalog, err := h.catalogService.UpdateCatalog(c.Request.Context(), id, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(catalog, "Catalog updated successfully"))
}

// DeleteCatalog deletes a catalog
// @Summary Delete service catalog
// @Description Delete a service catalog (Admin only)
// @Tags Services
// @Produce json
// @Param id path string true "Catalog ID"
// @Success 200 {object} model.APIResponse
// @Router /services/catalogs/{id} [delete]
func (h *ServiceCatalogHandler) DeleteCatalog(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	if err := h.catalogService.DeleteCatalog(c.Request.Context(), id); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil, "Catalog deleted successfully"))
}

// CreateProvider creates a new service provider
// @Summary Create service provider
// @Description Create a new service provider
// @Tags Services
// @Accept json
// @Produce json
// @Param request body model.ServiceProviderRequest true "Provider request"
// @Success 201 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /services/providers [post]
func (h *ServiceCatalogHandler) CreateProvider(c *gin.Context) {
	var req model.ServiceProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	provider, err := h.providerService.CreateProvider(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(provider, "Provider created successfully"))
}

// GetProvider retrieves a provider by ID
// @Summary Get service provider
// @Description Get a service provider by ID
// @Tags Services
// @Produce json
// @Param id path string true "Provider ID"
// @Success 200 {object} model.APIResponse
// @Router /services/providers/{id} [get]
func (h *ServiceCatalogHandler) GetProvider(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	provider, err := h.providerService.GetProvider(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(provider, "Provider retrieved successfully"))
}

// GetAllProviders retrieves all providers
// @Summary Get all service providers
// @Description Get all service providers with optional filters
// @Tags Services
// @Produce json
// @Param city query string false "Filter by city"
// @Param province query string false "Filter by province"
// @Param is_active query bool false "Filter by active status"
// @Param is_verified query bool false "Filter by verified status"
// @Param service_catalog_id query string false "Filter by service catalog ID"
// @Param search query string false "Search by business name or description"
// @Success 200 {object} model.APIResponse
// @Router /services/providers [get]
func (h *ServiceCatalogHandler) GetAllProviders(c *gin.Context) {
	filters := &repository.ServiceProviderFilters{}

	if city := c.Query("city"); city != "" {
		filters.City = &city
	}

	if province := c.Query("province"); province != "" {
		filters.Province = &province
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filters.IsActive = &isActive
	}

	if isVerifiedStr := c.Query("is_verified"); isVerifiedStr != "" {
		isVerified := isVerifiedStr == "true"
		filters.IsVerified = &isVerified
	}

	if serviceCatalogIDStr := c.Query("service_catalog_id"); serviceCatalogIDStr != "" {
		if serviceCatalogID, err := uuid.Parse(serviceCatalogIDStr); err == nil {
			filters.ServiceCatalogID = &serviceCatalogID
		}
	}

	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	providers, err := h.providerService.GetAllProviders(c.Request.Context(), filters)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(providers, "Providers retrieved successfully"))
}

// UpdateProvider updates a provider
// @Summary Update service provider
// @Description Update a service provider
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Param request body model.ServiceProviderRequest true "Provider request"
// @Success 200 {object} model.APIResponse
// @Router /services/providers/{id} [put]
func (h *ServiceCatalogHandler) UpdateProvider(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	var req model.ServiceProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	provider, err := h.providerService.UpdateProvider(c.Request.Context(), id, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(provider, "Provider updated successfully"))
}

// DeleteProvider deletes a provider
// @Summary Delete service provider
// @Description Delete a service provider (Admin only)
// @Tags Services
// @Produce json
// @Param id path string true "Provider ID"
// @Success 200 {object} model.APIResponse
// @Router /services/providers/{id} [delete]
func (h *ServiceCatalogHandler) DeleteProvider(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	if err := h.providerService.DeleteProvider(c.Request.Context(), id); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil, "Provider deleted successfully"))
}

// AddServiceToProvider adds a service to a provider
// @Summary Add service to provider
// @Description Add a service to a provider
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Param serviceId path string true "Service Catalog ID"
// @Param request body object true "Price" SchemaExample({"price": 100000})
// @Success 200 {object} model.APIResponse
// @Router /services/providers/{id}/services/{serviceId} [post]
func (h *ServiceCatalogHandler) AddServiceToProvider(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	serviceID, err := uuid.Parse(c.Param("serviceId"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	var req struct {
		Price float64 `json:"price" validate:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.providerService.AddServiceToProvider(c.Request.Context(), providerID, serviceID, req.Price); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil, "Service added to provider successfully"))
}

// RemoveServiceFromProvider removes a service from a provider
// @Summary Remove service from provider
// @Description Remove a service from a provider
// @Tags Services
// @Produce json
// @Param id path string true "Provider ID"
// @Param serviceId path string true "Service Catalog ID"
// @Success 200 {object} model.APIResponse
// @Router /services/providers/{id}/services/{serviceId} [delete]
func (h *ServiceCatalogHandler) RemoveServiceFromProvider(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	serviceID, err := uuid.Parse(c.Param("serviceId"))
	if err != nil {
		utils.HandleError(c, model.ErrInvalidInput)
		return
	}

	if err := h.providerService.RemoveServiceFromProvider(c.Request.Context(), providerID, serviceID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil, "Service removed from provider successfully"))
}

