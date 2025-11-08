package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userEntity "service-go/internal/modules/users/entity"
)

// MembershipTier represents the membership tier level
type MembershipTier string

const (
	MembershipTierBasic   MembershipTier = "basic"
	MembershipTierPremium MembershipTier = "premium"
	MembershipTierVIP     MembershipTier = "vip"
	MembershipTierElite   MembershipTier = "elite"
)

// MembershipStatus represents the membership status
type MembershipStatus string

const (
	MembershipStatusActive    MembershipStatus = "active"
	MembershipStatusExpired   MembershipStatus = "expired"
	MembershipStatusSuspended MembershipStatus = "suspended"
	MembershipStatusCancelled MembershipStatus = "cancelled"
	MembershipStatusTrial     MembershipStatus = "trial"
)

// SubscriptionType represents the subscription billing type
type SubscriptionType string

const (
	SubscriptionTypeMonthly SubscriptionType = "monthly"
	SubscriptionTypeYearly  SubscriptionType = "yearly"
)

// Membership represents a user's membership information
type Membership struct {
	ID                 uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID             uuid.UUID        `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User               userEntity.User   `json:"user" gorm:"foreignKey:UserID"`
	Tier               MembershipTier   `json:"tier" gorm:"not null;default:'basic'"`
	Status             MembershipStatus `json:"status" gorm:"not null;default:'trial'"`
	SubscriptionType   SubscriptionType `json:"subscription_type" gorm:"not null;default:'monthly'"`
	DiscountPercentage float64          `json:"discount_percentage" gorm:"default:0"`
	Points             int64            `json:"points" gorm:"default:0"`
	TotalSpent         float64          `json:"total_spent" gorm:"default:0"`
	OrdersCount        int64            `json:"orders_count" gorm:"default:0"`
	MonthlyPrice       float64          `json:"monthly_price" gorm:"default:0"`
	YearlyPrice        float64          `json:"yearly_price" gorm:"default:0"`
	CurrentPrice       float64          `json:"current_price" gorm:"default:0"`
	JoinedAt           time.Time        `json:"joined_at" gorm:"not null"`
	ExpiresAt          *time.Time       `json:"expires_at,omitempty"`
	LastOrderAt        *time.Time       `json:"last_order_at,omitempty"`
	NextBillingDate    *time.Time       `json:"next_billing_date,omitempty"`
	AutoRenew          bool             `json:"auto_renew" gorm:"default:true"`
	TrialEndsAt        *time.Time       `json:"trial_ends_at,omitempty"`
	CreatedAt          time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"index"`
}

// MembershipUsage represents the usage tracking for membership benefits
type MembershipUsage struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	MembershipID     uuid.UUID  `json:"membership_id" gorm:"type:uuid;not null"`
	Membership       Membership `json:"membership" gorm:"foreignKey:MembershipID"`
	FreeServicesUsed int        `json:"free_services_used" gorm:"default:0"`
	FreePickupsUsed  int        `json:"free_pickups_used" gorm:"default:0"`
	Month            int        `json:"month" gorm:"not null"`
	Year             int        `json:"year" gorm:"not null"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// CalculateDiscount calculates the discount amount based on membership tier
func (m *Membership) CalculateDiscount(amount float64) float64 {
	return amount * (m.DiscountPercentage / 100.0)
}

// CalculatePoints calculates points earned from an order
// Note: This method requires MembershipTierConfig which should be provided from dto package
func (m *Membership) CalculatePoints(amount float64, pointsMultiplier float64) int64 {
	return int64(amount * pointsMultiplier / 1000) // 1 point per 1000 IDR
}

// IsTrialExpired checks if the trial period has expired
func (m *Membership) IsTrialExpired() bool {
	if m.Status != MembershipStatusTrial || m.TrialEndsAt == nil {
		return false
	}
	return time.Now().After(*m.TrialEndsAt)
}

// IsExpired checks if the membership has expired
func (m *Membership) IsExpired() bool {
	if m.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*m.ExpiresAt)
}

// GetNextBillingDate calculates the next billing date
func (m *Membership) GetNextBillingDate() time.Time {
	if m.NextBillingDate != nil {
		return *m.NextBillingDate
	}

	now := time.Now()
	if m.SubscriptionType == SubscriptionTypeYearly {
		return now.AddDate(1, 0, 0)
	}
	return now.AddDate(0, 1, 0)
}

// GetCurrentPrice returns the current subscription price
// Note: This method uses stored prices, config should be provided from dto package
func (m *Membership) GetCurrentPrice() float64 {
	if m.SubscriptionType == SubscriptionTypeYearly {
		return m.YearlyPrice
	}
	return m.MonthlyPrice
}

// GetMembershipUsageKey returns a unique key for membership usage tracking
func (mu *MembershipUsage) GetMembershipUsageKey() string {
	return fmt.Sprintf("%s_%d_%d", mu.MembershipID.String(), mu.Year, mu.Month)
}

