package core

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdminPusat  UserRole = "admin_pusat"
	RoleAdminCabang UserRole = "admin_cabang"
	RoleKasir       UserRole = "kasir"
	RoleTeknisi     UserRole = "teknisi"
	RoleKurir       UserRole = "kurir"
	RolePelanggan   UserRole = "pelanggan"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FullName  string         `json:"full_name" gorm:"not null"`
	Name      string         `json:"name" gorm:"-"` // Computed field, maps to FullName
	Phone     string         `json:"phone" gorm:"not null"`
	Role      UserRole       `json:"role" gorm:"type:text;not null"`
	BranchID  *uuid.UUID     `json:"branch_id,omitempty" gorm:"type:uuid;references:id"`
	Branch    *Branch        `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	FCMToken  string         `json:"fcm_token,omitempty" gorm:"type:text"` // Firebase Cloud Messaging token
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// UserRequest represents the request payload for creating/updating a user
type UserRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
	FullName string   `json:"full_name" validate:"required"`
	Phone    string   `json:"phone" validate:"required,phone"`
	Role     UserRole `json:"role" validate:"required"`
	BranchID *string  `json:"branch_id,omitempty"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Phone     string     `json:"phone"`
	Role      UserRole   `json:"role"`
	BranchID  *uuid.UUID `json:"branch_id,omitempty"`
	Branch    *Branch    `json:"branch,omitempty"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Role:      u.Role,
		BranchID:  u.BranchID,
		Branch:    u.Branch,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// SetName sets the Name field from FullName
func (u *User) SetName() {
	u.Name = u.FullName
}
