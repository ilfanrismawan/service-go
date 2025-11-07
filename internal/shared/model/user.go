package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdminPusat   UserRole = "admin_pusat"
	RoleAdminCabang  UserRole = "admin_cabang"
	RoleKasir        UserRole = "kasir"
	RoleTeknisi      UserRole = "teknisi"
	RoleKurir        UserRole = "kurir"
	RolePelanggan    UserRole = "pelanggan"
	UserRoleProvider UserRole = "provider"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

// User represents a user entity stored in the database
type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email       string     `gorm:"uniqueIndex;not null"`
	Password    string     `gorm:"not null"`
	FullName    string     `gorm:"not null"`
	Name        string     `gorm:"-"` // computed, maps to FullName
	Phone       string     `gorm:"not null"`
	AvatarURL   string     `json:"avatar_url"`
	Role        UserRole   `json:"role" gorm:"type:varchar(50);not null"`
	Status      UserStatus `json:"status" gorm:"type:varchar(50);default:'active'"`
	LastLoginAt *time.Time `json:"last_login_at"`
	BranchID    *uuid.UUID `gorm:"type:uuid;references:id"`
	Branch      *Branch    `gorm:"foreignKey:BranchID"`
	IsActive    bool       `gorm:"default:true"`
	FCMToken    string     `json:"fcm_token,omitempty" gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// SetName maps FullName → Name (used in domain logic)
func (u *User) SetName() {
	u.Name = u.FullName
}

// UserRequest defines request payload for creating or updating user
type UserRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
	FullName string   `json:"full_name" validate:"required"`
	Phone    string   `json:"phone" validate:"required"`
	Role     UserRole `json:"role" validate:"required"`
	BranchID *string  `json:"branch_id,omitempty"`
}

// UserResponse defines how user data is returned to the client
type UserResponse struct {
	ID          uuid.UUID       `json:"id"`
	Email       string          `json:"email"`
	FullName    string          `json:"full_name"`
	Phone       string          `json:"phone"`
	Role        UserRole        `json:"role"`
	AvatarURL   string          `json:"avatar_url"`
	Status      UserStatus      `json:"status"`
	LastLoginAt *time.Time      `json:"last_login_at,omitempty"`
	BranchID    *uuid.UUID      `json:"branch_id,omitempty"`
	Branch      *BranchResponse `json:"branch,omitempty"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type UserUpdateRequest struct {
	FullName  string     `json:"full_name,omitempty"`
	Email     string     `json:"email,omitempty" validate:"omitempty,email"`
	Phone     string     `json:"phone,omitempty"`
	Password  string     `json:"password,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	Status    UserStatus `json:"status,omitempty"`
}

func (u *User) ToResponse() UserResponse {
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
