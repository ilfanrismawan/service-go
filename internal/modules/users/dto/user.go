package dto

import (
	"time"

	"github.com/google/uuid"

	userEntity "service-go/internal/modules/users/entity"
	branchDto "service-go/internal/modules/branches/dto"
)

// UserRequest defines request payload for creating or updating user
type UserRequest struct {
	Email    string            `json:"email" validate:"required,email"`
	Password string            `json:"password" validate:"required,min=6"`
	FullName string            `json:"full_name" validate:"required"`
	Phone    string            `json:"phone" validate:"required"`
	Role     userEntity.UserRole `json:"role" validate:"required"`
	BranchID *string          `json:"branch_id,omitempty"`
}

// UserResponse defines how user data is returned to the client
type UserResponse struct {
	ID          uuid.UUID                `json:"id"`
	Email       string                   `json:"email"`
	FullName    string                   `json:"full_name"`
	Phone       string                   `json:"phone"`
	Role        userEntity.UserRole      `json:"role"`
	AvatarURL   string                   `json:"avatar_url"`
	Status      userEntity.UserStatus    `json:"status"`
	LastLoginAt *time.Time               `json:"last_login_at,omitempty"`
	BranchID    *uuid.UUID               `json:"branch_id,omitempty"`
	Branch      *branchDto.BranchResponse `json:"branch,omitempty"`
	IsActive    bool                     `json:"is_active"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

// UserUpdateRequest defines request payload for updating user
type UserUpdateRequest struct {
	FullName  string              `json:"full_name,omitempty"`
	Email     string              `json:"email,omitempty" validate:"omitempty,email"`
	Phone     string              `json:"phone,omitempty"`
	Password  string              `json:"password,omitempty"`
	AvatarURL string              `json:"avatar_url,omitempty"`
	Status    userEntity.UserStatus `json:"status,omitempty"`
}

// ToResponse converts User entity to UserResponse DTO
func ToUserResponse(u *userEntity.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		FullName:    u.FullName,
		Email:       u.Email,
		Phone:       u.Phone,
		AvatarURL:   u.AvatarURL,
		Role:        u.Role,
		Status:      u.Status,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

