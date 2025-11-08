package dto

import (
	"time"

	"github.com/google/uuid"

	orderEntity "service-go/internal/modules/orders/entity"
	userDto "service-go/internal/modules/users/dto"
	branchDto "service-go/internal/modules/branches/dto"
)

// ServiceOrderRequest represents the request payload for creating a service order
type ServiceOrderRequest struct {
	IPhoneModel       string                `json:"iphone_model" validate:"required"`
	IPhoneColor       string                `json:"iphone_color" validate:"required"`
	IPhoneIMEI        string                `json:"iphone_imei" validate:"required"`
	IPhoneType        string                `json:"iphone_type" validate:"required"`
	ServiceType       orderEntity.ServiceType `json:"service_type" validate:"required"`
	Description       string                `json:"description" validate:"required"`
	Complaint         string                `json:"complaint" validate:"required"`
	PickupAddress     string                `json:"pickup_address" validate:"required"`
	PickupLocation    string                `json:"pickup_location" validate:"required"`
	PickupLatitude    float64               `json:"pickup_latitude" validate:"required"`
	PickupLongitude   float64               `json:"pickup_longitude" validate:"required"`
	EstimatedCost     float64               `json:"estimated_cost"`
	EstimatedDuration int                   `json:"estimated_duration"`
	BranchID          string                `json:"branch_id" validate:"required"`
}

// ServiceOrderResponse represents the response payload for service order data
type ServiceOrderResponse struct {
	ID                uuid.UUID                `json:"id"`
	OrderNumber       string                    `json:"order_number"`
	CustomerID        uuid.UUID                 `json:"customer_id"`
	Customer          userDto.UserResponse       `json:"customer"`
	BranchID          uuid.UUID                 `json:"branch_id"`
	Branch            branchDto.BranchResponse   `json:"branch"`
	TechnicianID      *uuid.UUID                `json:"technician_id,omitempty"`
	Technician        *userDto.UserResponse     `json:"technician,omitempty"`
	CourierID         *uuid.UUID                `json:"courier_id,omitempty"`
	Courier           *userDto.UserResponse     `json:"courier,omitempty"`
	IPhoneModel       string                    `json:"iphone_model"`
	IPhoneColor       string                    `json:"iphone_color"`
	IPhoneIMEI        string                    `json:"iphone_imei"`
	ServiceType       orderEntity.ServiceType   `json:"service_type"`
	Description       string                    `json:"description"`
	PickupAddress     string                    `json:"pickup_address"`
	PickupLatitude    float64                   `json:"pickup_latitude"`
	PickupLongitude   float64                   `json:"pickup_longitude"`
	Status            orderEntity.OrderStatus    `json:"status"`
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

// UpdateOrderStatusRequest represents the request payload for updating order status
type UpdateOrderStatusRequest struct {
	Status orderEntity.OrderStatus `json:"status" validate:"required"`
	Notes  string                  `json:"notes,omitempty"`
}

// ToServiceOrderResponse converts ServiceOrder entity to ServiceOrderResponse DTO
func ToServiceOrderResponse(so *orderEntity.ServiceOrder) ServiceOrderResponse {
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

	return resp
}

