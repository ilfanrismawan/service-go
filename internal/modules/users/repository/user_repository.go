package repository

import (
	"context"
	"service/internal/shared/database"
	userEntity "service-go/internal/modules/users/entity"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *gorm.DB
	// in-memory fallback used when db == nil
	inMemory bool
}

var (
	sharedUsers   = make(map[uuid.UUID]*userEntity.User)
	sharedUsersMu sync.RWMutex
)

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	if database.DB == nil {
		return &UserRepository{
			db:       nil,
			inMemory: true,
		}
	}

	return &UserRepository{
		db:       database.DB,
		inMemory: false,
	}
}

// Create creates a new user
func (r *UserRepository) Register(ctx context.Context, user *userEntity.User) error {
	if r.inMemory {
		sharedUsersMu.Lock()
		defer sharedUsersMu.Unlock()
		if user.ID == uuid.Nil {
			user.ID = uuid.New()
		}
		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now
		sharedUsers[user.ID] = user
		return nil
	}
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*userEntity.User, error) {
	if r.inMemory {
		sharedUsersMu.RLock()
		defer sharedUsersMu.RUnlock()
		u, ok := sharedUsers[id]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		return u, nil
	}
	var user userEntity.User
	err := r.db.WithContext(ctx).Preload("Branch").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	if r.inMemory {
		sharedUsersMu.RLock()
		defer sharedUsersMu.RUnlock()
		for _, u := range sharedUsers {
			if u.Email == email {
				return u, nil
			}
		}
		return nil, gorm.ErrRecordNotFound
	}
	var user userEntity.User
	err := r.db.WithContext(ctx).Preload("Branch").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *userEntity.User) error {
	if r.inMemory {
		sharedUsersMu.Lock()
		defer sharedUsersMu.Unlock()
		if _, ok := sharedUsers[user.ID]; !ok {
			return gorm.ErrRecordNotFound
		}
		user.UpdatedAt = time.Now()
		sharedUsers[user.ID] = user
		return nil
	}
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if r.inMemory {
		sharedUsersMu.Lock()
		defer sharedUsersMu.Unlock()
		if _, ok := sharedUsers[id]; !ok {
			return gorm.ErrRecordNotFound
		}
		delete(sharedUsers, id)
		return nil
	}
	return r.db.WithContext(ctx).Delete(&userEntity.User{}, "id = ?", id).Error
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, offset, limit int, role *userEntity.UserRole, branchID *uuid.UUID) ([]*userEntity.User, int64, error) {
	if r.inMemory {
		sharedUsersMu.RLock()
		defer sharedUsersMu.RUnlock()
		var list []*userEntity.User
		for _, u := range sharedUsers {
			if role != nil && u.Role != *role {
				continue
			}
			if branchID != nil {
				if u.BranchID == nil || *u.BranchID != *branchID {
					continue
				}
			}
			list = append(list, u)
		}
		// sort by CreatedAt desc
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		total := int64(len(list))
		if offset > len(list) {
			return []*userEntity.User{}, total, nil
		}
		end := offset + limit
		if end > len(list) {
			end = len(list)
		}
		return list[offset:end], total, nil
	}

	var users []*userEntity.User
	var total int64

	query := r.db.WithContext(ctx).Model(&userEntity.User{})

	if role != nil {
		query = query.Where("role = ?", *role)
	}

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get users with pagination
	err := query.Preload("Branch").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// GetByBranchID retrieves users by branch ID
func (r *UserRepository) GetByBranchID(ctx context.Context, branchID uuid.UUID) ([]*userEntity.User, error) {
	if r.inMemory {
		sharedUsersMu.RLock()
		defer sharedUsersMu.RUnlock()
		var users []*userEntity.User
		for _, u := range sharedUsers {
			if u.BranchID != nil && *u.BranchID == branchID {
				users = append(users, u)
			}
		}
		return users, nil
	}
	var users []*userEntity.User
	err := r.db.WithContext(ctx).Preload("Branch").Where("branch_id = ?", branchID).Find(&users).Error
	return users, err
}

// CheckEmailExists checks if email already exists
func (r *UserRepository) CheckEmailExists(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	if r.inMemory {
		sharedUsersMu.RLock()
		defer sharedUsersMu.RUnlock()
		for id, u := range sharedUsers {
			if u.Email == email {
				if excludeID != nil && id == *excludeID {
					continue
				}
				return true, nil
			}
		}
		return false, nil
	}
	var count int64
	query := r.db.WithContext(ctx).Model(&userEntity.User{}).Where("email = ?", email)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// CheckPhoneExists checks if phone number already exists
func (r *UserRepository) CheckPhoneExists(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error) {
	if r.inMemory {
		sharedUsersMu.RLock()
		defer sharedUsersMu.RUnlock()
		for id, u := range sharedUsers {
			if u.Phone == phone {
				if excludeID != nil && id == *excludeID {
					continue
				}
				return true, nil
			}
		}
		return false, nil
	}
	var count int64
	query := r.db.WithContext(ctx).Model(&userEntity.User{}).Where("phone = ?", phone)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// CountCustomersByDateRange counts customers within a date range
func (r *UserRepository) CountCustomersByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	if r.inMemory {
		sharedUsersMu.RLock()
		defer sharedUsersMu.RUnlock()
		var count int64
		for _, u := range sharedUsers {
			if u.Role == userEntity.RolePelanggan && !u.CreatedAt.Before(startDate) && !u.CreatedAt.After(endDate) {
				count++
			}
		}
		return count, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&userEntity.User{}).
		Where("role = ? AND created_at >= ? AND created_at <= ?", userEntity.RolePelanggan, startDate, endDate).
		Count(&count).Error
	return count, err
}

// CountNewCustomersByDateRange counts new customers registered within a date range
func (r *UserRepository) CountNewCustomersByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	if r.inMemory {
		sharedUsersMu.RLock()
		defer sharedUsersMu.RUnlock()
		var count int64
		for _, u := range sharedUsers {
			if u.Role == userEntity.RolePelanggan && !u.CreatedAt.Before(startDate) && !u.CreatedAt.After(endDate) {
				count++
			}
		}
		return count, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&userEntity.User{}).
		Where("role = ? AND created_at >= ? AND created_at <= ?", userEntity.RolePelanggan, startDate, endDate).
		Count(&count).Error
	return count, err
}
