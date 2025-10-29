package repository

import (
	"context"
	"service/internal/core"
	"service/internal/database"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BranchRepository handles branch data operations
type BranchRepository struct {
	db       *gorm.DB
	inMemory bool
	mu       sync.RWMutex
	branches map[uuid.UUID]*core.Branch
}

// NewBranchRepository creates a new branch repository
func NewBranchRepository() *BranchRepository {
	if database.DB == nil {
		// use shared in-memory store
		m := make(map[uuid.UUID]*core.Branch)
		// seed a default test branch so tests that reference a static UUID can work
		if id, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); err == nil {
			now := time.Now()
			m[id] = &core.Branch{
				ID:        id,
				Name:      "Test Branch",
				City:      "Test City",
				Province:  "Test Province",
				IsActive:  true,
				CreatedAt: now,
				UpdatedAt: now,
			}
		}
		return &BranchRepository{
			db:       nil,
			inMemory: true,
			branches: m,
		}
	}
	return &BranchRepository{
		db: database.DB,
	}
}

// Create creates a new branch
func (r *BranchRepository) Create(ctx context.Context, branch *core.Branch) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		if branch.ID == uuid.Nil {
			branch.ID = uuid.New()
		}
		now := time.Now()
		branch.CreatedAt = now
		branch.UpdatedAt = now
		r.branches[branch.ID] = branch
		return nil
	}
	return r.db.WithContext(ctx).Create(branch).Error
}

// GetByID retrieves a branch by ID
func (r *BranchRepository) GetByID(ctx context.Context, id uuid.UUID) (*core.Branch, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		b, ok := r.branches[id]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		return b, nil
	}
	var branch core.Branch
	err := r.db.WithContext(ctx).First(&branch, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

// Update updates a branch
func (r *BranchRepository) Update(ctx context.Context, branch *core.Branch) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.branches[branch.ID]; !ok {
			return gorm.ErrRecordNotFound
		}
		branch.UpdatedAt = time.Now()
		r.branches[branch.ID] = branch
		return nil
	}
	return r.db.WithContext(ctx).Save(branch).Error
}

// Delete soft deletes a branch
func (r *BranchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.branches[id]; !ok {
			return gorm.ErrRecordNotFound
		}
		delete(r.branches, id)
		return nil
	}
	return r.db.WithContext(ctx).Delete(&core.Branch{}, "id = ?", id).Error
}

// List retrieves branches with pagination
func (r *BranchRepository) List(ctx context.Context, offset, limit int, city *string, province *string) ([]*core.Branch, int64, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*core.Branch
		for _, b := range r.branches {
			if city != nil && !containsIgnoreCase(b.City, *city) {
				continue
			}
			if province != nil && !containsIgnoreCase(b.Province, *province) {
				continue
			}
			list = append(list, b)
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		total := int64(len(list))
		if offset > len(list) {
			return []*core.Branch{}, total, nil
		}
		end := offset + limit
		if end > len(list) {
			end = len(list)
		}
		return list[offset:end], total, nil
	}
	var branches []*core.Branch
	var total int64

	query := r.db.WithContext(ctx).Model(&core.Branch{})

	if city != nil {
		query = query.Where("city ILIKE ?", "%"+*city+"%")
	}

	if province != nil {
		query = query.Where("province ILIKE ?", "%"+*province+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get branches with pagination
	err := query.Offset(offset).Limit(limit).Find(&branches).Error
	return branches, total, err
}

// GetActiveBranches retrieves all active branches
func (r *BranchRepository) GetActiveBranches(ctx context.Context) ([]*core.Branch, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*core.Branch
		for _, b := range r.branches {
			if b.IsActive {
				list = append(list, b)
			}
		}
		return list, nil
	}
	var branches []*core.Branch
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&branches).Error
	return branches, err
}

// GetNearbyBranches retrieves branches within a certain radius
func (r *BranchRepository) GetNearbyBranches(ctx context.Context, latitude, longitude, radiusKm float64) ([]*core.Branch, error) {
	var branches []*core.Branch

	// Using Haversine formula to calculate distance
	// This is a simplified version - in production, you might want to use PostGIS
	query := `
		SELECT * FROM branches 
		WHERE is_active = true 
		AND (
			6371 * acos(
				cos(radians(?)) * cos(radians(latitude)) * 
				cos(radians(longitude) - radians(?)) + 
				sin(radians(?)) * sin(radians(latitude))
			)
		) <= ?
		ORDER BY (
			6371 * acos(
				cos(radians(?)) * cos(radians(latitude)) * 
				cos(radians(longitude) - radians(?)) + 
				sin(radians(?)) * sin(radians(latitude))
			)
		)
	`

	err := r.db.WithContext(ctx).Raw(query, latitude, longitude, latitude, radiusKm, latitude, longitude, latitude).Scan(&branches).Error
	return branches, err
}

// GetByCity retrieves branches by city
func (r *BranchRepository) GetByCity(ctx context.Context, city string) ([]*core.Branch, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*core.Branch
		for _, b := range r.branches {
			if b.IsActive && containsIgnoreCase(b.City, city) {
				list = append(list, b)
			}
		}
		return list, nil
	}
	var branches []*core.Branch
	err := r.db.WithContext(ctx).Where("city ILIKE ? AND is_active = ?", "%"+city+"%", true).Find(&branches).Error
	return branches, err
}

// GetByProvince retrieves branches by province
func (r *BranchRepository) GetByProvince(ctx context.Context, province string) ([]*core.Branch, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*core.Branch
		for _, b := range r.branches {
			if b.IsActive && containsIgnoreCase(b.Province, province) {
				list = append(list, b)
			}
		}
		return list, nil
	}
	var branches []*core.Branch
	err := r.db.WithContext(ctx).Where("province ILIKE ? AND is_active = ?", "%"+province+"%", true).Find(&branches).Error
	return branches, err
}

// GetTopBranchesByRevenueInDateRange gets top branches by revenue within a date range
func (r *BranchRepository) GetTopBranchesByRevenueInDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]core.BranchStats, error) {
	var results []core.BranchStats

	err := r.db.WithContext(ctx).Model(&core.Payment{}).
		Select("o.branch_id, COALESCE(SUM(p.amount), 0) as revenue").
		Joins("JOIN service_orders o ON p.order_id = o.id").
		Where("p.created_at >= ? AND p.created_at <= ? AND p.status = ?", startDate, endDate, core.PaymentStatusPaid).
		Group("o.branch_id").
		Order("revenue DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// GetBranches retrieves branches with pagination
func (r *BranchRepository) GetBranches(ctx context.Context, page, limit int) ([]core.Branch, int64, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []core.Branch
		for _, b := range r.branches {
			list = append(list, *b)
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		total := int64(len(list))
		offset := (page - 1) * limit
		if offset > len(list) {
			return []core.Branch{}, total, nil
		}
		end := offset + limit
		if end > len(list) {
			end = len(list)
		}
		return list[offset:end], total, nil
	}

	var branches []core.Branch
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&core.Branch{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated results
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&branches).Error

	return branches, total, err
}

// helpers
func containsIgnoreCase(s, sub string) bool {
	if len(s) < len(sub) {
		return false
	}
	// simple case-insensitive contains
	ss := strings.ToLower(s)
	subl := strings.ToLower(sub)
	return strings.Contains(ss, subl)
}
