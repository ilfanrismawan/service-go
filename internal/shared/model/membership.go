package model

import (
	"fmt"
	"time"
	userDTO "service/internal/users/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	User               userDTO.User             `json:"user" gorm:"foreignKey:UserID"`
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

// MembershipTierConfig represents configuration for each membership tier
type MembershipTierConfig struct {
	Tier               MembershipTier `json:"tier"`
	DiscountPercentage float64        `json:"discount_percentage"`
	MonthlyPrice       float64        `json:"monthly_price"`
	YearlyPrice        float64        `json:"yearly_price"`
	PointsMultiplier   float64        `json:"points_multiplier"`
	MaxFreeServices    int            `json:"max_free_services"`
	MaxFreePickups     int            `json:"max_free_pickups"`
	PrioritySupport    bool           `json:"priority_support"`
	ExtendedWarranty   int            `json:"extended_warranty_days"`
	ExclusiveOffers    bool           `json:"exclusive_offers"`
	FreeDiagnostics    bool           `json:"free_diagnostics"`
	Benefits           []string       `json:"benefits"`
}

// GetMembershipTierConfigs returns the configuration for all membership tiers
func GetMembershipTierConfigs() []MembershipTierConfig {
	return []MembershipTierConfig{
		{
			Tier:               MembershipTierBasic,
			DiscountPercentage: 5.0,
			MonthlyPrice:       50000,  // 50k IDR per month
			YearlyPrice:        500000, // 500k IDR per year (2 months free)
			PointsMultiplier:   1.0,
			MaxFreeServices:    0,
			MaxFreePickups:     0,
			PrioritySupport:    false,
			ExtendedWarranty:   0,
			ExclusiveOffers:    false,
			FreeDiagnostics:    false,
			Benefits: []string{
				"5% discount on all services",
				"Basic support",
				"Mobile app access",
				"Order tracking",
			},
		},
		{
			Tier:               MembershipTierPremium,
			DiscountPercentage: 10.0,
			MonthlyPrice:       100000,  // 100k IDR per month
			YearlyPrice:        1000000, // 1M IDR per year (2 months free)
			PointsMultiplier:   1.5,
			MaxFreeServices:    1,
			MaxFreePickups:     2,
			PrioritySupport:    true,
			ExtendedWarranty:   30,
			ExclusiveOffers:    false,
			FreeDiagnostics:    true,
			Benefits: []string{
				"10% discount on all services",
				"Priority support",
				"1 free service per month",
				"2 free pickups per month",
				"Free diagnostics",
				"30 days extended warranty",
				"Priority queue",
				"Mobile app premium features",
			},
		},
		{
			Tier:               MembershipTierVIP,
			DiscountPercentage: 15.0,
			MonthlyPrice:       200000,  // 200k IDR per month
			YearlyPrice:        2000000, // 2M IDR per year (2 months free)
			PointsMultiplier:   2.0,
			MaxFreeServices:    3,
			MaxFreePickups:     5,
			PrioritySupport:    true,
			ExtendedWarranty:   60,
			ExclusiveOffers:    true,
			FreeDiagnostics:    true,
			Benefits: []string{
				"15% discount on all services",
				"VIP support with dedicated manager",
				"3 free services per month",
				"5 free pickups per month",
				"Free diagnostics",
				"60 days extended warranty",
				"Exclusive offers & promotions",
				"Priority queue & fast track",
				"Free device cleaning",
				"Concierge service",
			},
		},
		{
			Tier:               MembershipTierElite,
			DiscountPercentage: 20.0,
			MonthlyPrice:       350000,  // 350k IDR per month
			YearlyPrice:        3500000, // 3.5M IDR per year (2 months free)
			PointsMultiplier:   3.0,
			MaxFreeServices:    5,
			MaxFreePickups:     10,
			PrioritySupport:    true,
			ExtendedWarranty:   90,
			ExclusiveOffers:    true,
			FreeDiagnostics:    true,
			Benefits: []string{
				"20% discount on all services",
				"Elite support with personal assistant",
				"5 free services per month",
				"10 free pickups per month",
				"Free diagnostics",
				"90 days extended warranty",
				"Exclusive offers & early access",
				"Highest priority queue",
				"Free device cleaning & maintenance",
				"White-glove concierge service",
				"Home service visits",
				"Exclusive events & workshops",
				"Lifetime warranty on repairs",
			},
		},
	}
}

// CalculateDiscount calculates the discount amount based on membership tier
func (m *Membership) CalculateDiscount(amount float64) float64 {
	return amount * (m.DiscountPercentage / 100.0)
}

// CalculatePoints calculates points earned from an order
func (m *Membership) CalculatePoints(amount float64) int64 {
	config := GetMembershipTierConfig(m.Tier)
	return int64(amount * config.PointsMultiplier / 1000) // 1 point per 1000 IDR
}

// GetMembershipTierConfig returns the configuration for a specific tier
func GetMembershipTierConfig(tier MembershipTier) MembershipTierConfig {
	configs := GetMembershipTierConfigs()
	for _, config := range configs {
		if config.Tier == tier {
			return config
		}
	}
	return configs[0] // Return bronze as default
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
func (m *Membership) GetCurrentPrice() float64 {
	config := GetMembershipTierConfig(m.Tier)
	if m.SubscriptionType == SubscriptionTypeYearly {
		return config.YearlyPrice
	}
	return config.MonthlyPrice
}

// ToResponse converts Membership to MembershipResponse
func (m *Membership) ToResponse() MembershipResponse {
	config := GetMembershipTierConfig(m.Tier)

	return MembershipResponse{
		ID:                 m.ID,
		UserID:             m.UserID,
		Tier:               m.Tier,
		Status:             m.Status,
		SubscriptionType:   m.SubscriptionType,
		DiscountPercentage: m.DiscountPercentage,
		Points:             m.Points,
		TotalSpent:         m.TotalSpent,
		OrdersCount:        m.OrdersCount,
		MonthlyPrice:       m.MonthlyPrice,
		YearlyPrice:        m.YearlyPrice,
		CurrentPrice:       m.CurrentPrice,
		JoinedAt:           m.JoinedAt,
		ExpiresAt:          m.ExpiresAt,
		LastOrderAt:        m.LastOrderAt,
		NextBillingDate:    m.NextBillingDate,
		AutoRenew:          m.AutoRenew,
		TrialEndsAt:        m.TrialEndsAt,
		Benefits:           config.Benefits,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

// MembershipResponse represents the response format for membership
type MembershipResponse struct {
	ID                 uuid.UUID        `json:"id"`
	UserID             uuid.UUID        `json:"user_id"`
	Tier               MembershipTier   `json:"tier"`
	Status             MembershipStatus `json:"status"`
	SubscriptionType   SubscriptionType `json:"subscription_type"`
	DiscountPercentage float64          `json:"discount_percentage"`
	Points             int64            `json:"points"`
	TotalSpent         float64          `json:"total_spent"`
	OrdersCount        int64            `json:"orders_count"`
	MonthlyPrice       float64          `json:"monthly_price"`
	YearlyPrice        float64          `json:"yearly_price"`
	CurrentPrice       float64          `json:"current_price"`
	JoinedAt           time.Time        `json:"joined_at"`
	ExpiresAt          *time.Time       `json:"expires_at,omitempty"`
	LastOrderAt        *time.Time       `json:"last_order_at,omitempty"`
	NextBillingDate    *time.Time       `json:"next_billing_date,omitempty"`
	AutoRenew          bool             `json:"auto_renew"`
	TrialEndsAt        *time.Time       `json:"trial_ends_at,omitempty"`
	Benefits           []string         `json:"benefits"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

// MembershipRequest represents the request payload for membership operations
type MembershipRequest struct {
	Tier             MembershipTier   `json:"tier" validate:"required,oneof=basic premium vip elite"`
	SubscriptionType SubscriptionType `json:"subscription_type" validate:"required,oneof=monthly yearly"`
	Status           MembershipStatus `json:"status,omitempty" validate:"omitempty,oneof=active expired suspended cancelled trial"`
	AutoRenew        *bool            `json:"auto_renew,omitempty"`
}

// SubscriptionRequest represents the request payload for subscription operations
type SubscriptionRequest struct {
	Tier             MembershipTier   `json:"tier" validate:"required,oneof=basic premium vip elite"`
	SubscriptionType SubscriptionType `json:"subscription_type" validate:"required,oneof=monthly yearly"`
	PaymentMethod    string           `json:"payment_method" validate:"required"`
	AutoRenew        bool             `json:"auto_renew"`
}

// CancelSubscriptionRequest represents the request payload for cancelling subscription
type CancelSubscriptionRequest struct {
	Reason string `json:"reason,omitempty"`
}

// MembershipOrderResult represents the result of processing an order with membership
type MembershipOrderResult struct {
	HasMembership         bool           `json:"has_membership"`
	OriginalAmount        float64        `json:"original_amount"`
	DiscountAmount        float64        `json:"discount_amount"`
	FinalAmount           float64        `json:"final_amount"`
	PointsEarned          int64          `json:"points_earned"`
	DiscountPercentage    float64        `json:"discount_percentage"`
	Tier                  MembershipTier `json:"tier,omitempty"`
	FreeServiceUsed       bool           `json:"free_service_used,omitempty"`
	FreePickupUsed        bool           `json:"free_pickup_used,omitempty"`
	RemainingFreeServices int            `json:"remaining_free_services,omitempty"`
	RemainingFreePickups  int            `json:"remaining_free_pickups,omitempty"`
}

// MembershipPointsResult represents the result of redeeming points
type MembershipPointsResult struct {
	PointsRedeemed  int64   `json:"points_redeemed"`
	RemainingPoints int64   `json:"remaining_points"`
	DiscountValue   float64 `json:"discount_value"`
}

// SubscriptionResult represents the result of subscription operations
type SubscriptionResult struct {
	MembershipID     uuid.UUID        `json:"membership_id"`
	Tier             MembershipTier   `json:"tier"`
	SubscriptionType SubscriptionType `json:"subscription_type"`
	Status           MembershipStatus `json:"status"`
	Price            float64          `json:"price"`
	NextBillingDate  time.Time        `json:"next_billing_date"`
	PaymentURL       string           `json:"payment_url,omitempty"`
	Message          string           `json:"message"`
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

// GetMembershipUsageKey returns a unique key for membership usage tracking
func (mu *MembershipUsage) GetMembershipUsageKey() string {
	return fmt.Sprintf("%s_%d_%d", mu.MembershipID.String(), mu.Year, mu.Month)
}
