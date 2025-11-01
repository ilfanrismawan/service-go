package service

import (
	"context"
	"errors"
	"service/internal/core"
	"service/internal/orders/repository"
	"service/internal/utils"

	"github.com/google/uuid"
)

// OrderService handles order business logic
type OrderService struct {
	orderRepo  *repository.ServiceOrderRepository
	userRepo   *repository.UserRepository
	branchRepo *repository.BranchRepository
}

// NewOrderService creates a new order service
func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:  repository.NewServiceOrderRepository(),
		userRepo:   repository.NewUserRepository(),
		branchRepo: repository.NewBranchRepository(),
	}
}

// CreateOrder creates a new service order
func (s *OrderService) CreateOrder(ctx context.Context, customerID uuid.UUID, req *core.ServiceOrderRequest) (*core.ServiceOrderResponse, error) {
	// Validate customer exists
	_, err := s.userRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	// Validate branch exists
	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch ID")
	}

	_, err = s.branchRepo.GetByID(ctx, branchID)
	if err != nil {
		return nil, core.ErrBranchNotFound
	}

	// Generate unique order number
	orderNumber := utils.GenerateOrderNumber()

	// Check if order number already exists (very unlikely but safety check)
	exists, err := s.orderRepo.CheckOrderNumberExists(ctx, orderNumber, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		// Regenerate if exists
		orderNumber = utils.GenerateOrderNumber()
	}

	// Create order entity
	order := &core.ServiceOrder{
		OrderNumber:       orderNumber,
		CustomerID:        customerID,
		BranchID:          branchID,
		IPhoneModel:       req.IPhoneModel,
		IPhoneColor:       req.IPhoneColor,
		IPhoneIMEI:        req.IPhoneIMEI,
		ServiceType:       req.ServiceType,
		Description:       req.Description,
		PickupAddress:     req.PickupAddress,
		PickupLatitude:    req.PickupLatitude,
		PickupLongitude:   req.PickupLongitude,
		Status:            core.StatusPendingPickup,
		EstimatedCost:     0,
		ActualCost:        0,
		EstimatedDuration: 0,
		ActualDuration:    0,
	}

	// Save to database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Return response with populated data
	response := order.ToResponse()
	return &response, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*core.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrOrderNotFound
	}

	response := order.ToResponse()
	return &response, nil
}

// GetOrderByNumber retrieves an order by order number
func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*core.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, core.ErrOrderNotFound
	}

	response := order.ToResponse()
	return &response, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, req *core.UpdateOrderStatusRequest) (*core.ServiceOrderResponse, error) {
	// Get existing order
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrOrderNotFound
	}

	// Update status
	order.Status = req.Status
	if req.Notes != "" {
		order.Notes = req.Notes
	}

	// Save changes
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	response := order.ToResponse()
	return &response, nil
}

// AssignTechnician assigns a technician to an order
func (s *OrderService) AssignTechnician(ctx context.Context, orderID, technicianID uuid.UUID) (*core.ServiceOrderResponse, error) {
	// Validate technician exists and has correct role
	technician, err := s.userRepo.GetByID(ctx, technicianID)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	if technician.Role != core.RoleTeknisi {
		return nil, errors.New("user is not a technician")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, core.ErrOrderNotFound
	}

	// Assign technician
	if err := s.orderRepo.AssignTechnician(ctx, orderID, technicianID); err != nil {
		return nil, err
	}

	// Update order status if needed
	if order.Status == core.StatusPendingPickup {
		order.Status = core.StatusInService
		s.orderRepo.Update(ctx, order)
	}

	response := order.ToResponse()
	return &response, nil
}

// AssignCourier assigns a courier to an order
func (s *OrderService) AssignCourier(ctx context.Context, orderID, courierID uuid.UUID) (*core.ServiceOrderResponse, error) {
	// Validate courier exists and has correct role
	courier, err := s.userRepo.GetByID(ctx, courierID)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	if courier.Role != core.RoleKurir {
		return nil, errors.New("user is not a courier")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, core.ErrOrderNotFound
	}

	// Assign courier
	if err := s.orderRepo.AssignCourier(ctx, orderID, courierID); err != nil {
		return nil, err
	}

	response := order.ToResponse()
	return &response, nil
}

// ListOrders retrieves orders with pagination and filters
func (s *OrderService) ListOrders(ctx context.Context, page, limit int, filters *repository.ServiceOrderFilters) (*core.PaginatedResponse, error) {
	offset := (page - 1) * limit

	orders, total, err := s.orderRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var responses []core.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	// Calculate pagination
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	pagination := core.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return &core.PaginatedResponse{
		Status:     "success",
		Data:       responses,
		Pagination: pagination,
		Message:    "Orders retrieved successfully",
		Timestamp:  core.GetCurrentTimestamp(),
	}, nil
}

// GetOrdersByCustomer retrieves orders for a specific customer
func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]core.ServiceOrderResponse, error) {
	orders, err := s.orderRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	var responses []core.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, nil
}

// GetOrdersByBranch retrieves orders for a specific branch with pagination
func (s *OrderService) GetOrdersByBranch(ctx context.Context, branchID uuid.UUID, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{BranchID: &branchID})
	if err != nil {
		return nil, 0, err
	}

	var responses []core.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, total, nil
}

// GetOrdersByStatus retrieves orders by status
func (s *OrderService) GetOrdersByStatus(ctx context.Context, status core.OrderStatus) ([]core.ServiceOrderResponse, error) {
	orders, err := s.orderRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	var responses []core.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, nil
}

// GetOrdersByTechnician retrieves orders assigned to a technician with pagination
func (s *OrderService) GetOrdersByTechnician(ctx context.Context, technicianID uuid.UUID, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{TechnicianID: &technicianID})
	if err != nil {
		return nil, 0, err
	}

	var responses []core.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, total, nil
}

// GetOrdersByCourier retrieves orders assigned to a courier with pagination
func (s *OrderService) GetOrdersByCourier(ctx context.Context, courierID uuid.UUID, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{CourierID: &courierID})
	if err != nil {
		return nil, 0, err
	}

	var responses []core.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, total, nil
}

// GetOrders retrieves orders for a user with pagination and optional filters
func (s *OrderService) GetOrders(ctx context.Context, userID uuid.UUID, page, limit int, status string, branchID *uuid.UUID) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	filters := &repository.ServiceOrderFilters{}
	if status != "" {
		st := core.OrderStatus(status)
		filters.Status = &st
	}
	if branchID != nil {
		filters.BranchID = branchID
	} else {
		// attempt to resolve user's branch
		if user, err := s.userRepo.GetByID(ctx, userID); err == nil && user.BranchID != nil {
			filters.BranchID = user.BranchID
		}
	}

	orders, total, err := s.orderRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	var responses []core.ServiceOrderResponse
	for _, o := range orders {
		responses = append(responses, o.ToResponse())
	}

	return responses, total, nil
}

// GetAllOrders retrieves all orders (admin) with pagination and optional filters
func (s *OrderService) GetAllOrders(ctx context.Context, page, limit int, status string, branchID *uuid.UUID) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	filters := &repository.ServiceOrderFilters{}
	if status != "" {
		st := core.OrderStatus(status)
		filters.Status = &st
	}
	if branchID != nil {
		filters.BranchID = branchID
	}

	orders, total, err := s.orderRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	var responses []core.ServiceOrderResponse
	for _, o := range orders {
		responses = append(responses, o.ToResponse())
	}

	return responses, total, nil
}

// UpdateOrder updates an order's details (admin)
func (s *OrderService) UpdateOrder(ctx context.Context, id uuid.UUID, req *core.ServiceOrderRequest) (*core.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrOrderNotFound
	}

	if req.BranchID != "" {
		if bid, err := uuid.Parse(req.BranchID); err == nil {
			order.BranchID = bid
		}
	}
	if req.IPhoneModel != "" {
		order.IPhoneModel = req.IPhoneModel
	}
	if req.IPhoneColor != "" {
		order.IPhoneColor = req.IPhoneColor
	}
	if req.IPhoneIMEI != "" {
		order.IPhoneIMEI = req.IPhoneIMEI
	}
	if req.ServiceType != "" {
		order.ServiceType = req.ServiceType
	}
	if req.Description != "" {
		order.Description = req.Description
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	resp := order.ToResponse()
	return &resp, nil
}

// DeleteOrder deletes an order (admin)
func (s *OrderService) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	if _, err := s.orderRepo.GetByID(ctx, id); err != nil {
		return core.ErrOrderNotFound
	}
	return s.orderRepo.Delete(ctx, id)
}

// GetOrdersByBranchID retrieves orders for a specific branch with pagination
func (s *OrderService) GetOrdersByBranchID(ctx context.Context, branchID uuid.UUID, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{BranchID: &branchID})
	if err != nil {
		return nil, 0, err
	}
	var responses []core.ServiceOrderResponse
	for _, o := range orders {
		responses = append(responses, o.ToResponse())
	}
	return responses, total, nil
}

// GetAvailableJobs returns orders which are available for couriers (pending and not assigned)
func (s *OrderService) GetAvailableJobs(ctx context.Context, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	status := core.StatusPendingPickup
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{Status: &status})
	if err != nil {
		return nil, 0, err
	}
	var responses []core.ServiceOrderResponse
	for _, o := range orders {
		if o.CourierID == nil {
			responses = append(responses, o.ToResponse())
		}
	}
	return responses, total, nil
}

// UpdateOrderCost updates the cost information of an order
func (s *OrderService) UpdateOrderCost(ctx context.Context, id uuid.UUID, estimatedCost, actualCost float64, estimatedDuration, actualDuration int) (*core.ServiceOrderResponse, error) {
	// Get existing order
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrOrderNotFound
	}

	// Update cost information
	order.EstimatedCost = estimatedCost
	order.ActualCost = actualCost
	order.EstimatedDuration = estimatedDuration
	order.ActualDuration = actualDuration

	// Save changes
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	response := order.ToResponse()
	return &response, nil
}
