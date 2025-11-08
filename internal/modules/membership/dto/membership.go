package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	membershipEntity "service-go/internal/modules/membership/entity"
)

// MembershipTierConfig represents configuration for each membership tier
type MembershipTierConfig struct {
	Tier               membershipEntity.MembershipTier `json:"tier"`
	DiscountPercentage float64                         `json:"discount_percentage"`
	MonthlyPrice       float64                         `json:"monthly_price"`
	YearlyPrice        float64                         `json:"yearly_price"`
	PointsMultiplier   float64                         `json:"points_multiplier"`
	MaxFreeServices    int                             `json:"max_free_services"`
	MaxFreePickups     int                             `json:"max_free_pickups"`
	PrioritySupport    bool                            `json:"priority_support"`
	ExtendedWarranty   int                             `json:"extended_warranty_days"`
	ExclusiveOffers    bool                            `json:"exclusive_offers"`
	FreeDiagnostics    bool                            `json:"free_diagnostics"`
	Benefits           []string                        `json:"benefits"`
}

// GetMembershipTierConfigs returns the configuration for all membership tiers
func GetMembershipTierConfigs() []MembershipTierConfig {
	return []MembershipTierConfig{
		{
			Tier:               membershipEntity.MembershipTierBasic,
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
			Tier:               membershipEntity.MembershipTierPremium,
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
			Tier:               membershipEntity.MembershipTierVIP,
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
			Tier:               membershipEntity.MembershipTierElite,
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

// GetMembershipTierConfig returns the configuration for a specific tier
func GetMembershipTierConfig(tier membershipEntity.MembershipTier) MembershipTierConfig {
	configs := GetMembershipTierConfigs()
	for _, config := range configs {
		if config.Tier == tier {
			return config
		}
	}
	return configs[0] // Return basic as default
}

// MembershipResponse represents the response format for membership
type MembershipResponse struct {
	ID                 uuid.UUID                        `json:"id"`
	UserID             uuid.UUID                        `json:"user_id"`
	Tier               membershipEntity.MembershipTier  `json:"tier"`
	Status             membershipEntity.MembershipStatus `json:"status"`
	SubscriptionType   membershipEntity.SubscriptionType `json:"subscription_type"`
	DiscountPercentage float64                          `json:"discount_percentage"`
	Points             int64                            `json:"points"`
	TotalSpent         float64                          `json:"total_spent"`
	OrdersCount        int64                            `json:"orders_count"`
	MonthlyPrice       float64                          `json:"monthly_price"`
	YearlyPrice        float64                          `json:"yearly_price"`
	CurrentPrice       float64                          `json:"current_price"`
	JoinedAt           time.Time                        `json:"joined_at"`
	ExpiresAt          *time.Time                       `json:"expires_at,omitempty"`
	LastOrderAt        *time.Time                       `json:"last_order_at,omitempty"`
	NextBillingDate    *time.Time                       `json:"next_billing_date,omitempty"`
	AutoRenew          bool                             `json:"auto_renew"`
	TrialEndsAt        *time.Time                       `json:"trial_ends_at,omitempty"`
	Benefits           []string                         `json:"benefits"`
	CreatedAt          time.Time                        `json:"created_at"`
	UpdatedAt          time.Time                        `json:"updated_at"`
}

// MembershipRequest represents the request payload for membership operations
type MembershipRequest struct {
	Tier             membershipEntity.MembershipTier   `json:"tier" validate:"required,oneof=basic premium vip elite"`
	SubscriptionType membershipEntity.SubscriptionType `json:"subscription_type" validate:"required,oneof=monthly yearly"`
	Status           membershipEntity.MembershipStatus  `json:"status,omitempty" validate:"omitempty,oneof=active expired suspended cancelled trial"`
	AutoRenew        *bool                             `json:"auto_renew,omitempty"`
}

// SubscriptionRequest represents the request payload for subscription operations
type SubscriptionRequest struct {
	Tier             membershipEntity.MembershipTier   `json:"tier" validate:"required,oneof=basic premium vip elite"`
	SubscriptionType membershipEntity.SubscriptionType `json:"subscription_type" validate:"required,oneof=monthly yearly"`
	PaymentMethod    string                            `json:"payment_method" validate:"required"`
	AutoRenew        bool                              `json:"auto_renew"`
}

// CancelSubscriptionRequest represents the request payload for cancelling subscription
type CancelSubscriptionRequest struct {
	Reason string `json:"reason,omitempty"`
}

// MembershipOrderResult represents the result of processing an order with membership
type MembershipOrderResult struct {
	HasMembership         bool                            `json:"has_membership"`
	OriginalAmount        float64                         `json:"original_amount"`
	DiscountAmount        float64                         `json:"discount_amount"`
	FinalAmount           float64                         `json:"final_amount"`
	PointsEarned          int64                           `json:"points_earned"`
	DiscountPercentage    float64                         `json:"discount_percentage"`
	Tier                  membershipEntity.MembershipTier `json:"tier,omitempty"`
	FreeServiceUsed       bool                            `json:"free_service_used,omitempty"`
	FreePickupUsed        bool                            `json:"free_pickup_used,omitempty"`
	RemainingFreeServices int                             `json:"remaining_free_services,omitempty"`
	RemainingFreePickups  int                             `json:"remaining_free_pickups,omitempty"`
}

// MembershipPointsResult represents the result of redeeming points
type MembershipPointsResult struct {
	PointsRedeemed  int64   `json:"points_redeemed"`
	RemainingPoints int64   `json:"remaining_points"`
	DiscountValue   float64 `json:"discount_value"`
}

// SubscriptionResult represents the result of subscription operations
type SubscriptionResult struct {
	MembershipID     uuid.UUID                        `json:"membership_id"`
	Tier             membershipEntity.MembershipTier  `json:"tier"`
	SubscriptionType membershipEntity.SubscriptionType `json:"subscription_type"`
	Status           membershipEntity.MembershipStatus `json:"status"`
	Price            float64                          `json:"price"`
	NextBillingDate  time.Time                        `json:"next_billing_date"`
	PaymentURL       string                           `json:"payment_url,omitempty"`
	Message          string                           `json:"message"`
}

// ToMembershipResponse converts Membership entity to MembershipResponse DTO
func ToMembershipResponse(m *membershipEntity.Membership) MembershipResponse {
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

