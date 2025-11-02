package service

import (
	"context"
	"errors"
<<<<<<< HEAD
	"service/internal/core"
	"service/internal/orders/repository"
	"service/internal/utils"
=======
<<<<<<<< HEAD:internal/service/order_service.go
	"service/internal/core"
	"service/internal/orders/repository"
	"service/internal/utils"
========
	branchRepo "service/internal/branches/repository"
	"service/internal/orders/dto"
	"service/internal/orders/repository"
	"service/internal/shared/model"
	"service/internal/shared/utils"
	userRepo "service/internal/users/repository"
>>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b:internal/orders/service/order_service.go
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b

	"github.com/google/uuid"
)

// OrderService handles order business logic
type OrderService struct {
	orderRepo  *repository.ServiceOrderRepository
<<<<<<< HEAD
	userRepo   *repository.UserRepository
	branchRepo *repository.BranchRepository
=======
	userRepo   *userRepo.UserRepository
	branchRepo *branchRepo.BranchRepository
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
}

// NewOrderService creates a new order service
func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:  repository.NewServiceOrderRepository(),
<<<<<<< HEAD
		userRepo:   repository.NewUserRepository(),
		branchRepo: repository.NewBranchRepository(),
=======
		userRepo:   userRepo.NewUserRepository(),
		branchRepo: branchRepo.NewBranchRepository(),
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}
}

// CreateOrder creates a new service order
<<<<<<< HEAD
func (s *OrderService) CreateOrder(ctx context.Context, customerID uuid.UUID, req *core.ServiceOrderRequest) (*core.ServiceOrderResponse, error) {
	// Validate customer exists
	_, err := s.userRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, core.ErrUserNotFound
=======
func (s *OrderService) CreateOrder(ctx context.Context, customerID uuid.UUID, req *dto.ServiceOrderRequest) (*dto.ServiceOrderResponse, error) {
	// Validate customer exists
	_, err := s.userRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, model.ErrUserNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	// Validate branch exists
	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch ID")
	}

	_, err = s.branchRepo.GetByID(ctx, branchID)
	if err != nil {
<<<<<<< HEAD
		return nil, core.ErrBranchNotFound
=======
		return nil, model.ErrBranchNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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
<<<<<<< HEAD
	order := &core.ServiceOrder{
=======
	order := &dto.ServiceOrder{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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
<<<<<<< HEAD
		Status:            core.StatusPendingPickup,
=======
		Status:            dto.StatusPendingPickup,
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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
<<<<<<< HEAD
func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*core.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrOrderNotFound
=======
func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*dto.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	response := order.ToResponse()
	return &response, nil
}

// GetOrderByNumber retrieves an order by order number
<<<<<<< HEAD
func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*core.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, core.ErrOrderNotFound
=======
func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*dto.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	response := order.ToResponse()
	return &response, nil
}

// UpdateOrderStatus updates the status of an order
<<<<<<< HEAD
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, req *core.UpdateOrderStatusRequest) (*core.ServiceOrderResponse, error) {
	// Get existing order
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrOrderNotFound
=======
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, req *dto.UpdateOrderStatusRequest) (*dto.ServiceOrderResponse, error) {
	// Get existing order
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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
<<<<<<< HEAD
func (s *OrderService) AssignTechnician(ctx context.Context, orderID, technicianID uuid.UUID) (*core.ServiceOrderResponse, error) {
	// Validate technician exists and has correct role
	technician, err := s.userRepo.GetByID(ctx, technicianID)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	if technician.Role != core.RoleTeknisi {
=======
func (s *OrderService) AssignTechnician(ctx context.Context, orderID, technicianID uuid.UUID) (*dto.ServiceOrderResponse, error) {
	// Validate technician exists and has correct role
	technician, err := s.userRepo.GetByID(ctx, technicianID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	if string(technician.Role) != string(model.RoleTeknisi) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		return nil, errors.New("user is not a technician")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
<<<<<<< HEAD
		return nil, core.ErrOrderNotFound
=======
		return nil, model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	// Assign technician
	if err := s.orderRepo.AssignTechnician(ctx, orderID, technicianID); err != nil {
		return nil, err
	}

	// Update order status if needed
<<<<<<< HEAD
	if order.Status == core.StatusPendingPickup {
		order.Status = core.StatusInService
=======
	if order.Status == dto.StatusPendingPickup {
		order.Status = dto.StatusInService
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		s.orderRepo.Update(ctx, order)
	}

	response := order.ToResponse()
	return &response, nil
}

// AssignCourier assigns a courier to an order
<<<<<<< HEAD
func (s *OrderService) AssignCourier(ctx context.Context, orderID, courierID uuid.UUID) (*core.ServiceOrderResponse, error) {
	// Validate courier exists and has correct role
	courier, err := s.userRepo.GetByID(ctx, courierID)
	if err != nil {
		return nil, core.ErrUserNotFound
	}

	if courier.Role != core.RoleKurir {
=======
func (s *OrderService) AssignCourier(ctx context.Context, orderID, courierID uuid.UUID) (*dto.ServiceOrderResponse, error) {
	// Validate courier exists and has correct role
	courier, err := s.userRepo.GetByID(ctx, courierID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	if string(courier.Role) != string(model.RoleKurir) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		return nil, errors.New("user is not a courier")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
<<<<<<< HEAD
		return nil, core.ErrOrderNotFound
=======
		return nil, model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	// Assign courier
	if err := s.orderRepo.AssignCourier(ctx, orderID, courierID); err != nil {
		return nil, err
	}

	response := order.ToResponse()
	return &response, nil
}

// ListOrders retrieves orders with pagination and filters
<<<<<<< HEAD
func (s *OrderService) ListOrders(ctx context.Context, page, limit int, filters *repository.ServiceOrderFilters) (*core.PaginatedResponse, error) {
=======
func (s *OrderService) ListOrders(ctx context.Context, page, limit int, filters *repository.ServiceOrderFilters) (*model.PaginatedResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	offset := (page - 1) * limit

	orders, total, err := s.orderRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, err
	}

	// Convert to response format
<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	// Calculate pagination
	totalPages := int((total + int64(limit) - 1) / int64(limit))
<<<<<<< HEAD
	pagination := core.PaginationResponse{
=======
	pagination := model.PaginationResponse{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

<<<<<<< HEAD
	return &core.PaginatedResponse{
=======
	return &model.PaginatedResponse{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		Status:     "success",
		Data:       responses,
		Pagination: pagination,
		Message:    "Orders retrieved successfully",
<<<<<<< HEAD
		Timestamp:  core.GetCurrentTimestamp(),
=======
		Timestamp:  model.GetCurrentTimestamp(),
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}, nil
}

// GetOrdersByCustomer retrieves orders for a specific customer
<<<<<<< HEAD
func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]core.ServiceOrderResponse, error) {
=======
func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]dto.ServiceOrderResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	orders, err := s.orderRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, nil
}

// GetOrdersByBranch retrieves orders for a specific branch with pagination
<<<<<<< HEAD
func (s *OrderService) GetOrdersByBranch(ctx context.Context, branchID uuid.UUID, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
=======
func (s *OrderService) GetOrdersByBranch(ctx context.Context, branchID uuid.UUID, page, limit int) ([]dto.ServiceOrderResponse, int64, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{BranchID: &branchID})
	if err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, total, nil
}

// GetOrdersByStatus retrieves orders by status
<<<<<<< HEAD
func (s *OrderService) GetOrdersByStatus(ctx context.Context, status core.OrderStatus) ([]core.ServiceOrderResponse, error) {
=======
func (s *OrderService) GetOrdersByStatus(ctx context.Context, status dto.OrderStatus) ([]dto.ServiceOrderResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	orders, err := s.orderRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, nil
}

// GetOrdersByTechnician retrieves orders assigned to a technician with pagination
<<<<<<< HEAD
func (s *OrderService) GetOrdersByTechnician(ctx context.Context, technicianID uuid.UUID, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
=======
func (s *OrderService) GetOrdersByTechnician(ctx context.Context, technicianID uuid.UUID, page, limit int) ([]dto.ServiceOrderResponse, int64, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{TechnicianID: &technicianID})
	if err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, total, nil
}

// GetOrdersByCourier retrieves orders assigned to a courier with pagination
<<<<<<< HEAD
func (s *OrderService) GetOrdersByCourier(ctx context.Context, courierID uuid.UUID, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
=======
func (s *OrderService) GetOrdersByCourier(ctx context.Context, courierID uuid.UUID, page, limit int) ([]dto.ServiceOrderResponse, int64, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{CourierID: &courierID})
	if err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, total, nil
}

// GetOrders retrieves orders for a user with pagination and optional filters
<<<<<<< HEAD
func (s *OrderService) GetOrders(ctx context.Context, userID uuid.UUID, page, limit int, status string, branchID *uuid.UUID) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	filters := &repository.ServiceOrderFilters{}
	if status != "" {
		st := core.OrderStatus(status)
=======
func (s *OrderService) GetOrders(ctx context.Context, userID uuid.UUID, page, limit int, status string, branchID *uuid.UUID) ([]dto.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	filters := &repository.ServiceOrderFilters{}
	if status != "" {
		st := dto.OrderStatus(status)
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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

<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, o := range orders {
		responses = append(responses, o.ToResponse())
	}

	return responses, total, nil
}

// GetAllOrders retrieves all orders (admin) with pagination and optional filters
<<<<<<< HEAD
func (s *OrderService) GetAllOrders(ctx context.Context, page, limit int, status string, branchID *uuid.UUID) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	filters := &repository.ServiceOrderFilters{}
	if status != "" {
		st := core.OrderStatus(status)
=======
func (s *OrderService) GetAllOrders(ctx context.Context, page, limit int, status string, branchID *uuid.UUID) ([]dto.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	filters := &repository.ServiceOrderFilters{}
	if status != "" {
		st := dto.OrderStatus(status)
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		filters.Status = &st
	}
	if branchID != nil {
		filters.BranchID = branchID
	}

	orders, total, err := s.orderRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, o := range orders {
		responses = append(responses, o.ToResponse())
	}

	return responses, total, nil
}

// UpdateOrder updates an order's details (admin)
<<<<<<< HEAD
func (s *OrderService) UpdateOrder(ctx context.Context, id uuid.UUID, req *core.ServiceOrderRequest) (*core.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrOrderNotFound
=======
func (s *OrderService) UpdateOrder(ctx context.Context, id uuid.UUID, req *dto.ServiceOrderRequest) (*dto.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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
<<<<<<< HEAD
		return core.ErrOrderNotFound
=======
		return model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}
	return s.orderRepo.Delete(ctx, id)
}

// GetOrdersByBranchID retrieves orders for a specific branch with pagination
<<<<<<< HEAD
func (s *OrderService) GetOrdersByBranchID(ctx context.Context, branchID uuid.UUID, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
=======
func (s *OrderService) GetOrdersByBranchID(ctx context.Context, branchID uuid.UUID, page, limit int) ([]dto.ServiceOrderResponse, int64, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{BranchID: &branchID})
	if err != nil {
		return nil, 0, err
	}
<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, o := range orders {
		responses = append(responses, o.ToResponse())
	}
	return responses, total, nil
}

// GetAvailableJobs returns orders which are available for couriers (pending and not assigned)
<<<<<<< HEAD
func (s *OrderService) GetAvailableJobs(ctx context.Context, page, limit int) ([]core.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	status := core.StatusPendingPickup
=======
func (s *OrderService) GetAvailableJobs(ctx context.Context, page, limit int) ([]dto.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	status := dto.StatusPendingPickup
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{Status: &status})
	if err != nil {
		return nil, 0, err
	}
<<<<<<< HEAD
	var responses []core.ServiceOrderResponse
=======
	var responses []dto.ServiceOrderResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, o := range orders {
		if o.CourierID == nil {
			responses = append(responses, o.ToResponse())
		}
	}
	return responses, total, nil
}

// UpdateOrderCost updates the cost information of an order
<<<<<<< HEAD
func (s *OrderService) UpdateOrderCost(ctx context.Context, id uuid.UUID, estimatedCost, actualCost float64, estimatedDuration, actualDuration int) (*core.ServiceOrderResponse, error) {
	// Get existing order
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrOrderNotFound
=======
func (s *OrderService) UpdateOrderCost(ctx context.Context, id uuid.UUID, estimatedCost, actualCost float64, estimatedDuration, actualDuration int) (*dto.ServiceOrderResponse, error) {
	// Get existing order
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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
