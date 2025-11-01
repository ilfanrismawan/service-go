package repository

import (
	"context"
	"service/internal/core"
	"service/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MembershipRepository handles membership data operations
type MembershipRepository struct {
	db *gorm.DB
}

// NewMembershipRepository creates a new membership repository
func NewMembershipRepository() *MembershipRepository {
	return &MembershipRepository{
		db: database.DB,
	}
}

// Create creates a new membership
func (r *MembershipRepository) Create(ctx context.Context, membership *core.Membership) error {
	return r.db.WithContext(ctx).Create(membership).Error
}

// GetByID gets a membership by ID
func (r *MembershipRepository) GetByID(ctx context.Context, id uuid.UUID) (*core.Membership, error) {
	var membership core.Membership
	err := r.db.WithContext(ctx).Preload("User").First(&membership, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

// GetByUserID gets a membership by user ID
func (r *MembershipRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*core.Membership, error) {
	var membership core.Membership
	err := r.db.WithContext(ctx).Preload("User").First(&membership, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

// Update updates a membership
func (r *MembershipRepository) Update(ctx context.Context, membership *core.Membership) error {
	return r.db.WithContext(ctx).Save(membership).Error
}

// Delete soft deletes a membership
func (r *MembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&core.Membership{}, "id = ?", id).Error
}

// List gets a list of memberships with pagination
func (r *MembershipRepository) List(ctx context.Context, page, limit int, tier *core.MembershipTier, status *core.MembershipStatus) ([]*core.Membership, int64, error) {
	var memberships []*core.Membership
	var total int64

	query := r.db.WithContext(ctx).Model(&core.Membership{}).Preload("User")

	// Apply filters
	if tier != nil {
		query = query.Where("tier = ?", *tier)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&memberships).Error; err != nil {
		return nil, 0, err
	}

	return memberships, total, nil
}

// GetActiveMemberships gets all active memberships
func (r *MembershipRepository) GetActiveMemberships(ctx context.Context) ([]*core.Membership, error) {
	var memberships []*core.Membership
	err := r.db.WithContext(ctx).Preload("User").Where("status = ?", core.MembershipStatusActive).Find(&memberships).Error
	return memberships, err
}

// GetMembershipsByTier gets memberships by tier
func (r *MembershipRepository) GetMembershipsByTier(ctx context.Context, tier core.MembershipTier) ([]*core.Membership, error) {
	var memberships []*core.Membership
	err := r.db.WithContext(ctx).Preload("User").Where("tier = ? AND status = ?", tier, core.MembershipStatusActive).Find(&memberships).Error
	return memberships, err
}

// UpdatePoints updates membership points
func (r *MembershipRepository) UpdatePoints(ctx context.Context, userID uuid.UUID, points int64) error {
	return r.db.WithContext(ctx).Model(&core.Membership{}).
		Where("user_id = ?", userID).
		Update("points", gorm.Expr("points + ?", points)).Error
}

// UpdateSpending updates membership spending and order count
func (r *MembershipRepository) UpdateSpending(ctx context.Context, userID uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Model(&core.Membership{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"total_spent":   gorm.Expr("total_spent + ?", amount),
			"orders_count":  gorm.Expr("orders_count + 1"),
			"last_order_at": gorm.Expr("NOW()"),
		}).Error
}

// GetTopSpenders gets top spending members
func (r *MembershipRepository) GetTopSpenders(ctx context.Context, limit int) ([]*core.Membership, error) {
	var memberships []*core.Membership
	err := r.db.WithContext(ctx).Preload("User").
		Where("status = ?", core.MembershipStatusActive).
		Order("total_spent DESC").
		Limit(limit).
		Find(&memberships).Error
	return memberships, err
}

// GetMembershipStats gets membership statistics
func (r *MembershipRepository) GetMembershipStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total memberships
	var totalMemberships int64
	if err := r.db.WithContext(ctx).Model(&core.Membership{}).Count(&totalMemberships).Error; err != nil {
		return nil, err
	}
	stats["total_memberships"] = totalMemberships

	// Active memberships
	var activeMemberships int64
	if err := r.db.WithContext(ctx).Model(&core.Membership{}).Where("status = ?", core.MembershipStatusActive).Count(&activeMemberships).Error; err != nil {
		return nil, err
	}
	stats["active_memberships"] = activeMemberships

	// Memberships by tier
	tierStats := make(map[string]int64)
	tiers := []core.MembershipTier{core.MembershipTierBasic, core.MembershipTierPremium, core.MembershipTierVIP, core.MembershipTierElite}

	for _, tier := range tiers {
		var count int64
		if err := r.db.WithContext(ctx).Model(&core.Membership{}).Where("tier = ? AND status = ?", tier, core.MembershipStatusActive).Count(&count).Error; err != nil {
			return nil, err
		}
		tierStats[string(tier)] = count
	}
	stats["memberships_by_tier"] = tierStats

	// Total points distributed
	var totalPoints int64
	if err := r.db.WithContext(ctx).Model(&core.Membership{}).Select("COALESCE(SUM(points), 0)").Scan(&totalPoints).Error; err != nil {
		return nil, err
	}
	stats["total_points"] = totalPoints

	// Total spending by members
	var totalSpending float64
	if err := r.db.WithContext(ctx).Model(&core.Membership{}).Select("COALESCE(SUM(total_spent), 0)").Scan(&totalSpending).Error; err != nil {
		return nil, err
	}
	stats["total_spending"] = totalSpending

	return stats, nil
}
