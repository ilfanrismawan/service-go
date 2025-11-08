package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServiceCategory represents a category of services (e.g., Beauty, Electronics, Health)
type ServiceCategory struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"` // URL or icon identifier
	ImageURL    string         `json:"image_url"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ServiceCategory) TableName() string {
	return "service_categories"
}

// ServiceCatalog represents a service that can be offered (e.g., Nail Art, iPhone Repair, Botox)
type ServiceCatalog struct {
	ID                  uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CategoryID          uuid.UUID      `json:"category_id" gorm:"type:uuid;not null"`
	Category            ServiceCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Name                string         `json:"name" gorm:"not null"`
	Description         string         `json:"description"`
	ImageURL            string         `json:"image_url"`
	BasePrice           float64        `json:"base_price" gorm:"default:0"`
	EstimatedDuration   int            `json:"estimated_duration" gorm:"default:0"` // in minutes
	RequiresPickup      bool           `json:"requires_pickup" gorm:"default:false"`
	RequiresDelivery    bool           `json:"requires_delivery" gorm:"default:false"`
	RequiresAppointment bool           `json:"requires_appointment" gorm:"default:false"`
	RequiresItem        bool           `json:"requires_item" gorm:"default:false"` // requires customer to bring item
	RequiresLocation    bool           `json:"requires_location" gorm:"default:true"` // service performed at location
	Metadata            JSONB          `json:"metadata" gorm:"type:jsonb"` // for service-specific fields
	IsActive            bool           `json:"is_active" gorm:"default:true"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ServiceCatalog) TableName() string {
	return "service_catalogs"
}

// ServiceProvider represents a provider/merchant that offers services
type ServiceProvider struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User            User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	BusinessName    string         `json:"business_name" gorm:"not null"`
	BusinessType    string         `json:"business_type"` // e.g., "Salon", "Clinic", "Workshop"
	Description     string         `json:"description"`
	Address         string         `json:"address" gorm:"not null"`
	City            string         `json:"city" gorm:"not null"`
	Province        string         `json:"province" gorm:"not null"`
	Phone           string         `json:"phone" gorm:"not null"`
	Email           string         `json:"email"`
	Latitude        float64        `json:"latitude" gorm:"type:decimal(10,6);not null"`
	Longitude       float64        `json:"longitude" gorm:"type:decimal(10,6);not null"`
	Rating          float64        `json:"rating" gorm:"default:0"`
	TotalReviews    int            `json:"total_reviews" gorm:"default:0"`
	ImageURL        string         `json:"image_url"`
	IsVerified      bool           `json:"is_verified" gorm:"default:false"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Many-to-many relationship with ServiceCatalog
	Services        []ServiceCatalog `json:"services,omitempty" gorm:"many2many:provider_services;"`
}

func (ServiceProvider) TableName() string {
	return "service_providers"
}

// ProviderService represents the many-to-many relationship between providers and services
type ProviderService struct {
	ProviderID      uuid.UUID      `gorm:"type:uuid;primary_key"`
	ServiceCatalogID uuid.UUID     `gorm:"type:uuid;primary_key"`
	Price           float64        `gorm:"not null"` // provider-specific price
	IsActive        bool           `gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (ProviderService) TableName() string {
	return "provider_services"
}

// JSONB is a custom type for handling JSONB in PostgreSQL
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

// ServiceCategoryRequest represents request payload for creating/updating category
type ServiceCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ImageURL    string `json:"image_url"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// ServiceCategoryResponse represents response payload for category
type ServiceCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	ImageURL    string    `json:"image_url"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (sc *ServiceCategory) ToResponse() ServiceCategoryResponse {
	return ServiceCategoryResponse{
		ID:          sc.ID,
		Name:        sc.Name,
		Description: sc.Description,
		Icon:        sc.Icon,
		ImageURL:    sc.ImageURL,
		SortOrder:   sc.SortOrder,
		IsActive:    sc.IsActive,
		CreatedAt:   sc.CreatedAt,
		UpdatedAt:   sc.UpdatedAt,
	}
}

// ServiceCatalogRequest represents request payload for creating/updating catalog
type ServiceCatalogRequest struct {
	CategoryID          string  `json:"category_id" validate:"required"`
	Name                string  `json:"name" validate:"required"`
	Description         string  `json:"description"`
	ImageURL            string  `json:"image_url"`
	BasePrice           float64 `json:"base_price"`
	EstimatedDuration   int     `json:"estimated_duration"`
	RequiresPickup      bool    `json:"requires_pickup"`
	RequiresDelivery    bool    `json:"requires_delivery"`
	RequiresAppointment bool    `json:"requires_appointment"`
	RequiresItem        bool    `json:"requires_item"`
	RequiresLocation    bool    `json:"requires_location"`
	Metadata            JSONB   `json:"metadata"`
	IsActive            bool    `json:"is_active"`
}

// ServiceCatalogResponse represents response payload for catalog
type ServiceCatalogResponse struct {
	ID                  uuid.UUID              `json:"id"`
	CategoryID          uuid.UUID              `json:"category_id"`
	Category            *ServiceCategoryResponse `json:"category,omitempty"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	ImageURL            string                 `json:"image_url"`
	BasePrice           float64                `json:"base_price"`
	EstimatedDuration   int                    `json:"estimated_duration"`
	RequiresPickup      bool                   `json:"requires_pickup"`
	RequiresDelivery    bool                   `json:"requires_delivery"`
	RequiresAppointment bool                   `json:"requires_appointment"`
	RequiresItem        bool                   `json:"requires_item"`
	RequiresLocation    bool                   `json:"requires_location"`
	Metadata            JSONB                  `json:"metadata"`
	IsActive            bool                   `json:"is_active"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

func (sc *ServiceCatalog) ToResponse() ServiceCatalogResponse {
	resp := ServiceCatalogResponse{
		ID:                  sc.ID,
		CategoryID:          sc.CategoryID,
		Name:                sc.Name,
		Description:         sc.Description,
		ImageURL:            sc.ImageURL,
		BasePrice:           sc.BasePrice,
		EstimatedDuration:   sc.EstimatedDuration,
		RequiresPickup:      sc.RequiresPickup,
		RequiresDelivery:    sc.RequiresDelivery,
		RequiresAppointment: sc.RequiresAppointment,
		RequiresItem:        sc.RequiresItem,
		RequiresLocation:    sc.RequiresLocation,
		Metadata:            sc.Metadata,
		IsActive:            sc.IsActive,
		CreatedAt:           sc.CreatedAt,
		UpdatedAt:           sc.UpdatedAt,
	}
	
	if sc.Category.ID != uuid.Nil {
		categoryResp := sc.Category.ToResponse()
		resp.Category = &categoryResp
	}
	
	return resp
}

// ServiceProviderRequest represents request payload for creating/updating provider
type ServiceProviderRequest struct {
	UserID       string   `json:"user_id" validate:"required"`
	BusinessName string   `json:"business_name" validate:"required"`
	BusinessType string   `json:"business_type"`
	Description  string   `json:"description"`
	Address      string   `json:"address" validate:"required"`
	City         string   `json:"city" validate:"required"`
	Province     string   `json:"province" validate:"required"`
	Phone        string   `json:"phone" validate:"required"`
	Email        string   `json:"email"`
	Latitude     float64  `json:"latitude" validate:"required"`
	Longitude    float64  `json:"longitude" validate:"required"`
	ImageURL     string   `json:"image_url"`
	ServiceIDs   []string `json:"service_ids"` // IDs of services this provider offers
}

// ServiceProviderResponse represents response payload for provider
type ServiceProviderResponse struct {
	ID           uuid.UUID                `json:"id"`
	UserID       uuid.UUID                `json:"user_id"`
	User         *UserResponse            `json:"user,omitempty"`
	BusinessName string                  `json:"business_name"`
	BusinessType string                  `json:"business_type"`
	Description  string                  `json:"description"`
	Address      string                  `json:"address"`
	City         string                  `json:"city"`
	Province     string                  `json:"province"`
	Phone        string                  `json:"phone"`
	Email        string                  `json:"email"`
	Latitude     float64                 `json:"latitude"`
	Longitude    float64                 `json:"longitude"`
	Rating       float64                 `json:"rating"`
	TotalReviews int                     `json:"total_reviews"`
	ImageURL     string                  `json:"image_url"`
	IsVerified   bool                    `json:"is_verified"`
	IsActive     bool                    `json:"is_active"`
	Services     []ServiceCatalogResponse `json:"services,omitempty"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

func (sp *ServiceProvider) ToResponse() ServiceProviderResponse {
	resp := ServiceProviderResponse{
		ID:           sp.ID,
		UserID:       sp.UserID,
		BusinessName: sp.BusinessName,
		BusinessType: sp.BusinessType,
		Description:  sp.Description,
		Address:      sp.Address,
		City:         sp.City,
		Province:     sp.Province,
		Phone:        sp.Phone,
		Email:        sp.Email,
		Latitude:     sp.Latitude,
		Longitude:    sp.Longitude,
		Rating:       sp.Rating,
		TotalReviews: sp.TotalReviews,
		ImageURL:     sp.ImageURL,
		IsVerified:   sp.IsVerified,
		IsActive:     sp.IsActive,
		CreatedAt:    sp.CreatedAt,
		UpdatedAt:    sp.UpdatedAt,
	}
	
	if sp.User.ID != uuid.Nil {
		userResp := sp.User.ToResponse()
		resp.User = &userResp
	}
	
	if len(sp.Services) > 0 {
		resp.Services = make([]ServiceCatalogResponse, len(sp.Services))
		for i, service := range sp.Services {
			resp.Services[i] = service.ToResponse()
		}
	}
	
	return resp
}

