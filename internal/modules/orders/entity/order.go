package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userEntity "service-go/internal/modules/users/entity"
	branchEntity "service-go/internal/modules/branches/entity"
)

// OrderStatus represents the status of an order
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

// ServiceType represents the type of service
type ServiceType string

const (
	ServiceTypeScreenRepair       ServiceType = "screen_repair"
	ServiceTypeBatteryReplacement ServiceType = "battery_replacement"
	ServiceTypeWaterDamage        ServiceType = "water_damage"
	ServiceTypeSoftwareIssue      ServiceType = "software_issue"
	ServiceTypeHardwareRepair     ServiceType = "hardware_repair"
	ServiceTypeOther              ServiceType = "other"
)

// ServiceOrder represents a service order entity
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

