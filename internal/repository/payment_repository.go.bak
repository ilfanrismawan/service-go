package repository

import (
	"context"
	"service/internal/core"
	"service/internal/database"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentRepository handles payment data operations
type PaymentRepository struct {
	db       *gorm.DB
	inMemory bool
	mu       sync.RWMutex
	payments map[uuid.UUID]*core.Payment
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository() *PaymentRepository {
	if database.DB == nil {
		m := make(map[uuid.UUID]*core.Payment)
		return &PaymentRepository{
			db:       nil,
			inMemory: true,
			payments: m,
		}
	}
	return &PaymentRepository{
		db: database.DB,
	}
}

// Create creates a new payment
func (r *PaymentRepository) Create(ctx context.Context, payment *core.Payment) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		if payment.ID == uuid.Nil {
			payment.ID = uuid.New()
		}
		now := time.Now()
		payment.CreatedAt = now
		payment.UpdatedAt = now
		r.payments[payment.ID] = payment
		return nil
	}
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*core.Payment, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		p, ok := r.payments[id]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		return p, nil
	}
	var payment core.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Order.Customer").
		Preload("Order.Branch").
		First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByInvoiceNumber retrieves a payment by invoice number
func (r *PaymentRepository) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*core.Payment, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for _, p := range r.payments {
			if p.InvoiceNumber == invoiceNumber {
				return p, nil
			}
		}
		return nil, gorm.ErrRecordNotFound
	}
	var payment core.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Order.Customer").
		Preload("Order.Branch").
		First(&payment, "invoice_number = ?", invoiceNumber).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByTransactionID retrieves a payment by its transaction ID
func (r *PaymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*core.Payment, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for _, p := range r.payments {
			if p.TransactionID == transactionID {
				return p, nil
			}
		}
		return nil, gorm.ErrRecordNotFound
	}
	var payment core.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Order.Customer").
		Preload("Order.Branch").
		First(&payment, "transaction_id = ?", transactionID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// Update updates a payment
func (r *PaymentRepository) Update(ctx context.Context, payment *core.Payment) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.payments[payment.ID]; !ok {
			return gorm.ErrRecordNotFound
		}
		payment.UpdatedAt = time.Now()
		r.payments[payment.ID] = payment
		return nil
	}
	return r.db.WithContext(ctx).Save(payment).Error
}

// Delete soft deletes a payment
func (r *PaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if r.inMemory {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.payments[id]; !ok {
			return gorm.ErrRecordNotFound
		}
		delete(r.payments, id)
		return nil
	}
	return r.db.WithContext(ctx).Delete(&core.Payment{}, "id = ?", id).Error
}

// List retrieves payments with pagination
func (r *PaymentRepository) List(ctx context.Context, offset, limit int, filters *PaymentFilters) ([]*core.Payment, int64, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*core.Payment
		for _, p := range r.payments {
			if filters != nil {
				if filters.OrderID != nil && p.OrderID != *filters.OrderID {
					continue
				}
				if filters.Status != nil && p.Status != *filters.Status {
					continue
				}
				if filters.PaymentMethod != nil && p.PaymentMethod != *filters.PaymentMethod {
					continue
				}
			}
			list = append(list, p)
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		total := int64(len(list))
		if offset > len(list) {
			return []*core.Payment{}, total, nil
		}
		end := offset + limit
		if end > len(list) {
			end = len(list)
		}
		return list[offset:end], total, nil
	}
	var payments []*core.Payment
	var total int64

	query := r.db.WithContext(ctx).Model(&core.Payment{})

	if filters != nil {
		if filters.OrderID != nil {
			query = query.Where("order_id = ?", *filters.OrderID)
		}
		if filters.Status != nil {
			query = query.Where("status = ?", *filters.Status)
		}
		if filters.PaymentMethod != nil {
			query = query.Where("payment_method = ?", *filters.PaymentMethod)
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

	// Get payments with pagination
	err := query.
		Preload("Order").
		Preload("Order.Customer").
		Preload("Order.Branch").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, total, err
}

// GetByOrderID retrieves payments by order ID
func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*core.Payment, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*core.Payment
		for _, p := range r.payments {
			if p.OrderID == orderID {
				list = append(list, p)
			}
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		return list, nil
	}
	var payments []*core.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Order.Customer").
		Preload("Order.Branch").
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// GetByStatus retrieves payments by status
func (r *PaymentRepository) GetByStatus(ctx context.Context, status core.PaymentStatus) ([]*core.Payment, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*core.Payment
		for _, p := range r.payments {
			if p.Status == status {
				list = append(list, p)
			}
		}
		sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
		return list, nil
	}
	var payments []*core.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Order.Customer").
		Preload("Order.Branch").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// CheckInvoiceExists checks if invoice number already exists
func (r *PaymentRepository) CheckInvoiceExists(ctx context.Context, invoiceNumber string, excludeID *uuid.UUID) (bool, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for id, p := range r.payments {
			if p.InvoiceNumber == invoiceNumber {
				if excludeID != nil && id == *excludeID {
					continue
				}
				return true, nil
			}
		}
		return false, nil
	}
	var count int64
	query := r.db.WithContext(ctx).Model(&core.Payment{}).Where("invoice_number = ?", invoiceNumber)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetTotalRevenueByDateRange gets total revenue within a date range
func (r *PaymentRepository) GetTotalRevenueByDateRange(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var total float64
		for _, p := range r.payments {
			if !p.CreatedAt.Before(startDate) && !p.CreatedAt.After(endDate) && p.Status == core.PaymentStatusPaid {
				total += p.Amount
			}
		}
		return total, nil
	}
	var total float64
	err := r.db.WithContext(ctx).Model(&core.Payment{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("created_at >= ? AND created_at <= ? AND status = ?", startDate, endDate, core.PaymentStatusPaid).
		Scan(&total).Error
	return total, err
}

// GetRevenueByBranchInDateRange gets revenue by branch within a date range
func (r *PaymentRepository) GetRevenueByBranchInDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]float64, error) {
	var results []struct {
		BranchID string  `json:"branch_id"`
		Revenue  float64 `json:"revenue"`
	}
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		m := make(map[string]float64)
		for _, p := range r.payments {
			if !p.CreatedAt.Before(startDate) && !p.CreatedAt.After(endDate) && p.Status == core.PaymentStatusPaid {
				m[p.OrderID.String()] += p.Amount
			}
		}
		return m, nil
	}

	err := r.db.WithContext(ctx).Model(&core.Payment{}).
		Select("o.branch_id, COALESCE(SUM(p.amount), 0) as revenue").
		Joins("JOIN service_orders o ON p.order_id = o.id").
		Where("p.created_at >= ? AND p.created_at <= ? AND p.status = ?", startDate, endDate, core.PaymentStatusPaid).
		Group("o.branch_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	branchMap := make(map[string]float64)
	for _, result := range results {
		branchMap[result.BranchID] = result.Revenue
	}

	return branchMap, nil
}

// GetRevenueByPaymentMethodInDateRange gets revenue by payment method within a date range
func (r *PaymentRepository) GetRevenueByPaymentMethodInDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]float64, error) {
	var results []struct {
		PaymentMethod string  `json:"payment_method"`
		Revenue       float64 `json:"revenue"`
	}
	if r.inMemory {
		r.mu.RLock()
		defer r.mu.RUnlock()
		m := make(map[string]float64)
		for _, p := range r.payments {
			if !p.CreatedAt.Before(startDate) && !p.CreatedAt.After(endDate) && p.Status == core.PaymentStatusPaid {
				m[string(p.PaymentMethod)] += p.Amount
			}
		}
		return m, nil
	}

	err := r.db.WithContext(ctx).Model(&core.Payment{}).
		Select("payment_method, COALESCE(SUM(amount), 0) as revenue").
		Where("created_at >= ? AND created_at <= ? AND status = ?", startDate, endDate, core.PaymentStatusPaid).
		Group("payment_method").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	methodMap := make(map[string]float64)
	for _, result := range results {
		methodMap[result.PaymentMethod] = result.Revenue
	}

	return methodMap, nil
}

// PaymentFilters represents filters for payment queries
type PaymentFilters struct {
	OrderID       *uuid.UUID
	Status        *core.PaymentStatus
	PaymentMethod *core.PaymentMethod
	DateFrom      *string
	DateTo        *string
}
