package dto

import (
	"time"

	"github.com/google/uuid"

	userDto "service-go/internal/modules/users/dto"
	branchDto "service-go/internal/modules/branches/dto"
)

// RatingRequest represents the request payload for creating a rating
type RatingRequest struct {
	OrderID  string `json:"order_id" validate:"required"`
	Rating   int    `json:"rating" validate:"required,min=1,max=5"`
	Review   string `json:"review,omitempty"`
	IsPublic bool   `json:"is_public" gorm:"default:true"`
}

// RatingResponse represents the response payload for rating data
type RatingResponse struct {
	ID           uuid.UUID                `json:"id"`
	OrderID      uuid.UUID                `json:"order_id"`
	CustomerID   uuid.UUID                `json:"customer_id"`
	Customer     userDto.UserResponse     `json:"customer"`
	BranchID     *uuid.UUID               `json:"branch_id,omitempty"`
	Branch       *branchDto.BranchResponse `json:"branch,omitempty"`
	TechnicianID *uuid.UUID               `json:"technician_id,omitempty"`
	Technician   *userDto.UserResponse    `json:"technician,omitempty"`
	Rating       int                     `json:"rating"`
	Review       string                  `json:"review,omitempty"`
	IsPublic     bool                    `json:"is_public"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
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

// ToRatingResponse converts Rating entity to RatingResponse DTO
func ToRatingResponse(r *ratingEntity.Rating) RatingResponse {
	response := RatingResponse{
		ID:         r.ID,
		OrderID:    r.OrderID,
		CustomerID: r.CustomerID,
		Rating:     r.Rating,
		Review:     r.Review,
		IsPublic:   r.IsPublic,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}

	if r.BranchID != nil {
		response.BranchID = r.BranchID
		if r.Branch != nil {
			branchResp := branchDto.ToBranchResponse(r.Branch)
			response.Branch = &branchResp
		}
	}

	if r.TechnicianID != nil {
		response.TechnicianID = r.TechnicianID
		if r.Technician != nil {
			techResp := userDto.ToUserResponse(r.Technician)
			response.Technician = &techResp
		}
	}

	return response
}

