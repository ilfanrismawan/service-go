package dto

import (
	branchDTO "service/internal/branches/dto"
	"service/internal/shared/model"
	userDTO "service/internal/users/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus = model.OrderStatus
type ServiceType = model.ServiceType

// ServiceOrder represents an iPhone service order
type ServiceOrder struct {
	ID                uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderNumber       string           `json:"order_number" gorm:"uniqueIndex;not null"`
	CustomerID        uuid.UUID        `json:"customer_id" gorm:"type:uuid;not null"`
	UserID            uuid.UUID        `json:"user_id" gorm:"type:uuid;not null"` // Alias for CustomerID
	Customer          userDTO.User     `json:"customer" gorm:"foreignKey:CustomerID"`
	BranchID          uuid.UUID        `json:"branch_id" gorm:"type:uuid;not null"`
	Branch            branchDTO.Branch `json:"branch" gorm:"foreignKey:BranchID"`
	TechnicianID      *uuid.UUID       `json:"technician_id,omitempty" gorm:"type:uuid"`
	Technician        *userDTO.User    `json:"technician,omitempty" gorm:"foreignKey:TechnicianID"`
	CourierID         *uuid.UUID       `json:"courier_id,omitempty" gorm:"type:uuid"`
	Courier           *userDTO.User    `json:"courier,omitempty" gorm:"foreignKey:CourierID"`
	IPhoneModel       string           `json:"iphone_model" gorm:"not null"`
	IPhoneColor       string           `json:"iphone_color" gorm:"not null"`
	IPhoneIMEI        string           `json:"iphone_imei" gorm:"not null"`
	IPhoneType        string           `json:"iphone_type" gorm:"not null"` // Alias for IPhoneModel
	ServiceType       ServiceType      `json:"service_type" gorm:"not null"`
	Description       string           `json:"description" gorm:"not null"`
	Complaint         string           `json:"complaint" gorm:"not null"` // Alias for Description
	PickupAddress     string           `json:"pickup_address" gorm:"not null"`
	PickupLocation    string           `json:"pickup_location" gorm:"not null"` // Alias for PickupAddress
	PickupLatitude    float64          `json:"pickup_latitude" gorm:"not null"`
	PickupLongitude   float64          `json:"pickup_longitude" gorm:"not null"`
	Status            OrderStatus      `json:"status" gorm:"not null;default:'pending_pickup'"`
	EstimatedCost     float64          `json:"estimated_cost" gorm:"default:0"`
	ActualCost        float64          `json:"actual_cost" gorm:"default:0"`
	ServiceCost       float64          `json:"service_cost" gorm:"default:0"`       // Alias for ActualCost
	EstimatedDuration int              `json:"estimated_duration" gorm:"default:0"` // in hours
	ActualDuration    int              `json:"actual_duration" gorm:"default:0"`    // in hours
	PickupPhoto       string           `json:"pickup_photo,omitempty"`
	ServicePhoto      string           `json:"service_photo,omitempty"`
	DeliveryPhoto     string           `json:"delivery_photo,omitempty"`
	Notes             string           `json:"notes,omitempty"`
	// Legal & Compliance fields
	InvoiceNumber   string         `json:"invoice_number,omitempty" gorm:"uniqueIndex"` // Nomor faktur pajak
	TaxAmount       float64        `json:"tax_amount" gorm:"default:0"`                 // PPN 11%
	CustomerIDCard  string         `json:"customer_id_card,omitempty"`                  // KTP customer (encrypted)
	CustomerNPWP    string         `json:"customer_npwp,omitempty"`                     // NPWP customer
	TermsAccepted   bool           `json:"terms_accepted" gorm:"default:false"`
	PrivacyAccepted bool           `json:"privacy_accepted" gorm:"default:false"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for ServiceOrder
func (ServiceOrder) TableName() string {
	return "service_orders"
}

// ServiceOrderRequest represents the request payload for creating a service order
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

// ServiceOrderResponse represents the response payload for service order data
type ServiceOrderResponse struct {
	ID                uuid.UUID                `json:"id"`
	OrderNumber       string                   `json:"order_number"`
	CustomerID        uuid.UUID                `json:"customer_id"`
	Customer          userDTO.UserResponse     `json:"customer"`
	BranchID          uuid.UUID                `json:"branch_id"`
	Branch            branchDTO.BranchResponse `json:"branch"`
	TechnicianID      *uuid.UUID               `json:"technician_id,omitempty"`
	Technician        *userDTO.UserResponse    `json:"technician,omitempty"`
	CourierID         *uuid.UUID               `json:"courier_id,omitempty"`
	Courier           *userDTO.UserResponse    `json:"courier,omitempty"`
	IPhoneModel       string                   `json:"iphone_model"`
	IPhoneColor       string                   `json:"iphone_color"`
	IPhoneIMEI        string                   `json:"iphone_imei"`
	ServiceType       ServiceType              `json:"service_type"`
	Description       string                   `json:"description"`
	PickupAddress     string                   `json:"pickup_address"`
	PickupLatitude    float64                  `json:"pickup_latitude"`
	PickupLongitude   float64                  `json:"pickup_longitude"`
	Status            OrderStatus              `json:"status"`
	EstimatedCost     float64                  `json:"estimated_cost"`
	ActualCost        float64                  `json:"actual_cost"`
	EstimatedDuration int                      `json:"estimated_duration"`
	ActualDuration    int                      `json:"actual_duration"`
	PickupPhoto       string                   `json:"pickup_photo,omitempty"`
	ServicePhoto      string                   `json:"service_photo,omitempty"`
	DeliveryPhoto     string                   `json:"delivery_photo,omitempty"`
	Notes             string                   `json:"notes,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

// ToResponse converts ServiceOrder to ServiceOrderResponse
func (so *ServiceOrder) ToResponse() ServiceOrderResponse {
	response := ServiceOrderResponse{
		ID:                so.ID,
		OrderNumber:       so.OrderNumber,
		CustomerID:        so.CustomerID,
		Customer:          so.Customer.ToResponse(),
		BranchID:          so.BranchID,
		Branch:            so.Branch.ToResponse(),
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

	if so.TechnicianID != nil {
		response.TechnicianID = so.TechnicianID
		technicianResponse := so.Technician.ToResponse()
		response.Technician = &technicianResponse
	}

	if so.CourierID != nil {
		response.CourierID = so.CourierID
		courierResponse := so.Courier.ToResponse()
		response.Courier = &courierResponse
	}

	return response
}

// SetAliasFields sets the alias fields from their corresponding main fields
func (so *ServiceOrder) SetAliasFields() {
	so.UserID = so.CustomerID
	so.IPhoneType = so.IPhoneModel
	so.Complaint = so.Description
	so.PickupLocation = so.PickupAddress
	so.ServiceCost = so.ActualCost
}

// UpdateOrderStatusRequest represents the request payload for updating order status
type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required"`
	Notes  string      `json:"notes,omitempty"`
}
