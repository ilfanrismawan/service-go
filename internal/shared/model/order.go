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
	BranchID          uuid.UUID   `gorm:"type:uuid;not null"`
	TechnicianID      *uuid.UUID  `gorm:"type:uuid"`
	CourierID         *uuid.UUID  `gorm:"type:uuid"`
	IPhoneModel       string      `gorm:"not null"`
	IPhoneColor       string      `gorm:"not null"`
	IPhoneIMEI        string      `gorm:"not null"`
	IPhoneType        string      `gorm:"not null"`
	ServiceType       ServiceType `gorm:"not null"`
	Description       string      `gorm:"not null"`
	Complaint         string      `gorm:"not null"`
	PickupAddress     string      `gorm:"not null"`
	PickupLocation    string      `gorm:"not null"`
	PickupLatitude    float64     `gorm:"not null"`
	PickupLongitude   float64     `gorm:"not null"`
	Status            OrderStatus `gorm:"not null;default:'pending_pickup'"`
	EstimatedCost     float64     `gorm:"default:0"`
	ActualCost        float64     `gorm:"default:0"`
	ServiceCost       float64     `gorm:"default:0"`
	EstimatedDuration int         `gorm:"default:0"`
	ActualDuration    int         `gorm:"default:0"`
	PickupPhoto       string
	ServicePhoto      string
	DeliveryPhoto     string
	Notes             string
	InvoiceNumber     string  `gorm:"uniqueIndex"`
	TaxAmount         float64 `gorm:"default:0"`
	CustomerIDCard    string
	CustomerNPWP      string
	TermsAccepted     bool `gorm:"default:false"`
	PrivacyAccepted   bool `gorm:"default:false"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (ServiceOrder) TableName() string {
	return "service_orders"
}

func (so *ServiceOrder) SetAliasFields() {
	so.IPhoneType = so.IPhoneModel
	so.Complaint = so.Description
	so.PickupLocation = so.PickupAddress
	so.ServiceCost = so.ActualCost
}

type ServiceOrderRequest struct {
	IPhoneModel       string      `json:"iphone_model" validate:"required"`
	IPhoneColor       string      `json:"iphone_color" validate:"required"`
	IPhoneIMEI        string      `json:"iphone_imei" validate:"required"`
	IPhoneType        string      `json:"iphone_type" validate:"required"`
	ServiceType       ServiceType `json:"service_type" validate:"required"`
	Description       string      `json:"description" validate:"required"`
	Complaint         string      `json:"complaint" validate:"required"`
	PickupAddress     string      `json:"pickup_address" validate:"required"`
	PickupLocation    string      `json:"pickup_location" validate:"required"`
	PickupLatitude    float64     `json:"pickup_latitude" validate:"required"`
	PickupLongitude   float64     `json:"pickup_longitude" validate:"required"`
	EstimatedCost     float64     `json:"estimated_cost"`
	EstimatedDuration int         `json:"estimated_duration"`
	BranchID          string      `json:"branch_id" validate:"required"`
}

type ServiceOrderResponse struct {
	ID                uuid.UUID      `json:"id"`
	OrderNumber       string         `json:"order_number"`
	CustomerID        uuid.UUID      `json:"customer_id"`
	Customer          UserResponse   `json:"customer"`
	BranchID          uuid.UUID      `json:"branch_id"`
	Branch            BranchResponse `json:"branch"`
	TechnicianID      *uuid.UUID     `json:"technician_id,omitempty"`
	Technician        *UserResponse  `json:"technician,omitempty"`
	CourierID         *uuid.UUID     `json:"courier_id,omitempty"`
	Courier           *UserResponse  `json:"courier,omitempty"`
	IPhoneModel       string         `json:"iphone_model"`
	IPhoneColor       string         `json:"iphone_color"`
	IPhoneIMEI        string         `json:"iphone_imei"`
	ServiceType       ServiceType    `json:"service_type"`
	Description       string         `json:"description"`
	PickupAddress     string         `json:"pickup_address"`
	PickupLatitude    float64        `json:"pickup_latitude"`
	PickupLongitude   float64        `json:"pickup_longitude"`
	Status            OrderStatus    `json:"status"`
	EstimatedCost     float64        `json:"estimated_cost"`
	ActualCost        float64        `json:"actual_cost"`
	EstimatedDuration int            `json:"estimated_duration"`
	ActualDuration    int            `json:"actual_duration"`
	PickupPhoto       string         `json:"pickup_photo,omitempty"`
	ServicePhoto      string         `json:"service_photo,omitempty"`
	DeliveryPhoto     string         `json:"delivery_photo,omitempty"`
	Notes             string         `json:"notes,omitempty"`
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
		BranchID:          so.BranchID,
		TechnicianID:      so.TechnicianID,
		CourierID:         so.CourierID,
		IPhoneModel:       so.IPhoneModel,
		IPhoneColor:       so.IPhoneColor,
		IPhoneIMEI:        so.IPhoneIMEI,
		ServiceType:       so.ServiceType,
		Description:       so.Description,
		PickupAddress:     so.PickupAddress,
		PickupLatitude:    so.PickupLatitude,
		PickupLongitude:   so.PickupLongitude,
		Status:            so.Status,
		EstimatedCost:     so.EstimatedCost,
		ActualCost:        so.ActualCost,
		EstimatedDuration: so.EstimatedDuration,
		ActualDuration:    so.ActualDuration,
		PickupPhoto:       so.PickupPhoto,
		ServicePhoto:      so.ServicePhoto,
		DeliveryPhoto:     so.DeliveryPhoto,
		Notes:             so.Notes,
		CreatedAt:         so.CreatedAt,
		UpdatedAt:         so.UpdatedAt,
	}

	// opsional: isi relasi jika sudah dipreload
	if so.CustomerID != uuid.Nil {
		resp.CustomerID = so.CustomerID
	}
	return resp
}
