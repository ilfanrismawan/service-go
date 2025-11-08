package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userEntity "service-go/internal/modules/users/entity"
	branchEntity "service-go/internal/modules/branches/entity"
)

// Rating represents a customer rating and review
type Rating struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderID      uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	Order        ServiceOrder   `json:"order" gorm:"foreignKey:OrderID"`
	CustomerID   uuid.UUID      `json:"customer_id" gorm:"type:uuid;not null"`
	Customer     userEntity.User `json:"customer" gorm:"foreignKey:CustomerID"`
	BranchID     *uuid.UUID     `json:"branch_id,omitempty" gorm:"type:uuid"`
	Branch       *branchEntity.Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	TechnicianID *uuid.UUID     `json:"technician_id,omitempty" gorm:"type:uuid"`
	Technician   *userEntity.User `json:"technician,omitempty" gorm:"foreignKey:TechnicianID"`
	Rating       int            `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"` // 1-5 stars
	Review       string         `json:"review" gorm:"type:text"`
	IsPublic     bool           `json:"is_public" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for Rating
func (Rating) TableName() string {
	return "ratings"
}

