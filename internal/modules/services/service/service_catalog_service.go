package service

import (
	"context"
	"errors"
	"service/internal/modules/services/repository"
	userRepo "service/internal/modules/users/repository"
	"service/internal/shared/model"

	"github.com/google/uuid"
)

// ServiceCategoryService handles service category business logic
type ServiceCategoryService struct {
	categoryRepo *repository.ServiceCategoryRepository
}

// NewServiceCategoryService creates a new service category service
func NewServiceCategoryService() *ServiceCategoryService {
	return &ServiceCategoryService{
		categoryRepo: repository.NewServiceCategoryRepository(),
	}
}

// CreateCategory creates a new service category
func (s *ServiceCategoryService) CreateCategory(ctx context.Context, req *model.ServiceCategoryRequest) (*model.ServiceCategoryResponse, error) {
	category := &model.ServiceCategory{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	response := category.ToResponse()
	return &response, nil
}

// GetCategory retrieves a category by ID
func (s *ServiceCategoryService) GetCategory(ctx context.Context, id uuid.UUID) (*model.ServiceCategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrCategoryNotFound
	}

	response := category.ToResponse()
	return &response, nil
}

// GetAllCategories retrieves all categories
func (s *ServiceCategoryService) GetAllCategories(ctx context.Context, includeInactive bool) ([]model.ServiceCategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll(ctx, includeInactive)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ServiceCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = category.ToResponse()
	}

	return responses, nil
}

// UpdateCategory updates a category
func (s *ServiceCategoryService) UpdateCategory(ctx context.Context, id uuid.UUID, req *model.ServiceCategoryRequest) (*model.ServiceCategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrCategoryNotFound
	}

	category.Name = req.Name
	category.Description = req.Description
	category.Icon = req.Icon
	category.ImageURL = req.ImageURL
	category.SortOrder = req.SortOrder
	category.IsActive = req.IsActive

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	response := category.ToResponse()
	return &response, nil
}

// DeleteCategory deletes a category
func (s *ServiceCategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return model.ErrCategoryNotFound
	}

	return s.categoryRepo.Delete(ctx, id)
}

// ServiceCatalogService handles service catalog business logic
type ServiceCatalogService struct {
	catalogRepo *repository.ServiceCatalogRepository
	categoryRepo *repository.ServiceCategoryRepository
}

// NewServiceCatalogService creates a new service catalog service
func NewServiceCatalogService() *ServiceCatalogService {
	return &ServiceCatalogService{
		catalogRepo: repository.NewServiceCatalogRepository(),
		categoryRepo: repository.NewServiceCategoryRepository(),
	}
}

// CreateCatalog creates a new service catalog
func (s *ServiceCatalogService) CreateCatalog(ctx context.Context, req *model.ServiceCatalogRequest) (*model.ServiceCatalogResponse, error) {
	// Validate category exists
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	_, err = s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, model.ErrCategoryNotFound
	}

	catalog := &model.ServiceCatalog{
		CategoryID:          categoryID,
		Name:                req.Name,
		Description:         req.Description,
		ImageURL:            req.ImageURL,
		BasePrice:           req.BasePrice,
		EstimatedDuration:   req.EstimatedDuration,
		RequiresPickup:      req.RequiresPickup,
		RequiresDelivery:    req.RequiresDelivery,
		RequiresAppointment: req.RequiresAppointment,
		RequiresItem:        req.RequiresItem,
		RequiresLocation:    req.RequiresLocation,
		Metadata:            req.Metadata,
		IsActive:            req.IsActive,
	}

	if err := s.catalogRepo.Create(ctx, catalog); err != nil {
		return nil, err
	}

	// Reload with category
	catalog, err = s.catalogRepo.GetByID(ctx, catalog.ID)
	if err != nil {
		return nil, err
	}

	response := catalog.ToResponse()
	return &response, nil
}

// GetCatalog retrieves a catalog by ID
func (s *ServiceCatalogService) GetCatalog(ctx context.Context, id uuid.UUID) (*model.ServiceCatalogResponse, error) {
	catalog, err := s.catalogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrCatalogNotFound
	}

	response := catalog.ToResponse()
	return &response, nil
}

// GetAllCatalogs retrieves all catalogs with optional filters
func (s *ServiceCatalogService) GetAllCatalogs(ctx context.Context, filters *repository.ServiceCatalogFilters) ([]model.ServiceCatalogResponse, error) {
	catalogs, err := s.catalogRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ServiceCatalogResponse, len(catalogs))
	for i, catalog := range catalogs {
		responses[i] = catalog.ToResponse()
	}

	return responses, nil
}

// GetCatalogsByCategory retrieves all catalogs in a category
func (s *ServiceCatalogService) GetCatalogsByCategory(ctx context.Context, categoryID uuid.UUID, includeInactive bool) ([]model.ServiceCatalogResponse, error) {
	catalogs, err := s.catalogRepo.GetByCategoryID(ctx, categoryID, includeInactive)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ServiceCatalogResponse, len(catalogs))
	for i, catalog := range catalogs {
		responses[i] = catalog.ToResponse()
	}

	return responses, nil
}

// UpdateCatalog updates a catalog
func (s *ServiceCatalogService) UpdateCatalog(ctx context.Context, id uuid.UUID, req *model.ServiceCatalogRequest) (*model.ServiceCatalogResponse, error) {
	catalog, err := s.catalogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrCatalogNotFound
	}

	// Validate category if changed
	if req.CategoryID != "" {
		categoryID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category ID")
		}

		_, err = s.categoryRepo.GetByID(ctx, categoryID)
		if err != nil {
			return nil, model.ErrCategoryNotFound
		}

		catalog.CategoryID = categoryID
	}

	catalog.Name = req.Name
	catalog.Description = req.Description
	catalog.ImageURL = req.ImageURL
	catalog.BasePrice = req.BasePrice
	catalog.EstimatedDuration = req.EstimatedDuration
	catalog.RequiresPickup = req.RequiresPickup
	catalog.RequiresDelivery = req.RequiresDelivery
	catalog.RequiresAppointment = req.RequiresAppointment
	catalog.RequiresItem = req.RequiresItem
	catalog.RequiresLocation = req.RequiresLocation
	catalog.Metadata = req.Metadata
	catalog.IsActive = req.IsActive

	if err := s.catalogRepo.Update(ctx, catalog); err != nil {
		return nil, err
	}

	// Reload with category
	catalog, err = s.catalogRepo.GetByID(ctx, catalog.ID)
	if err != nil {
		return nil, err
	}

	response := catalog.ToResponse()
	return &response, nil
}

// DeleteCatalog deletes a catalog
func (s *ServiceCatalogService) DeleteCatalog(ctx context.Context, id uuid.UUID) error {
	_, err := s.catalogRepo.GetByID(ctx, id)
	if err != nil {
		return model.ErrCatalogNotFound
	}

	return s.catalogRepo.Delete(ctx, id)
}

// ServiceProviderService handles service provider business logic
type ServiceProviderService struct {
	providerRepo *repository.ServiceProviderRepository
	catalogRepo  *repository.ServiceCatalogRepository
	userRepo     *userRepo.UserRepository
}

// NewServiceProviderService creates a new service provider service
func NewServiceProviderService() *ServiceProviderService {
	return &ServiceProviderService{
		providerRepo: repository.NewServiceProviderRepository(),
		catalogRepo:  repository.NewServiceCatalogRepository(),
		userRepo:     userRepo.NewUserRepository(),
	}
}

// CreateProvider creates a new service provider
func (s *ServiceProviderService) CreateProvider(ctx context.Context, req *model.ServiceProviderRequest) (*model.ServiceProviderResponse, error) {
	// Validate user exists and has provider role
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Validate user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	// Validate user has provider role
	if user.Role != model.UserRoleProvider {
		return nil, errors.New("user does not have provider role")
	}

	provider := &model.ServiceProvider{
		UserID:       userID,
		BusinessName: req.BusinessName,
		BusinessType: req.BusinessType,
		Description:  req.Description,
		Address:      req.Address,
		City:         req.City,
		Province:     req.Province,
		Phone:        req.Phone,
		Email:        req.Email,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		ImageURL:     req.ImageURL,
		IsActive:     true,
	}

	if err := s.providerRepo.Create(ctx, provider); err != nil {
		return nil, err
	}

	// Add services if provided
	if len(req.ServiceIDs) > 0 {
		for _, serviceIDStr := range req.ServiceIDs {
			serviceID, err := uuid.Parse(serviceIDStr)
			if err != nil {
				continue
			}

			// Get catalog to get base price
			catalog, err := s.catalogRepo.GetByID(ctx, serviceID)
			if err != nil {
				continue
			}

			// Add service with base price (can be updated later)
			_ = s.providerRepo.AddService(ctx, provider.ID, serviceID, catalog.BasePrice)
		}
	}

	// Reload with relations
	provider, err = s.providerRepo.GetByID(ctx, provider.ID)
	if err != nil {
		return nil, err
	}

	response := provider.ToResponse()
	return &response, nil
}

// GetProvider retrieves a provider by ID
func (s *ServiceProviderService) GetProvider(ctx context.Context, id uuid.UUID) (*model.ServiceProviderResponse, error) {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrProviderNotFound
	}

	response := provider.ToResponse()
	return &response, nil
}

// GetProviderByUserID retrieves a provider by user ID
func (s *ServiceProviderService) GetProviderByUserID(ctx context.Context, userID uuid.UUID) (*model.ServiceProviderResponse, error) {
	provider, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, model.ErrProviderNotFound
	}

	response := provider.ToResponse()
	return &response, nil
}

// GetAllProviders retrieves all providers with optional filters
func (s *ServiceProviderService) GetAllProviders(ctx context.Context, filters *repository.ServiceProviderFilters) ([]model.ServiceProviderResponse, error) {
	providers, err := s.providerRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ServiceProviderResponse, len(providers))
	for i, provider := range providers {
		responses[i] = provider.ToResponse()
	}

	return responses, nil
}

// UpdateProvider updates a provider
func (s *ServiceProviderService) UpdateProvider(ctx context.Context, id uuid.UUID, req *model.ServiceProviderRequest) (*model.ServiceProviderResponse, error) {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrProviderNotFound
	}

	provider.BusinessName = req.BusinessName
	provider.BusinessType = req.BusinessType
	provider.Description = req.Description
	provider.Address = req.Address
	provider.City = req.City
	provider.Province = req.Province
	provider.Phone = req.Phone
	provider.Email = req.Email
	provider.Latitude = req.Latitude
	provider.Longitude = req.Longitude
	provider.ImageURL = req.ImageURL

	if err := s.providerRepo.Update(ctx, provider); err != nil {
		return nil, err
	}

	// Update services if provided
	if len(req.ServiceIDs) > 0 {
		// Get current services
		currentServices, _ := s.providerRepo.GetProviderServices(ctx, id)
		currentServiceMap := make(map[uuid.UUID]bool)
		for _, svc := range currentServices {
			currentServiceMap[svc.ID] = true
		}

		// Add new services
		for _, serviceIDStr := range req.ServiceIDs {
			serviceID, err := uuid.Parse(serviceIDStr)
			if err != nil {
				continue
			}

			if !currentServiceMap[serviceID] {
				catalog, err := s.catalogRepo.GetByID(ctx, serviceID)
				if err == nil {
					_ = s.providerRepo.AddService(ctx, id, serviceID, catalog.BasePrice)
				}
			}
		}
	}

	// Reload with relations
	provider, err = s.providerRepo.GetByID(ctx, provider.ID)
	if err != nil {
		return nil, err
	}

	response := provider.ToResponse()
	return &response, nil
}

// DeleteProvider deletes a provider
func (s *ServiceProviderService) DeleteProvider(ctx context.Context, id uuid.UUID) error {
	_, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return model.ErrProviderNotFound
	}

	return s.providerRepo.Delete(ctx, id)
}

// AddServiceToProvider adds a service to a provider
func (s *ServiceProviderService) AddServiceToProvider(ctx context.Context, providerID uuid.UUID, serviceCatalogID uuid.UUID, price float64) error {
	_, err := s.providerRepo.GetByID(ctx, providerID)
	if err != nil {
		return model.ErrProviderNotFound
	}

	_, err = s.catalogRepo.GetByID(ctx, serviceCatalogID)
	if err != nil {
		return model.ErrCatalogNotFound
	}

	return s.providerRepo.AddService(ctx, providerID, serviceCatalogID, price)
}

// RemoveServiceFromProvider removes a service from a provider
func (s *ServiceProviderService) RemoveServiceFromProvider(ctx context.Context, providerID uuid.UUID, serviceCatalogID uuid.UUID) error {
	_, err := s.providerRepo.GetByID(ctx, providerID)
	if err != nil {
		return model.ErrProviderNotFound
	}

	return s.providerRepo.RemoveService(ctx, providerID, serviceCatalogID)
}

