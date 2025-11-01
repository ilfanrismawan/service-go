package core

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Warranty represents warranty tracking for service orders
type Warranty struct {
	ID               uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderID          uuid.UUID   `json:"order_id" gorm:"type:uuid;not null"`
	Order            ServiceOrder `json:"order" gorm:"foreignKey:OrderID"`
	WarrantyDays     int         `json:"warranty_days" gorm:"not null"` // Garansi dalam hari
	StartDate        time.Time   `json:"start_date" gorm:"not null"`
	EndDate          time.Time   `json:"end_date" gorm:"not null"`
	IsActive         bool        `json:"is_active" gorm:"default:true"`
	NotificationSent bool       `json:"notification_sent" gorm:"default:false"` // Notifikasi sebelum habis sudah dikirim
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for Warranty
func (Warranty) TableName() string {
	return "warranties"
}

// IsExpired checks if warranty has expired
func (w *Warranty) IsExpired() bool {
	return time.Now().After(w.EndDate)
}

// DaysRemaining returns days remaining until warranty expires
func (w *Warranty) DaysRemaining() int {
	if w.IsExpired() {
		return 0
	}
	remaining := time.Until(w.EndDate)
	return int(remaining.Hours() / 24)
}

// ShouldNotify checks if warranty notification should be sent (7 days before expiry)
func (w *Warranty) ShouldNotify() bool {
	if w.NotificationSent {
		return false
	}
	daysRemaining := w.DaysRemaining()
	return daysRemaining > 0 && daysRemaining <= 7
}

