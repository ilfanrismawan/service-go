package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	branchEntity "service-go/internal/modules/branches/entity"
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

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

// User represents a user entity stored in the database
type User struct {
	ID          uuid.UUID            `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email       string               `gorm:"uniqueIndex;not null"`
	Password    string               `gorm:"not null"`
	FullName    string               `gorm:"not null"`
	Name        string               `gorm:"-"` // computed, maps to FullName
	Phone       string               `gorm:"not null"`
	AvatarURL   string               `json:"avatar_url"`
	Role        UserRole             `json:"role" gorm:"type:varchar(50);not null"`
	Status      UserStatus           `json:"status" gorm:"type:varchar(50);default:'active'"`
	LastLoginAt *time.Time           `json:"last_login_at"`
	BranchID    *uuid.UUID           `gorm:"type:uuid;references:id"`
	Branch      *branchEntity.Branch  `gorm:"foreignKey:BranchID"`
	IsActive    bool                 `gorm:"default:true"`
	FCMToken    string               `json:"fcm_token,omitempty" gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt       `gorm:"index"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// SetName maps FullName → Name (used in domain logic)
func (u *User) SetName() {
	u.Name = u.FullName
}

