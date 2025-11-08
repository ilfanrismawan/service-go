package repository

import (
	"context"
	"service/internal/shared/database"
	"service/internal/shared/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServiceCategoryRepository handles service category data operations
type ServiceCategoryRepository struct {
	db *gorm.DB
}

// NewServiceCategoryRepository creates a new service category repository
func NewServiceCategoryRepository() *ServiceCategoryRepository {
	return &ServiceCategoryRepository{
		db: database.DB,
	}
}

// Create creates a new service category
func (r *ServiceCategoryRepository) Create(ctx context.Context, category *model.ServiceCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetByID retrieves a category by ID
func (r *ServiceCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ServiceCategory, error) {
	var category model.ServiceCategory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll retrieves all categories
func (r *ServiceCategoryRepository) GetAll(ctx context.Context, includeInactive bool) ([]model.ServiceCategory, error) {
	var categories []model.ServiceCategory
	query := r.db.WithContext(ctx).Order("sort_order ASC, name ASC")
	
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	
	err := query.Find(&categories).Error
	return categories, err
}

// Update updates a category
func (r *ServiceCategoryRepository) Update(ctx context.Context, category *model.ServiceCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete soft deletes a category
func (r *ServiceCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ServiceCategory{}, id).Error
}

// ServiceCatalogRepository handles service catalog data operations
type ServiceCatalogRepository struct {
	db *gorm.DB
}

// NewServiceCatalogRepository creates a new service catalog repository
func NewServiceCatalogRepository() *ServiceCatalogRepository {
	return &ServiceCatalogRepository{
		db: database.DB,
	}
}

// Create creates a new service catalog
func (r *ServiceCatalogRepository) Create(ctx context.Context, catalog *model.ServiceCatalog) error {
	return r.db.WithContext(ctx).Create(catalog).Error
}

// GetByID retrieves a catalog by ID with category
func (r *ServiceCatalogRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ServiceCatalog, error) {
	var catalog model.ServiceCatalog
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("id = ?", id).
		First(&catalog).Error
	if err != nil {
		return nil, err
	}
	return &catalog, nil
}

// GetAll retrieves all catalogs with optional filters
func (r *ServiceCatalogRepository) GetAll(ctx context.Context, filters *ServiceCatalogFilters) ([]model.ServiceCatalog, error) {
	var catalogs []model.ServiceCatalog
	query := r.db.WithContext(ctx).Preload("Category")
	
	if filters != nil {
		if filters.CategoryID != nil {
			query = query.Where("category_id = ?", *filters.CategoryID)
		}
		if filters.IsActive != nil {
			query = query.Where("is_active = ?", *filters.IsActive)
		}
		if filters.RequiresAppointment != nil {
			query = query.Where("requires_appointment = ?", *filters.RequiresAppointment)
		}
		if filters.Search != nil && *filters.Search != "" {
			search := "%" + *filters.Search + "%"
			query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
		}
	}
	
	err := query.Order("name ASC").Find(&catalogs).Error
	return catalogs, err
}

// GetByCategoryID retrieves all catalogs in a category
func (r *ServiceCatalogRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID, includeInactive bool) ([]model.ServiceCatalog, error) {
	var catalogs []model.ServiceCatalog
	query := r.db.WithContext(ctx).
		Preload("Category").
		Where("category_id = ?", categoryID)
	
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	
	err := query.Order("name ASC").Find(&catalogs).Error
	return catalogs, err
}

// Update updates a catalog
func (r *ServiceCatalogRepository) Update(ctx context.Context, catalog *model.ServiceCatalog) error {
	return r.db.WithContext(ctx).Save(catalog).Error
}

// Delete soft deletes a catalog
func (r *ServiceCatalogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ServiceCatalog{}, id).Error
}

// ServiceCatalogFilters represents filters for querying service catalogs
type ServiceCatalogFilters struct {
	CategoryID          *uuid.UUID
	IsActive            *bool
	RequiresAppointment *bool
	Search              *string
}

// ServiceProviderRepository handles service provider data operations
type ServiceProviderRepository struct {
	db *gorm.DB
}

// NewServiceProviderRepository creates a new service provider repository
func NewServiceProviderRepository() *ServiceProviderRepository {
	return &ServiceProviderRepository{
		db: database.DB,
	}
}

// Create creates a new service provider
func (r *ServiceProviderRepository) Create(ctx context.Context, provider *model.ServiceProvider) error {
	return r.db.WithContext(ctx).Create(provider).Error
}

// GetByID retrieves a provider by ID with relations
func (r *ServiceProviderRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ServiceProvider, error) {
	var provider model.ServiceProvider
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Services").
		Preload("Services.Category").
		Where("id = ?", id).
		First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// GetByUserID retrieves a provider by user ID
func (r *ServiceProviderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.ServiceProvider, error) {
	var provider model.ServiceProvider
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Services").
		Preload("Services.Category").
		Where("user_id = ?", userID).
		First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// GetAll retrieves all providers with optional filters
func (r *ServiceProviderRepository) GetAll(ctx context.Context, filters *ServiceProviderFilters) ([]model.ServiceProvider, error) {
	var providers []model.ServiceProvider
	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("Services").
		Preload("Services.Category")
	
	if filters != nil {
		if filters.City != nil && *filters.City != "" {
			query = query.Where("city = ?", *filters.City)
		}
		if filters.Province != nil && *filters.Province != "" {
			query = query.Where("province = ?", *filters.Province)
		}
		if filters.IsActive != nil {
			query = query.Where("is_active = ?", *filters.IsActive)
		}
		if filters.IsVerified != nil {
			query = query.Where("is_verified = ?", *filters.IsVerified)
		}
		if filters.ServiceCatalogID != nil {
			query = query.Joins("JOIN provider_services ON provider_services.provider_id = service_providers.id").
				Where("provider_services.service_catalog_id = ?", *filters.ServiceCatalogID)
		}
		if filters.Search != nil && *filters.Search != "" {
			search := "%" + *filters.Search + "%"
			query = query.Where("business_name ILIKE ? OR description ILIKE ?", search, search)
		}
	}
	
	err := query.Order("rating DESC, business_name ASC").Find(&providers).Error
	return providers, err
}

// Update updates a provider
func (r *ServiceProviderRepository) Update(ctx context.Context, provider *model.ServiceProvider) error {
	return r.db.WithContext(ctx).Save(provider).Error
}

// Delete soft deletes a provider
func (r *ServiceProviderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ServiceProvider{}, id).Error
}

// AddService adds a service to a provider
func (r *ServiceProviderRepository) AddService(ctx context.Context, providerID uuid.UUID, serviceCatalogID uuid.UUID, price float64) error {
	providerService := model.ProviderService{
		ProviderID:      providerID,
		ServiceCatalogID: serviceCatalogID,
		Price:           price,
		IsActive:        true,
	}
	return r.db.WithContext(ctx).Create(&providerService).Error
}

// RemoveService removes a service from a provider
func (r *ServiceProviderRepository) RemoveService(ctx context.Context, providerID uuid.UUID, serviceCatalogID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("provider_id = ? AND service_catalog_id = ?", providerID, serviceCatalogID).
		Delete(&model.ProviderService{}).Error
}

// GetProviderServices retrieves all services for a provider
func (r *ServiceProviderRepository) GetProviderServices(ctx context.Context, providerID uuid.UUID) ([]model.ServiceCatalog, error) {
	var services []model.ServiceCatalog
	err := r.db.WithContext(ctx).
		Joins("JOIN provider_services ON provider_services.service_catalog_id = service_catalogs.id").
		Where("provider_services.provider_id = ? AND provider_services.is_active = ?", providerID, true).
		Preload("Category").
		Find(&services).Error
	return services, err
}

// ServiceProviderFilters represents filters for querying service providers
type ServiceProviderFilters struct {
	City            *string
	Province        *string
	IsActive        *bool
	IsVerified      *bool
	ServiceCatalogID *uuid.UUID
	Search          *string
}

