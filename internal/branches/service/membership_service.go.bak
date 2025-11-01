package service

import (
	"context"
	"fmt"
	"service/internal/core"
	"service/internal/repository"
	"time"

	"github.com/google/uuid"
)

// MembershipService handles membership business logic
type MembershipService struct {
	membershipRepo *repository.MembershipRepository
	userRepo       *repository.UserRepository
}

// NewMembershipService creates a new membership service
func NewMembershipService() *MembershipService {
	return &MembershipService{
		membershipRepo: repository.NewMembershipRepository(),
		userRepo:       repository.NewUserRepository(),
	}
}

// CreateMembership creates a new membership for a user
func (s *MembershipService) CreateMembership(ctx context.Context, userID uuid.UUID, tier core.MembershipTier) (*core.Membership, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user already has membership
	existingMembership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err == nil && existingMembership != nil {
		return nil, fmt.Errorf("user already has membership")
	}

	// Get tier configuration
	config := core.GetMembershipTierConfig(tier)

	// Create membership with trial period
	now := time.Now()
	trialEndsAt := now.AddDate(0, 0, 7) // 7 days trial

	membership := &core.Membership{
		UserID:             userID,
		Tier:               tier,
		Status:             core.MembershipStatusTrial,
		SubscriptionType:   core.SubscriptionTypeMonthly,
		DiscountPercentage: config.DiscountPercentage,
		Points:             0,
		TotalSpent:         0,
		OrdersCount:        0,
		MonthlyPrice:       config.MonthlyPrice,
		YearlyPrice:        config.YearlyPrice,
		CurrentPrice:       config.MonthlyPrice,
		JoinedAt:           now,
		ExpiresAt:          nil,
		NextBillingDate:    nil,
		AutoRenew:          true,
		TrialEndsAt:        &trialEndsAt,
	}

	if err := s.membershipRepo.Create(ctx, membership); err != nil {
		return nil, fmt.Errorf("failed to create membership: %w", err)
	}

	return membership, nil
}

// SubscribeToMembership creates a new subscription for a user
func (s *MembershipService) SubscribeToMembership(ctx context.Context, userID uuid.UUID, req *core.SubscriptionRequest) (*core.SubscriptionResult, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user already has active membership
	existingMembership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err == nil && existingMembership != nil && existingMembership.Status == core.MembershipStatusActive {
		return nil, fmt.Errorf("user already has active membership")
	}

	// Get tier configuration
	config := core.GetMembershipTierConfig(req.Tier)

	// Calculate price
	var price float64
	if req.SubscriptionType == core.SubscriptionTypeYearly {
		price = config.YearlyPrice
	} else {
		price = config.MonthlyPrice
	}

	// Create or update membership
	now := time.Now()
	var membership *core.Membership

	if existingMembership != nil {
		// Update existing membership
		membership = existingMembership
		membership.Tier = req.Tier
		membership.SubscriptionType = req.SubscriptionType
		membership.Status = core.MembershipStatusActive
		membership.DiscountPercentage = config.DiscountPercentage
		membership.MonthlyPrice = config.MonthlyPrice
		membership.YearlyPrice = config.YearlyPrice
		membership.CurrentPrice = price
		membership.AutoRenew = req.AutoRenew
		membership.TrialEndsAt = nil
	} else {
		// Create new membership
		membership = &core.Membership{
			UserID:             userID,
			Tier:               req.Tier,
			Status:             core.MembershipStatusActive,
			SubscriptionType:   req.SubscriptionType,
			DiscountPercentage: config.DiscountPercentage,
			Points:             0,
			TotalSpent:         0,
			OrdersCount:        0,
			MonthlyPrice:       config.MonthlyPrice,
			YearlyPrice:        config.YearlyPrice,
			CurrentPrice:       price,
			JoinedAt:           now,
			AutoRenew:          req.AutoRenew,
		}
	}

	// Set billing dates
	if req.SubscriptionType == core.SubscriptionTypeYearly {
		nextBilling := now.AddDate(1, 0, 0)
		membership.NextBillingDate = &nextBilling
		expiresAt := now.AddDate(1, 0, 0)
		membership.ExpiresAt = &expiresAt
	} else {
		nextBilling := now.AddDate(0, 1, 0)
		membership.NextBillingDate = &nextBilling
		expiresAt := now.AddDate(0, 1, 0)
		membership.ExpiresAt = &expiresAt
	}

	// Save membership
	if existingMembership != nil {
		if err := s.membershipRepo.Update(ctx, membership); err != nil {
			return nil, fmt.Errorf("failed to update membership: %w", err)
		}
	} else {
		if err := s.membershipRepo.Create(ctx, membership); err != nil {
			return nil, fmt.Errorf("failed to create membership: %w", err)
		}
	}

	// TODO: Integrate with payment gateway to get payment URL
	paymentURL := fmt.Sprintf("/payment/subscription/%s", membership.ID.String())

	return &core.SubscriptionResult{
		MembershipID:     membership.ID,
		Tier:             membership.Tier,
		SubscriptionType: membership.SubscriptionType,
		Status:           membership.Status,
		Price:            price,
		NextBillingDate:  *membership.NextBillingDate,
		PaymentURL:       paymentURL,
		Message:          "Subscription created successfully",
	}, nil
}

// CancelSubscription cancels a user's subscription
func (s *MembershipService) CancelSubscription(ctx context.Context, userID uuid.UUID, req *core.CancelSubscriptionRequest) error {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("membership not found: %w", err)
	}

	if membership.Status != core.MembershipStatusActive {
		return fmt.Errorf("membership is not active")
	}

	// Update membership status
	membership.Status = core.MembershipStatusCancelled
	membership.AutoRenew = false

	if err := s.membershipRepo.Update(ctx, membership); err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return nil
}

// StartTrial starts a trial membership for a user
func (s *MembershipService) StartTrial(ctx context.Context, userID uuid.UUID, tier core.MembershipTier) (*core.Membership, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user already has membership
	existingMembership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err == nil && existingMembership != nil {
		return nil, fmt.Errorf("user already has membership")
	}

	// Get tier configuration
	config := core.GetMembershipTierConfig(tier)

	// Create trial membership
	now := time.Now()
	trialEndsAt := now.AddDate(0, 0, 7) // 7 days trial

	membership := &core.Membership{
		UserID:             userID,
		Tier:               tier,
		Status:             core.MembershipStatusTrial,
		SubscriptionType:   core.SubscriptionTypeMonthly,
		DiscountPercentage: config.DiscountPercentage,
		Points:             0,
		TotalSpent:         0,
		OrdersCount:        0,
		MonthlyPrice:       config.MonthlyPrice,
		YearlyPrice:        config.YearlyPrice,
		CurrentPrice:       0, // Free during trial
		JoinedAt:           now,
		ExpiresAt:          &trialEndsAt,
		NextBillingDate:    nil,
		AutoRenew:          false,
		TrialEndsAt:        &trialEndsAt,
	}

	if err := s.membershipRepo.Create(ctx, membership); err != nil {
		return nil, fmt.Errorf("failed to create trial membership: %w", err)
	}

	return membership, nil
}

// GetMembership gets a membership by user ID
func (s *MembershipService) GetMembership(ctx context.Context, userID uuid.UUID) (*core.Membership, error) {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("membership not found: %w", err)
	}
	return membership, nil
}

// UpdateMembership updates a membership
func (s *MembershipService) UpdateMembership(ctx context.Context, userID uuid.UUID, req *core.MembershipRequest) (*core.Membership, error) {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("membership not found: %w", err)
	}

	// Update tier if provided
	if req.Tier != "" {
		config := core.GetMembershipTierConfig(req.Tier)
		membership.Tier = req.Tier
		membership.DiscountPercentage = config.DiscountPercentage
	}

	// Update status if provided
	if req.Status != "" {
		membership.Status = req.Status
	}

	if err := s.membershipRepo.Update(ctx, membership); err != nil {
		return nil, fmt.Errorf("failed to update membership: %w", err)
	}

	return membership, nil
}

// ProcessOrderWithMembership processes an order with membership benefits
func (s *MembershipService) ProcessOrderWithMembership(ctx context.Context, userID uuid.UUID, orderAmount float64, isService bool, isPickup bool) (*core.MembershipOrderResult, error) {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		// If no membership, return without discount
		return &core.MembershipOrderResult{
			HasMembership:      false,
			OriginalAmount:     orderAmount,
			DiscountAmount:     0,
			FinalAmount:        orderAmount,
			PointsEarned:       0,
			DiscountPercentage: 0,
		}, nil
	}

	// Check if membership is active or in trial
	if membership.Status != core.MembershipStatusActive && membership.Status != core.MembershipStatusTrial {
		return &core.MembershipOrderResult{
			HasMembership:      false,
			OriginalAmount:     orderAmount,
			DiscountAmount:     0,
			FinalAmount:        orderAmount,
			PointsEarned:       0,
			DiscountPercentage: 0,
		}, nil
	}

	// Check if trial has expired
	if membership.Status == core.MembershipStatusTrial && membership.IsTrialExpired() {
		return &core.MembershipOrderResult{
			HasMembership:      false,
			OriginalAmount:     orderAmount,
			DiscountAmount:     0,
			FinalAmount:        orderAmount,
			PointsEarned:       0,
			DiscountPercentage: 0,
		}, nil
	}

	config := core.GetMembershipTierConfig(membership.Tier)

	// Check for free services and pickups
	var freeServiceUsed, freePickupUsed bool
	var remainingFreeServices, remainingFreePickups int

	if isService && config.MaxFreeServices > 0 {
		// TODO: Check current month usage for free services
		// For now, assume we can use free service
		freeServiceUsed = true
		remainingFreeServices = config.MaxFreeServices - 1
	}

	if isPickup && config.MaxFreePickups > 0 {
		// TODO: Check current month usage for free pickups
		// For now, assume we can use free pickup
		freePickupUsed = true
		remainingFreePickups = config.MaxFreePickups - 1
	}

	// Calculate discount
	discountAmount := membership.CalculateDiscount(orderAmount)
	finalAmount := orderAmount - discountAmount

	// If free service is used, make it free
	if freeServiceUsed {
		finalAmount = 0
		discountAmount = orderAmount
	}

	// Calculate points earned
	pointsEarned := membership.CalculatePoints(orderAmount)

	// Update membership spending
	if err := s.membershipRepo.UpdateSpending(ctx, userID, orderAmount); err != nil {
		return nil, fmt.Errorf("failed to update membership spending: %w", err)
	}

	// Update points
	if err := s.membershipRepo.UpdatePoints(ctx, userID, pointsEarned); err != nil {
		return nil, fmt.Errorf("failed to update membership points: %w", err)
	}

	return &core.MembershipOrderResult{
		HasMembership:         true,
		OriginalAmount:        orderAmount,
		DiscountAmount:        discountAmount,
		FinalAmount:           finalAmount,
		PointsEarned:          pointsEarned,
		DiscountPercentage:    membership.DiscountPercentage,
		Tier:                  membership.Tier,
		FreeServiceUsed:       freeServiceUsed,
		FreePickupUsed:        freePickupUsed,
		RemainingFreeServices: remainingFreeServices,
		RemainingFreePickups:  remainingFreePickups,
	}, nil
}

// ListMemberships gets a list of memberships with pagination
func (s *MembershipService) ListMemberships(ctx context.Context, page, limit int, tier *core.MembershipTier, status *core.MembershipStatus) ([]*core.Membership, int64, error) {
	return s.membershipRepo.List(ctx, page, limit, tier, status)
}

// GetMembershipStats gets membership statistics
func (s *MembershipService) GetMembershipStats(ctx context.Context) (map[string]interface{}, error) {
	return s.membershipRepo.GetMembershipStats(ctx)
}

// GetTopSpenders gets top spending members
func (s *MembershipService) GetTopSpenders(ctx context.Context, limit int) ([]*core.Membership, error) {
	return s.membershipRepo.GetTopSpenders(ctx, limit)
}

// RedeemPoints redeems membership points for discount
func (s *MembershipService) RedeemPoints(ctx context.Context, userID uuid.UUID, pointsToRedeem int64) (*core.MembershipPointsResult, error) {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("membership not found: %w", err)
	}

	if membership.Status != core.MembershipStatusActive {
		return nil, fmt.Errorf("membership is not active")
	}

	if membership.Points < pointsToRedeem {
		return nil, fmt.Errorf("insufficient points")
	}

	// Calculate discount value (1 point = 100 IDR)
	discountValue := float64(pointsToRedeem) * 100

	// Update points
	newPoints := membership.Points - pointsToRedeem
	membership.Points = newPoints

	if err := s.membershipRepo.Update(ctx, membership); err != nil {
		return nil, fmt.Errorf("failed to update points: %w", err)
	}

	return &core.MembershipPointsResult{
		PointsRedeemed:  pointsToRedeem,
		RemainingPoints: newPoints,
		DiscountValue:   discountValue,
	}, nil
}

// GetMembershipTiers returns all available membership tiers
func (s *MembershipService) GetMembershipTiers(ctx context.Context) []core.MembershipTierConfig {
	return core.GetMembershipTierConfigs()
}

// UpgradeMembership upgrades a user's membership tier
func (s *MembershipService) UpgradeMembership(ctx context.Context, userID uuid.UUID, newTier core.MembershipTier, subscriptionType core.SubscriptionType) (*core.SubscriptionResult, error) {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("membership not found: %w", err)
	}

	if membership.Status != core.MembershipStatusActive && membership.Status != core.MembershipStatusTrial {
		return nil, fmt.Errorf("membership is not active")
	}

	// Get new tier configuration
	config := core.GetMembershipTierConfig(newTier)

	// Calculate price difference
	var priceDifference float64
	if subscriptionType == core.SubscriptionTypeYearly {
		priceDifference = config.YearlyPrice - membership.YearlyPrice
	} else {
		priceDifference = config.MonthlyPrice - membership.MonthlyPrice
	}

	// Update membership
	membership.Tier = newTier
	membership.SubscriptionType = subscriptionType
	membership.DiscountPercentage = config.DiscountPercentage
	membership.MonthlyPrice = config.MonthlyPrice
	membership.YearlyPrice = config.YearlyPrice

	if subscriptionType == core.SubscriptionTypeYearly {
		membership.CurrentPrice = config.YearlyPrice
	} else {
		membership.CurrentPrice = config.MonthlyPrice
	}

	// Update billing dates
	now := time.Now()
	if subscriptionType == core.SubscriptionTypeYearly {
		nextBilling := now.AddDate(1, 0, 0)
		membership.NextBillingDate = &nextBilling
		expiresAt := now.AddDate(1, 0, 0)
		membership.ExpiresAt = &expiresAt
	} else {
		nextBilling := now.AddDate(0, 1, 0)
		membership.NextBillingDate = &nextBilling
		expiresAt := now.AddDate(0, 1, 0)
		membership.ExpiresAt = &expiresAt
	}

	if err := s.membershipRepo.Update(ctx, membership); err != nil {
		return nil, fmt.Errorf("failed to upgrade membership: %w", err)
	}

	// TODO: Integrate with payment gateway for upgrade payment
	paymentURL := fmt.Sprintf("/payment/upgrade/%s", membership.ID.String())

	return &core.SubscriptionResult{
		MembershipID:     membership.ID,
		Tier:             membership.Tier,
		SubscriptionType: membership.SubscriptionType,
		Status:           membership.Status,
		Price:            priceDifference,
		NextBillingDate:  *membership.NextBillingDate,
		PaymentURL:       paymentURL,
		Message:          "Membership upgraded successfully",
	}, nil
}

// GetMembershipUsage gets the current month usage for a membership
func (s *MembershipService) GetMembershipUsage(ctx context.Context, userID uuid.UUID) (*core.MembershipUsage, error) {
	// TODO: Implement usage tracking
	// For now, return empty usage
	now := time.Now()
	return &core.MembershipUsage{
		MembershipID:     userID,
		FreeServicesUsed: 0,
		FreePickupsUsed:  0,
		Month:            int(now.Month()),
		Year:             now.Year(),
	}, nil
}
