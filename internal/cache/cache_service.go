package cache

import (
	"context"
	"encoding/json"
	"fmt"
	branchEntity "service-go/internal/modules/branches/entity"
	membershipEntity "service-go/internal/modules/membership/entity"
	orderDTO "service-go/internal/modules/orders/dto"
	orderEntity "service-go/internal/modules/orders/entity"
	"service-go/internal/shared/database"
	"time"

	"github.com/google/uuid"
)

const (
	// Cache TTL
	CacheTTL_Branch       = 30 * time.Minute
	CacheTTL_Membership   = 15 * time.Minute
	CacheTTL_ServicePrice = 1 * time.Hour

	// Cache keys prefix
	CacheKey_Branch       = "branch:%s"
	CacheKey_BranchList   = "branches:list"
	CacheKey_Membership   = "membership:%s"
	CacheKey_ServicePrice = "service_price:%s"
)

// CacheService handles caching operations
type CacheService struct{}

// NewCacheService creates a new cache service
func NewCacheService() *CacheService {
	return &CacheService{}
}

// GetBranch retrieves branch from cache or returns nil
func (c *CacheService) GetBranch(ctx context.Context, branchID uuid.UUID) (*branchEntity.Branch, error) {
	if database.Redis == nil {
		return nil, fmt.Errorf("redis not initialized")
	}

	key := fmt.Sprintf(CacheKey_Branch, branchID.String())
	data, err := database.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var branch branchEntity.Branch
	if err := json.Unmarshal([]byte(data), &branch); err != nil {
		return nil, err
	}

	return &branch, nil
}

// SetBranch caches a branch
func (c *CacheService) SetBranch(ctx context.Context, branch *branchEntity.Branch) error {
	if database.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	key := fmt.Sprintf(CacheKey_Branch, branch.ID.String())
	data, err := json.Marshal(branch)
	if err != nil {
		return err
	}

	return database.Redis.Set(ctx, key, data, CacheTTL_Branch).Err()
}

// InvalidateBranch invalidates branch cache
func (c *CacheService) InvalidateBranch(ctx context.Context, branchID uuid.UUID) error {
	if database.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	key := fmt.Sprintf(CacheKey_Branch, branchID.String())
	return database.Redis.Del(ctx, key).Err()
}

// InvalidateBranchList invalidates branch list cache
func (c *CacheService) InvalidateBranchList(ctx context.Context) error {
	if database.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	return database.Redis.Del(ctx, CacheKey_BranchList).Err()
}

// GetMembership retrieves membership from cache
func (c *CacheService) GetMembership(ctx context.Context, userID uuid.UUID) (*membershipEntity.Membership, error) {
	if database.Redis == nil {
		return nil, fmt.Errorf("redis not initialized")
	}

	key := fmt.Sprintf(CacheKey_Membership, userID.String())
	data, err := database.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var membership membershipEntity.Membership
	if err := json.Unmarshal([]byte(data), &membership); err != nil {
		return nil, err
	}

	return &membership, nil
}

// SetMembership caches a membership
func (c *CacheService) SetMembership(ctx context.Context, membership *membershipEntity.Membership) error {
	if database.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	key := fmt.Sprintf(CacheKey_Membership, membership.UserID.String())
	data, err := json.Marshal(membership)
	if err != nil {
		return err
	}

	return database.Redis.Set(ctx, key, data, CacheTTL_Membership).Err()
}

// InvalidateMembership invalidates membership cache
func (c *CacheService) InvalidateMembership(ctx context.Context, userID uuid.UUID) error {
	if database.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	key := fmt.Sprintf(CacheKey_Membership, userID.String())
	return database.Redis.Del(ctx, key).Err()
}

// GetServicePrice retrieves service price estimate from cache
func (c *CacheService) GetServicePrice(ctx context.Context, serviceType orderEntity.ServiceType) (*orderDTO.ServiceEstimate, error) {
	if database.Redis == nil {
		return nil, fmt.Errorf("redis not initialized")
	}

	key := fmt.Sprintf(CacheKey_ServicePrice, string(serviceType))
	data, err := database.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var estimate orderDTO.ServiceEstimate
	if err := json.Unmarshal([]byte(data), &estimate); err != nil {
		return nil, err
	}

	return &estimate, nil
}

// SetServicePrice caches service price estimate
func (c *CacheService) SetServicePrice(ctx context.Context, serviceType orderEntity.ServiceType, estimate *orderDTO.ServiceEstimate) error {
	if database.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	key := fmt.Sprintf(CacheKey_ServicePrice, string(serviceType))
	data, err := json.Marshal(estimate)
	if err != nil {
		return err
	}

	return database.Redis.Set(ctx, key, data, CacheTTL_ServicePrice).Err()
}
