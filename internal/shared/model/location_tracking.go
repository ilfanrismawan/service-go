package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LocationTracking represents real-time location tracking history for orders
type LocationTracking struct {
	ID              uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderID         uuid.UUID   `json:"order_id" gorm:"type:uuid;not null;index"`
	Order           *ServiceOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	UserID          uuid.UUID   `json:"user_id" gorm:"type:uuid;not null;index"` // Courier/Provider ID
	User            *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Latitude        float64     `json:"latitude" gorm:"type:decimal(10,6);not null"`
	Longitude       float64     `json:"longitude" gorm:"type:decimal(10,6);not null"`
	Accuracy        float64     `json:"accuracy" gorm:"type:decimal(10,2);default:0"` // GPS accuracy in meters
	Speed           float64     `json:"speed" gorm:"type:decimal(10,2);default:0"` // Speed in km/h
	Heading         float64     `json:"heading" gorm:"type:decimal(5,2);default:0"` // Direction in degrees (0-360)
	Timestamp       time.Time   `json:"timestamp" gorm:"not null;index"`
	CreatedAt       time.Time   `json:"created_at"`
}

func (LocationTracking) TableName() string {
	return "location_tracking"
}

// CurrentLocation represents current location of courier/provider for an order
type CurrentLocation struct {
	ID              uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderID         uuid.UUID   `json:"order_id" gorm:"type:uuid;not null;uniqueIndex"`
	Order           *ServiceOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	UserID          uuid.UUID   `json:"user_id" gorm:"type:uuid;not null;index"` // Courier/Provider ID
	User            *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Latitude        float64     `json:"latitude" gorm:"type:decimal(10,6);not null"`
	Longitude       float64     `json:"longitude" gorm:"type:decimal(10,6);not null"`
	Accuracy        float64     `json:"accuracy" gorm:"type:decimal(10,2);default:0"`
	Speed           float64     `json:"speed" gorm:"type:decimal(10,2);default:0"`
	Heading         float64     `json:"heading" gorm:"type:decimal(5,2);default:0"`
	ETA             int         `json:"eta" gorm:"default:0"` // Estimated time of arrival in minutes
	Distance        float64     `json:"distance" gorm:"type:decimal(10,2);default:0"` // Distance to destination in km
	UpdatedAt       time.Time   `json:"updated_at" gorm:"not null;index"`
	CreatedAt       time.Time   `json:"created_at"`
}

func (CurrentLocation) TableName() string {
	return "current_locations"
}

// LocationUpdateRequest represents request payload for updating location
type LocationUpdateRequest struct {
	Latitude    float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude   float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Accuracy    float64 `json:"accuracy,omitempty"` // GPS accuracy in meters
	Speed       float64 `json:"speed,omitempty"` // Speed in km/h
	Heading     float64 `json:"heading,omitempty"` // Direction in degrees (0-360)
}

// LocationUpdateResponse represents response payload for location update
type LocationUpdateResponse struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	UserID      uuid.UUID `json:"user_id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Accuracy    float64   `json:"accuracy"`
	Speed       float64   `json:"speed"`
	Heading     float64   `json:"heading"`
	ETA         int       `json:"eta"` // in minutes
	Distance    float64   `json:"distance"` // in km
	UpdatedAt   time.Time `json:"updated_at"`
}

func (cl *CurrentLocation) ToResponse() LocationUpdateResponse {
	return LocationUpdateResponse{
		ID:          cl.ID,
		OrderID:     cl.OrderID,
		UserID:      cl.UserID,
		Latitude:    cl.Latitude,
		Longitude:   cl.Longitude,
		Accuracy:    cl.Accuracy,
		Speed:       cl.Speed,
		Heading:     cl.Heading,
		ETA:         cl.ETA,
		Distance:    cl.Distance,
		UpdatedAt:   cl.UpdatedAt,
	}
}

// LocationHistoryResponse represents location history response
type LocationHistoryResponse struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	UserID      uuid.UUID `json:"user_id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Accuracy    float64   `json:"accuracy"`
	Speed       float64   `json:"speed"`
	Heading     float64   `json:"heading"`
	Timestamp   time.Time `json:"timestamp"`
	CreatedAt   time.Time `json:"created_at"`
}

func (lt *LocationTracking) ToHistoryResponse() LocationHistoryResponse {
	return LocationHistoryResponse{
		ID:          lt.ID,
		OrderID:     lt.OrderID,
		UserID:      lt.UserID,
		Latitude:    lt.Latitude,
		Longitude:   lt.Longitude,
		Accuracy:    lt.Accuracy,
		Speed:       lt.Speed,
		Heading:     lt.Heading,
		Timestamp:   lt.Timestamp,
		CreatedAt:   lt.CreatedAt,
	}
}

// ETACalculationRequest represents request for ETA calculation
type ETACalculationRequest struct {
	CurrentLatitude  float64 `json:"current_latitude" validate:"required"`
	CurrentLongitude float64 `json:"current_longitude" validate:"required"`
	DestinationLatitude  float64 `json:"destination_latitude" validate:"required"`
	DestinationLongitude float64 `json:"destination_longitude" validate:"required"`
	Speed            float64 `json:"speed,omitempty"` // Current speed in km/h
}

// ETAResponse represents ETA calculation response
type ETAResponse struct {
	ETA         int     `json:"eta"` // in minutes
	Distance    float64 `json:"distance"` // in km
	EstimatedArrival time.Time `json:"estimated_arrival"`
}

