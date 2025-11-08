package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPendingPickup OrderStatus = "pending_pickup"
	StatusOnPickup      OrderStatus = "on_pickup"
	StatusInService     OrderStatus = "in_service"
	StatusReady         OrderStatus = "ready"
	StatusDelivered     OrderStatus = "delivered"
	StatusCompleted     OrderStatus = "completed"
	StatusCancelled     OrderStatus = "cancelled"
)

type ServiceType string

const (
	ServiceTypeScreenRepair       ServiceType = "screen_repair"
	ServiceTypeBatteryReplacement ServiceType = "battery_replacement"
	ServiceTypeWaterDamage        ServiceType = "water_damage"
	ServiceTypeSoftwareIssue      ServiceType = "software_issue"
	ServiceTypeHardwareRepair     ServiceType = "hardware_repair"
	ServiceTypeOther              ServiceType = "other"
)

type ServiceOrder struct {
	ID                uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderNumber       string      `gorm:"uniqueIndex;not null"`
	CustomerID        uuid.UUID   `gorm:"type:uuid;not null"`
	
	// New multi-service fields
	ServiceCatalogID  *uuid.UUID  `gorm:"type:uuid"` // Reference to ServiceCatalog
	ServiceProviderID *uuid.UUID  `gorm:"type:uuid"` // Reference to ServiceProvider
	
	// Legacy fields (kept for backward compatibility)
	BranchID          *uuid.UUID  `gorm:"type:uuid"` // Made optional, can come from provider
	TechnicianID      *uuid.UUID  `gorm:"type:uuid"`
	CourierID         *uuid.UUID  `gorm:"type:uuid"`
	
	// Legacy iPhone-specific fields (kept for backward compatibility, now optional)
	IPhoneModel       *string     `gorm:"type:varchar(100)"` // Made optional
	IPhoneColor       *string     `gorm:"type:varchar(50)"`   // Made optional
	IPhoneIMEI        *string     `gorm:"type:varchar(50)"`  // Made optional
	IPhoneType        *string     `gorm:"type:varchar(100)"` // Made optional
	
	// Generic service fields
	ServiceType       ServiceType `gorm:"type:varchar(50)"` // Kept for backward compatibility, can be derived from ServiceCatalog
	ServiceName       string      `gorm:"type:varchar(255)"` // Name from ServiceCatalog
	Description       string      `gorm:"not null"`
	Complaint         string      `gorm:"type:text"` // Made optional for non-repair services
	
	// Generic item/device info (for services that require items)
	ItemModel         *string     `gorm:"type:varchar(100)"` // Generic item model
	ItemColor         *string     `gorm:"type:varchar(50)"`  // Generic item color
	ItemSerial        *string     `gorm:"type:varchar(100)"` // Generic item serial/IMEI
	ItemType          *string     `gorm:"type:varchar(100)"` // Generic item type
	
	// Appointment fields (for services that require appointment)
	AppointmentDate   *time.Time  `gorm:"type:timestamp"`
	AppointmentTime   *time.Time  `gorm:"type:time"`
	
	// Location fields (made optional, not all services need pickup)
	PickupAddress     *string     `gorm:"type:text"` // Made optional
	PickupLocation    *string     `gorm:"type:text"` // Made optional
	PickupLatitude    *float64    `gorm:"type:decimal(10,6)"` // Made optional
	PickupLongitude   *float64    `gorm:"type:decimal(10,6)"` // Made optional
	ServiceLocation   string      `gorm:"type:text"` // Location where service is performed (provider location or branch)
	
	// On-demand service fields
	IsOnDemand        bool        `gorm:"default:false"` // Service datang ke customer
	CurrentLatitude   *float64    `gorm:"type:decimal(10,6)"` // Current location courier/provider
	CurrentLongitude  *float64    `gorm:"type:decimal(10,6)"` // Current location courier/provider
	ETA               *int        `gorm:"default:0"` // Estimated time of arrival in minutes
	LastLocationUpdate *time.Time `gorm:"type:timestamp"` // Last location update time
	
	// Status and pricing
	Status            OrderStatus `gorm:"type:varchar(50);not null;default:'pending_pickup'"`
	EstimatedCost     float64     `gorm:"type:decimal(15,2);default:0"`
	ActualCost        float64     `gorm:"type:decimal(15,2);default:0"`
	ServiceCost       float64     `gorm:"type:decimal(15,2);default:0"`
	EstimatedDuration int         `gorm:"default:0"` // in minutes
	ActualDuration    int         `gorm:"default:0"` // in minutes
	
	// Photos
	PickupPhoto       string      `gorm:"type:text"`
	ServicePhoto      string      `gorm:"type:text"`
	DeliveryPhoto     string      `gorm:"type:text"`
	
	// Additional fields
	Notes             string      `gorm:"type:text"`
	InvoiceNumber     string      `gorm:"type:varchar(100);uniqueIndex"`
	TaxAmount         float64     `gorm:"type:decimal(15,2);default:0"`
	CustomerIDCard    string      `gorm:"type:varchar(50)"`
	CustomerNPWP      string      `gorm:"type:varchar(50)"`
	TermsAccepted     bool        `gorm:"default:false"`
	PrivacyAccepted   bool        `gorm:"default:false"`
	
	// Metadata for service-specific fields (JSONB)
	Metadata          JSONB       `gorm:"type:jsonb"` // For service-specific metadata
	
	// Timestamps
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	
	// Relations (will be preloaded)
	ServiceCatalog    *ServiceCatalog `gorm:"foreignKey:ServiceCatalogID"`
	ServiceProvider   *ServiceProvider `gorm:"foreignKey:ServiceProviderID"`
	Customer          *User           `gorm:"foreignKey:CustomerID"`
	Branch            *Branch         `gorm:"foreignKey:BranchID"`
	Technician        *User           `gorm:"foreignKey:TechnicianID"`
	Courier           *User           `gorm:"foreignKey:CourierID"`
}

func (ServiceOrder) TableName() string {
	return "service_orders"
}

func (so *ServiceOrder) SetAliasFields() {
	// Legacy compatibility: map iPhone fields to generic item fields if not set
	if so.ItemModel == nil && so.IPhoneModel != nil {
		so.ItemModel = so.IPhoneModel
	}
	if so.ItemColor == nil && so.IPhoneColor != nil {
		so.ItemColor = so.IPhoneColor
	}
	if so.ItemSerial == nil && so.IPhoneIMEI != nil {
		so.ItemSerial = so.IPhoneIMEI
	}
	if so.ItemType == nil && so.IPhoneType != nil {
		so.ItemType = so.IPhoneType
	}
	if so.ItemType == nil && so.IPhoneModel != nil {
		so.ItemType = so.IPhoneModel
	}
	
	// Legacy compatibility
	if so.Complaint == "" {
		so.Complaint = so.Description
	}
	if so.PickupLocation == nil && so.PickupAddress != nil {
		so.PickupLocation = so.PickupAddress
	}
	if so.ServiceCost == 0 {
		so.ServiceCost = so.ActualCost
	}
}

type ServiceOrderRequest struct {
	// New multi-service fields
	ServiceCatalogID  string      `json:"service_catalog_id"` // Required for new orders
	ServiceProviderID string      `json:"service_provider_id"` // Optional, can be derived from catalog
	
	// Legacy fields (for backward compatibility)
	BranchID          string      `json:"branch_id"` // Optional if ServiceProviderID is provided
	ServiceType       ServiceType `json:"service_type"` // Optional if ServiceCatalogID is provided
	
	// Legacy iPhone fields (optional, for backward compatibility)
	IPhoneModel       string      `json:"iphone_model"`
	IPhoneColor       string      `json:"iphone_color"`
	IPhoneIMEI        string      `json:"iphone_imei"`
	IPhoneType        string      `json:"iphone_type"`
	
	// Generic item fields (for services that require items)
	ItemModel         string      `json:"item_model"`
	ItemColor         string      `json:"item_color"`
	ItemSerial        string      `json:"item_serial"`
	ItemType          string      `json:"item_type"`
	
	// Service details
	Description       string      `json:"description" validate:"required"`
	Complaint         string      `json:"complaint"`
	
	// Appointment fields (for services that require appointment)
	AppointmentDate   string      `json:"appointment_date"` // ISO 8601 format
	AppointmentTime   string      `json:"appointment_time"` // HH:MM format
	
	// Location fields (optional, depends on service type)
	PickupAddress     string      `json:"pickup_address"`
	PickupLocation    string      `json:"pickup_location"`
	PickupLatitude    *float64    `json:"pickup_latitude"`
	PickupLongitude   *float64    `json:"pickup_longitude"`
	ServiceLocation   string      `json:"service_location"` // Where service will be performed
	
	// Pricing and duration
	EstimatedCost     float64     `json:"estimated_cost"`
	EstimatedDuration int         `json:"estimated_duration"` // in minutes
	
	// Metadata for service-specific fields
	Metadata          JSONB       `json:"metadata"`
}

type ServiceOrderResponse struct {
	ID                uuid.UUID      `json:"id"`
	OrderNumber       string         `json:"order_number"`
	CustomerID        uuid.UUID      `json:"customer_id"`
	Customer          *UserResponse  `json:"customer,omitempty"`
	
	// New multi-service fields
	ServiceCatalogID  *uuid.UUID              `json:"service_catalog_id,omitempty"`
	ServiceCatalog    *ServiceCatalogResponse `json:"service_catalog,omitempty"`
	ServiceProviderID *uuid.UUID              `json:"service_provider_id,omitempty"`
	ServiceProvider   *ServiceProviderResponse `json:"service_provider,omitempty"`
	
	// Legacy fields
	BranchID          *uuid.UUID      `json:"branch_id,omitempty"`
	Branch            *BranchResponse `json:"branch,omitempty"`
	TechnicianID      *uuid.UUID     `json:"technician_id,omitempty"`
	Technician        *UserResponse  `json:"technician,omitempty"`
	CourierID         *uuid.UUID     `json:"courier_id,omitempty"`
	Courier           *UserResponse  `json:"courier,omitempty"`
	
	// Legacy iPhone fields (for backward compatibility)
	IPhoneModel       *string        `json:"iphone_model,omitempty"`
	IPhoneColor       *string        `json:"iphone_color,omitempty"`
	IPhoneIMEI        *string        `json:"iphone_imei,omitempty"`
	IPhoneType        *string        `json:"iphone_type,omitempty"`
	
	// Generic item fields
	ItemModel         *string        `json:"item_model,omitempty"`
	ItemColor         *string        `json:"item_color,omitempty"`
	ItemSerial        *string        `json:"item_serial,omitempty"`
	ItemType          *string        `json:"item_type,omitempty"`
	
	// Service details
	ServiceType       ServiceType    `json:"service_type"`
	ServiceName       string         `json:"service_name"`
	Description       string         `json:"description"`
	Complaint         string         `json:"complaint,omitempty"`
	
	// Appointment
	AppointmentDate   *time.Time     `json:"appointment_date,omitempty"`
	AppointmentTime   *time.Time     `json:"appointment_time,omitempty"`
	
	// Location
	PickupAddress     *string        `json:"pickup_address,omitempty"`
	PickupLocation    *string        `json:"pickup_location,omitempty"`
	PickupLatitude    *float64       `json:"pickup_latitude,omitempty"`
	PickupLongitude   *float64       `json:"pickup_longitude,omitempty"`
	ServiceLocation   string         `json:"service_location"`
	
	// On-demand service
	IsOnDemand        bool           `json:"is_on_demand"`
	CurrentLatitude   *float64       `json:"current_latitude,omitempty"`
	CurrentLongitude  *float64       `json:"current_longitude,omitempty"`
	ETA               *int           `json:"eta,omitempty"` // in minutes
	LastLocationUpdate *time.Time    `json:"last_location_update,omitempty"`
	
	// Status and pricing
	Status            OrderStatus    `json:"status"`
	EstimatedCost     float64        `json:"estimated_cost"`
	ActualCost        float64        `json:"actual_cost"`
	ServiceCost       float64        `json:"service_cost"`
	EstimatedDuration int            `json:"estimated_duration"`
	ActualDuration    int            `json:"actual_duration"`
	
	// Photos
	PickupPhoto       string         `json:"pickup_photo,omitempty"`
	ServicePhoto      string         `json:"service_photo,omitempty"`
	DeliveryPhoto     string         `json:"delivery_photo,omitempty"`
	
	// Additional
	Notes             string         `json:"notes,omitempty"`
	InvoiceNumber     string         `json:"invoice_number,omitempty"`
	TaxAmount         float64        `json:"tax_amount"`
	Metadata          JSONB         `json:"metadata,omitempty"`
	
	// Timestamps
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required"`
	Notes  string      `json:"notes,omitempty"`
}

func (so *ServiceOrder) ToResponse() ServiceOrderResponse {
	resp := ServiceOrderResponse{
		ID:                so.ID,
		OrderNumber:       so.OrderNumber,
		CustomerID:        so.CustomerID,
		ServiceCatalogID:  so.ServiceCatalogID,
		ServiceProviderID: so.ServiceProviderID,
		BranchID:          so.BranchID,
		TechnicianID:      so.TechnicianID,
		CourierID:         so.CourierID,
		IPhoneModel:       so.IPhoneModel,
		IPhoneColor:       so.IPhoneColor,
		IPhoneIMEI:        so.IPhoneIMEI,
		IPhoneType:        so.IPhoneType,
		ItemModel:         so.ItemModel,
		ItemColor:         so.ItemColor,
		ItemSerial:        so.ItemSerial,
		ItemType:          so.ItemType,
		ServiceType:       so.ServiceType,
		ServiceName:       so.ServiceName,
		Description:       so.Description,
		Complaint:         so.Complaint,
		AppointmentDate:   so.AppointmentDate,
		AppointmentTime:   so.AppointmentTime,
		PickupAddress:     so.PickupAddress,
		PickupLocation:    so.PickupLocation,
		PickupLatitude:    so.PickupLatitude,
		PickupLongitude:   so.PickupLongitude,
		ServiceLocation:   so.ServiceLocation,
		IsOnDemand:        so.IsOnDemand,
		CurrentLatitude:   so.CurrentLatitude,
		CurrentLongitude:  so.CurrentLongitude,
		ETA:               so.ETA,
		LastLocationUpdate: so.LastLocationUpdate,
		Status:            so.Status,
		EstimatedCost:     so.EstimatedCost,
		ActualCost:        so.ActualCost,
		ServiceCost:       so.ServiceCost,
		EstimatedDuration: so.EstimatedDuration,
		ActualDuration:    so.ActualDuration,
		PickupPhoto:       so.PickupPhoto,
		ServicePhoto:      so.ServicePhoto,
		DeliveryPhoto:     so.DeliveryPhoto,
		Notes:             so.Notes,
		InvoiceNumber:     so.InvoiceNumber,
		TaxAmount:         so.TaxAmount,
		Metadata:          so.Metadata,
		CreatedAt:         so.CreatedAt,
		UpdatedAt:         so.UpdatedAt,
	}

	// Fill relations if preloaded
	if so.Customer != nil {
		customerResp := so.Customer.ToResponse()
		resp.Customer = &customerResp
	}
	if so.Branch != nil {
		branchResp := so.Branch.ToResponse()
		resp.Branch = &branchResp
	}
	if so.Technician != nil {
		technicianResp := so.Technician.ToResponse()
		resp.Technician = &technicianResp
	}
	if so.Courier != nil {
		courierResp := so.Courier.ToResponse()
		resp.Courier = &courierResp
	}
	if so.ServiceCatalog != nil {
		catalogResp := so.ServiceCatalog.ToResponse()
		resp.ServiceCatalog = &catalogResp
	}
	if so.ServiceProvider != nil {
		providerResp := so.ServiceProvider.ToResponse()
		resp.ServiceProvider = &providerResp
	}

	return resp
}
