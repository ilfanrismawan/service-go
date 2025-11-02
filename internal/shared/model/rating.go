package model

import (
	"time"
	userDTO "service/internal/users/dto"
	branchDTO "service/internal/branches/dto"
	orderDTO "service/internal/orders/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Rating represents a customer rating and review
type Rating struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderID     uuid.UUID   `json:"order_id" gorm:"type:uuid;not null"`
	Order       orderDTO.ServiceOrder `json:"order" gorm:"foreignKey:OrderID"`
	CustomerID  uuid.UUID   `json:"customer_id" gorm:"type:uuid;not null"`
	Customer    userDTO.User        `json:"customer" gorm:"foreignKey:CustomerID"`
	BranchID    *uuid.UUID  `json:"branch_id,omitempty" gorm:"type:uuid"`
	Branch      *branchDTO.Branch     `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	TechnicianID *uuid.UUID  `json:"technician_id,omitempty" gorm:"type:uuid"`
	Technician   *userDTO.User       `json:"technician,omitempty" gorm:"foreignKey:TechnicianID"`
	Rating      int         `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"` // 1-5 stars
	Review      string      `json:"review" gorm:"type:text"`
	IsPublic    bool        `json:"is_public" gorm:"default:true"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for Rating
func (Rating) TableName() string {
	return "ratings"
}

// RatingRequest represents the request payload for creating a rating
type RatingRequest struct {
	OrderID     string `json:"order_id" validate:"required"`
	Rating      int    `json:"rating" validate:"required,min=1,max=5"`
	Review      string `json:"review,omitempty"`
	IsPublic    bool   `json:"is_public" gorm:"default:true"`
}

// RatingResponse represents the response payload for rating data
type RatingResponse struct {
	ID          uuid.UUID   `json:"id"`
	OrderID     uuid.UUID   `json:"order_id"`
	CustomerID  uuid.UUID   `json:"customer_id"`
	Customer    userDTO.UserResponse `json:"customer"`
	BranchID    *uuid.UUID  `json:"branch_id,omitempty"`
	Branch      *branchDTO.BranchResponse `json:"branch,omitempty"`
	TechnicianID *uuid.UUID  `json:"technician_id,omitempty"`
	Technician   *userDTO.UserResponse `json:"technician,omitempty"`
	Rating      int         `json:"rating"`
	Review      string      `json:"review,omitempty"`
	IsPublic    bool        `json:"is_public"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// ToResponse converts Rating to RatingResponse
func (r *Rating) ToResponse() RatingResponse {
	response := RatingResponse{
		ID:         r.ID,
		OrderID:    r.OrderID,
		CustomerID: r.CustomerID,
		Customer:   r.Customer.ToResponse(),
		Rating:     r.Rating,
		Review:     r.Review,
		IsPublic:   r.IsPublic,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}

	if r.BranchID != nil {
		response.BranchID = r.BranchID
		if r.Branch != nil {
			branchResp := r.Branch.ToResponse()
			response.Branch = &branchResp
		}
	}

	if r.TechnicianID != nil {
		response.TechnicianID = r.TechnicianID
		if r.Technician != nil {
			techResp := r.Technician.ToResponse()
			response.Technician = &techResp
		}
	}

	return response
}

// AverageRating represents average rating statistics
type AverageRating struct {
	AverageRating float64 `json:"average_rating"`
	TotalRatings  int64   `json:"total_ratings"`
	Rating5       int64   `json:"rating_5"`
	Rating4       int64   `json:"rating_4"`
	Rating3       int64   `json:"rating_3"`
	Rating2       int64   `json:"rating_2"`
	Rating1       int64   `json:"rating_1"`
}

