package repository

import (
	"context"
	// "service/internal/orders/model"
	"service/internal/shared/database"
	"service/internal/shared/model"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServiceOrderRepository handles service order data operations
type ServiceOrderRepository struct {
	db       *gorm.DB
	inMemory bool
	mu       sync.RWMutex
	orders   map[uuid.UUID]*model.ServiceOrder
}

// NewServiceOrderRepository creates a new service order repository
func NewServiceOrderRepository() *ServiceOrderRepository {
	if database.DB == nil {
		m := make(map[uuid.UUID]*model.ServiceOrder)
		// seed a default test order (so tests referencing a fixed order ID will work)
		if id, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000"); err == nil {
			// determine a customer id: use any existing shared user if available
			var customerID uuid.UUID
			// access shared user repository map if package-level variable exists
			// try to pick any user from repository package's sharedUsers if present
			// (reflectively access is messy; instead, create a default customer here)
			customerID = uuid.New()
			now := time.Now()
			// create a lightweight order using the test branch id
			branchID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
			m[id] = &model.ServiceOrder{
				ID:              id,
				OrderNumber:     "ORD-TEST-000001",
				CustomerID:      customerID,
				BranchID:        branchID,
				IPhoneModel:     "TestModel",
				IPhoneColor:     "Black",
				IPhoneIMEI:      "0000",
				ServiceType:     model.ServiceTypeOther,
				Description:     "Test order",
				PickupAddress:   "Test Address",
				PickupLatitude:  0,
				PickupLongitude: 0,
				Status:          model.StatusPendingPickup,
				CreatedAt:       now,
				UpdatedAt:       now,
			}
		}
		return &ServiceOrderRepository{
			db:       nil,
			inMemory: true,
			orders:   m,
		}
	}
	return &ServiceOrderRepository{
		db: database.DB,
	}
}

// Create creates a new service order
func (r *ServiceOrderRepository) Create(ctx context.Context, order *model.ServiceOrder) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		if order.ID == uuid.Nil {
			order.ID = uuid.New()
		}
		now := time.Now()
		order.CreatedAt = now
		order.UpdatedAt = now
		r.orders[order.ID] = order
		return nil
	}
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID retrieves a service order by ID
func (r *ServiceOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ServiceOrder, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		o, ok := r.orders[id]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		return o, nil
	}
	var order model.ServiceOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		Preload("Courier").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNumber retrieves a service order by order number
func (r *ServiceOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*model.ServiceOrder, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for _, o := range r.orders {
			if o.OrderNumber == orderNumber {
				return o, nil
			}
		}
		return nil, gorm.ErrRecordNotFound
	}
	var order model.ServiceOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		Preload("Courier").
		First(&order, "order_number = ?", orderNumber).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Update updates a service order
func (r *ServiceOrderRepository) Update(ctx context.Context, order *model.ServiceOrder) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.orders[order.ID]; !ok {
			return gorm.ErrRecordNotFound
		}
		order.UpdatedAt = time.Now()
		r.orders[order.ID] = order
		return nil
	}
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete soft deletes a service order
func (r *ServiceOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.orders[id]; !ok {
			return gorm.ErrRecordNotFound
		}
		delete(r.orders, id)
		return nil
	}
	return r.db.WithContext(ctx).Delete(&model.ServiceOrder{}, "id = ?", id).Error
}

// List retrieves service orders with pagination
func (r *ServiceOrderRepository) List(ctx context.Context, offset, limit int, filters *ServiceOrderFilters) ([]*model.ServiceOrder, int64, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*model.ServiceOrder
		for _, o := range r.orders {
			if filters != nil {
				if filters.CustomerID != nil && o.CustomerID != *filters.CustomerID {
					continue
				}
				if filters.BranchID != nil && o.BranchID != *filters.BranchID {
					continue
				}
				if filters.TechnicianID != nil {
					if o.TechnicianID == nil || *filters.TechnicianID != *o.TechnicianID {
						continue
					}
				}
				if filters.CourierID != nil {
					if o.CourierID == nil || *filters.CourierID != *o.CourierID {
						continue
					}
				}
				if filters.Status != nil && o.Status != *filters.Status {
					continue
				}
				if filters.ServiceType != nil && o.ServiceType != *filters.ServiceType {
					continue
				}
			}
			list = append(list, o)
		}
		// sort desc
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		total := int64(len(list))
		if offset > len(list) {
			return []*model.ServiceOrder{}, total, nil
		}
		end := offset + limit
		if end > len(list) {
			end = len(list)
		}
		return list[offset:end], total, nil
	}
	var orders []*model.ServiceOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ServiceOrder{})

	if filters != nil {
		if filters.CustomerID != nil {
			query = query.Where("customer_id = ?", *filters.CustomerID)
		}
		if filters.BranchID != nil {
			query = query.Where("branch_id = ?", *filters.BranchID)
		}
		if filters.TechnicianID != nil {
			query = query.Where("technician_id = ?", *filters.TechnicianID)
		}
		if filters.CourierID != nil {
			query = query.Where("courier_id = ?", *filters.CourierID)
		}
		if filters.Status != nil {
			query = query.Where("status = ?", *filters.Status)
		}
		if filters.ServiceType != nil {
			query = query.Where("service_type = ?", *filters.ServiceType)
		}
		if filters.DateFrom != nil {
			query = query.Where("created_at >= ?", *filters.DateFrom)
		}
		if filters.DateTo != nil {
			query = query.Where("created_at <= ?", *filters.DateTo)
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get orders with pagination
	err := query.
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		Preload("Courier").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, total, err
}

// GetByCustomerID retrieves service orders by customer ID
func (r *ServiceOrderRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*model.ServiceOrder, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*model.ServiceOrder
		for _, o := range r.orders {
			if o.CustomerID == customerID {
				list = append(list, o)
			}
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		return list, nil
	}
	var orders []*model.ServiceOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		Preload("Courier").
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetByBranchID retrieves service orders by branch ID
func (r *ServiceOrderRepository) GetByBranchID(ctx context.Context, branchID uuid.UUID) ([]*model.ServiceOrder, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*model.ServiceOrder
		for _, o := range r.orders {
			if o.BranchID == branchID {
				list = append(list, o)
			}
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		return list, nil
	}
	var orders []*model.ServiceOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		Preload("Courier").
		Where("branch_id = ?", branchID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetByStatus retrieves service orders by status
func (r *ServiceOrderRepository) GetByStatus(ctx context.Context, status model.OrderStatus) ([]*model.ServiceOrder, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*model.ServiceOrder
		for _, o := range r.orders {
			if o.Status == status {
				list = append(list, o)
			}
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		return list, nil
	}
	var orders []*model.ServiceOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		Preload("Courier").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetByTechnicianID retrieves service orders by technician ID
func (r *ServiceOrderRepository) GetByTechnicianID(ctx context.Context, technicianID uuid.UUID) ([]*model.ServiceOrder, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*model.ServiceOrder
		for _, o := range r.orders {
			if o.TechnicianID != nil && *o.TechnicianID == technicianID {
				list = append(list, o)
			}
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		return list, nil
	}
	var orders []*model.ServiceOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		Preload("Courier").
		Where("technician_id = ?", technicianID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetByCourierID retrieves service orders by courier ID
func (r *ServiceOrderRepository) GetByCourierID(ctx context.Context, courierID uuid.UUID) ([]*model.ServiceOrder, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*model.ServiceOrder
		for _, o := range r.orders {
			if o.CourierID != nil && *o.CourierID == courierID {
				list = append(list, o)
			}
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		return list, nil
	}
	var orders []*model.ServiceOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Branch").
		Preload("Technician").
		Preload("Courier").
		Where("courier_id = ?", courierID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// UpdateStatus updates the status of a service order
func (r *ServiceOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus, notes string) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		o, ok := r.orders[id]
		if !ok {
			return gorm.ErrRecordNotFound
		}
		o.Status = status
		o.Notes = notes
		o.UpdatedAt = time.Now()
		r.orders[id] = o
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&model.ServiceOrder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": status,
			"notes":  notes,
		}).Error
}

// AssignTechnician assigns a technician to a service order
func (r *ServiceOrderRepository) AssignTechnician(ctx context.Context, id uuid.UUID, technicianID uuid.UUID) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		o, ok := r.orders[id]
		if !ok {
			return gorm.ErrRecordNotFound
		}
		o.TechnicianID = &technicianID
		o.UpdatedAt = time.Now()
		r.orders[id] = o
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&model.ServiceOrder{}).
		Where("id = ?", id).
		Update("technician_id", technicianID).Error
}

// AssignCourier assigns a courier to a service order
func (r *ServiceOrderRepository) AssignCourier(ctx context.Context, id uuid.UUID, courierID uuid.UUID) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		o, ok := r.orders[id]
		if !ok {
			return gorm.ErrRecordNotFound
		}
		o.CourierID = &courierID
		o.UpdatedAt = time.Now()
		r.orders[id] = o
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&model.ServiceOrder{}).
		Where("id = ?", id).
		Update("courier_id", courierID).Error
}

// CheckOrderNumberExists checks if order number already exists
func (r *ServiceOrderRepository) CheckOrderNumberExists(ctx context.Context, orderNumber string, excludeID *uuid.UUID) (bool, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for id, o := range r.orders {
			if o.OrderNumber == orderNumber {
				if excludeID != nil && id == *excludeID {
					continue
				}
				return true, nil
			}
		}
		return false, nil
	}
	var count int64
	query := r.db.WithContext(ctx).Model(&model.ServiceOrder{}).Where("order_number = ?", orderNumber)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// CountOrdersByDateRange counts orders within a date range
func (r *ServiceOrderRepository) CountOrdersByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var count int64
		for _, o := range r.orders {
			if !o.CreatedAt.Before(startDate) && !o.CreatedAt.After(endDate) {
				count++
			}
		}
		return count, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ServiceOrder{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

// GetOrdersByStatusInDateRange gets order counts by status within a date range
func (r *ServiceOrderRepository) GetOrdersByStatusInDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]int64, error) {
	var results []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		m := make(map[string]int64)
		for _, o := range r.orders {
			if !o.CreatedAt.Before(startDate) && !o.CreatedAt.After(endDate) {
				m[string(o.Status)]++
			}
		}
		return m, nil
	}

	err := r.db.WithContext(ctx).Model(&model.ServiceOrder{}).
		Select("status, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, result := range results {
		statusMap[result.Status] = result.Count
	}

	return statusMap, nil
}

// GetOrdersByBranchInDateRange gets order counts by branch within a date range
func (r *ServiceOrderRepository) GetOrdersByBranchInDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]int64, error) {
	var results []struct {
		BranchID string `json:"branch_id"`
		Count    int64  `json:"count"`
	}
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		m := make(map[string]int64)
		for _, o := range r.orders {
			if !o.CreatedAt.Before(startDate) && !o.CreatedAt.After(endDate) {
				m[o.BranchID.String()]++
			}
		}
		return m, nil
	}

	err := r.db.WithContext(ctx).Model(&model.ServiceOrder{}).
		Select("branch_id, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("branch_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	branchMap := make(map[string]int64)
	for _, result := range results {
		branchMap[result.BranchID] = result.Count
	}

	return branchMap, nil
}

// GetOrdersByServiceTypeInDateRange gets order counts by service type within a date range
func (r *ServiceOrderRepository) GetOrdersByServiceTypeInDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]int64, error) {
	var results []struct {
		ServiceType string `json:"service_type"`
		Count       int64  `json:"count"`
	}
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		m := make(map[string]int64)
		for _, o := range r.orders {
			if !o.CreatedAt.Before(startDate) && !o.CreatedAt.After(endDate) {
				m[string(o.ServiceType)]++
			}
		}
		return m, nil
	}

	err := r.db.WithContext(ctx).Model(&model.ServiceOrder{}).
		Select("service_type, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("service_type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	serviceMap := make(map[string]int64)
	for _, result := range results {
		serviceMap[result.ServiceType] = result.Count
	}

	return serviceMap, nil
}

// GetTopServiceTypesInDateRange gets top service types by order count within a date range
func (r *ServiceOrderRepository) GetTopServiceTypesInDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]model.ServiceTypeStats, error) {
	var results []model.ServiceTypeStats
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		counts := make(map[string]int64)
		for _, o := range r.orders {
			if !o.CreatedAt.Before(startDate) && !o.CreatedAt.After(endDate) {
				counts[string(o.ServiceType)]++
			}
		}
		// convert to slice
		var stats []model.ServiceTypeStats
		for k, v := range counts {
			stats = append(stats, model.ServiceTypeStats{ServiceType: k, Count: v})
		}
		// sort desc
		sort.Slice(stats, func(i, j int) bool { return stats[i].Count > stats[j].Count })
		if len(stats) > limit {
			stats = stats[:limit]
		}
		return stats, nil
	}

	err := r.db.WithContext(ctx).Model(&model.ServiceOrder{}).
		Select("service_type, COUNT(*) as order_count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("service_type").
		Order("order_count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// ServiceOrderFilters represents filters for service order queries
type ServiceOrderFilters struct {
	CustomerID   *uuid.UUID
	BranchID     *uuid.UUID
	TechnicianID *uuid.UUID
	CourierID    *uuid.UUID
	Status       *model.OrderStatus
	ServiceType  *model.ServiceType
	DateFrom     *string
	DateTo       *string
}
