package service

import (
	"context"
	"errors"
	"time"
	branchRepo "service/internal/modules/branches/repository"
	"service/internal/modules/orders/repository"
	serviceRepo "service/internal/modules/services/repository"
	userRepo "service/internal/modules/users/repository"
	"service/internal/shared/model"
	"service/internal/shared/utils"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo        *repository.ServiceOrderRepository
	userRepo         *userRepo.UserRepository
	branchRepo       *branchRepo.BranchRepository
	catalogRepo      *serviceRepo.ServiceCatalogRepository
	providerRepo     *serviceRepo.ServiceProviderRepository
}

// NewOrderService creates a new order service
func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:    repository.NewServiceOrderRepository(),
		userRepo:     userRepo.NewUserRepository(),
		branchRepo:   branchRepo.NewBranchRepository(),
		catalogRepo:  serviceRepo.NewServiceCatalogRepository(),
		providerRepo: serviceRepo.NewServiceProviderRepository(),
	}
}

// CreateOrder creates a new service order
func (s *OrderService) CreateOrder(ctx context.Context, customerID uuid.UUID, req *model.ServiceOrderRequest) (*model.ServiceOrderResponse, error) {
	// Validate customer exists
	_, err := s.userRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	var serviceCatalog *model.ServiceCatalog
	var serviceProvider *model.ServiceProvider
	var branchID *uuid.UUID

	// New multi-service flow: validate ServiceCatalog
	if req.ServiceCatalogID != "" {
		catalogID, err := uuid.Parse(req.ServiceCatalogID)
		if err != nil {
			return nil, errors.New("invalid service catalog ID")
		}

		serviceCatalog, err = s.catalogRepo.GetByID(ctx, catalogID)
		if err != nil {
			return nil, model.ErrCatalogNotFound
		}

		if !serviceCatalog.IsActive {
			return nil, errors.New("service catalog is not active")
		}

		// Validate ServiceProvider if provided
		if req.ServiceProviderID != "" {
			providerID, err := uuid.Parse(req.ServiceProviderID)
			if err != nil {
				return nil, errors.New("invalid service provider ID")
			}

			serviceProvider, err = s.providerRepo.GetByID(ctx, providerID)
			if err != nil {
				return nil, model.ErrProviderNotFound
			}

			if !serviceProvider.IsActive {
				return nil, errors.New("service provider is not active")
			}

			// Use provider's location as service location
			branchID = nil // Provider-based service doesn't need branch
		} else if req.BranchID != "" {
			// Legacy: use branch if provided
			parsedBranchID, err := uuid.Parse(req.BranchID)
			if err != nil {
				return nil, errors.New("invalid branch ID")
			}

			_, err = s.branchRepo.GetByID(ctx, parsedBranchID)
			if err != nil {
				return nil, model.ErrBranchNotFound
			}
			branchID = &parsedBranchID
		}
	} else {
		// Legacy flow: validate branch (backward compatibility)
		if req.BranchID == "" {
			return nil, errors.New("either service_catalog_id or branch_id is required")
		}

		parsedBranchID, err := uuid.Parse(req.BranchID)
		if err != nil {
			return nil, errors.New("invalid branch ID")
		}

		_, err = s.branchRepo.GetByID(ctx, parsedBranchID)
		if err != nil {
			return nil, model.ErrBranchNotFound
		}
		branchID = &parsedBranchID
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

	// Parse appointment date/time if provided
	var appointmentDate *time.Time
	var appointmentTime *time.Time
	if req.AppointmentDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.AppointmentDate)
		if err == nil {
			appointmentDate = &parsedDate
		}
	}
	if req.AppointmentTime != "" {
		parsedTime, err := time.Parse("15:04", req.AppointmentTime)
		if err == nil {
			appointmentTime = &parsedTime
		}
	}

	// Map item fields (prioritize generic fields, fallback to iPhone fields for backward compatibility)
	var itemModel, itemColor, itemSerial, itemType *string
	if req.ItemModel != "" {
		itemModel = &req.ItemModel
	} else if req.IPhoneModel != "" {
		itemModel = &req.IPhoneModel
	}
	if req.ItemColor != "" {
		itemColor = &req.ItemColor
	} else if req.IPhoneColor != "" {
		itemColor = &req.IPhoneColor
	}
	if req.ItemSerial != "" {
		itemSerial = &req.ItemSerial
	} else if req.IPhoneIMEI != "" {
		itemSerial = &req.IPhoneIMEI
	}
	if req.ItemType != "" {
		itemType = &req.ItemType
	} else if req.IPhoneType != "" {
		itemType = &req.IPhoneType
	} else if req.IPhoneModel != "" {
		itemType = &req.IPhoneModel
	}

	// Map iPhone fields for backward compatibility
	var iphoneModel, iphoneColor, iphoneIMEI, iphoneType *string
	if req.IPhoneModel != "" {
		iphoneModel = &req.IPhoneModel
	}
	if req.IPhoneColor != "" {
		iphoneColor = &req.IPhoneColor
	}
	if req.IPhoneIMEI != "" {
		iphoneIMEI = &req.IPhoneIMEI
	}
	if req.IPhoneType != "" {
		iphoneType = &req.IPhoneType
	}

	// Map location fields
	var pickupAddress, pickupLocation *string
	var pickupLatitude, pickupLongitude *float64
	if req.PickupAddress != "" {
		pickupAddress = &req.PickupAddress
	}
	if req.PickupLocation != "" {
		pickupLocation = &req.PickupLocation
	}
	if req.PickupLatitude != nil {
		pickupLatitude = req.PickupLatitude
	}
	if req.PickupLongitude != nil {
		pickupLongitude = req.PickupLongitude
	}

	// Determine service location and on-demand status
	serviceLocation := req.ServiceLocation
	isOnDemand := false
	
	if serviceCatalog != nil {
		// If RequiresLocation is false, service datang ke customer (on-demand)
		isOnDemand = !serviceCatalog.RequiresLocation
		
		if isOnDemand {
			// Service datang ke customer, use customer location
			if req.PickupAddress != "" {
				serviceLocation = req.PickupAddress
			} else if req.PickupLocation != "" {
				serviceLocation = req.PickupLocation
			}
		} else {
			// Service di lokasi provider/branch
			if serviceLocation == "" && serviceProvider != nil {
				serviceLocation = serviceProvider.Address
			} else if serviceLocation == "" && branchID != nil {
				branch, _ := s.branchRepo.GetByID(ctx, *branchID)
				if branch != nil {
					serviceLocation = branch.Address
				}
			}
		}
	} else {
		// Legacy: use provider/branch location
		if serviceLocation == "" && serviceProvider != nil {
			serviceLocation = serviceProvider.Address
		} else if serviceLocation == "" && branchID != nil {
			branch, _ := s.branchRepo.GetByID(ctx, *branchID)
			if branch != nil {
				serviceLocation = branch.Address
			}
		}
	}

	// Determine service type and name
	serviceType := req.ServiceType
	serviceName := ""
	if serviceCatalog != nil {
		serviceName = serviceCatalog.Name
		if serviceType == "" {
			serviceType = model.ServiceTypeOther // Default if not provided
		}
	}

	// Determine initial status based on service requirements
	initialStatus := model.StatusPendingPickup
	if serviceCatalog != nil && serviceCatalog.RequiresAppointment {
		initialStatus = model.StatusPendingPickup // Will be updated when appointment is confirmed
	}

	// Determine estimated cost and duration
	estimatedCost := req.EstimatedCost
	estimatedDuration := req.EstimatedDuration
	if serviceCatalog != nil {
		if estimatedCost == 0 {
			estimatedCost = serviceCatalog.BasePrice
		}
		if estimatedDuration == 0 {
			estimatedDuration = serviceCatalog.EstimatedDuration
		}
	}

	// Set ServiceCatalogID and ServiceProviderID
	var catalogID, providerID *uuid.UUID
	if serviceCatalog != nil {
		catalogID = &serviceCatalog.ID
	}
	if serviceProvider != nil {
		providerID = &serviceProvider.ID
	}

	// Create order entity
	order := &model.ServiceOrder{
		OrderNumber:       orderNumber,
		CustomerID:        customerID,
		ServiceCatalogID:  catalogID,
		ServiceProviderID:  providerID,
		BranchID:          branchID,
		IPhoneModel:       iphoneModel,
		IPhoneColor:       iphoneColor,
		IPhoneIMEI:        iphoneIMEI,
		IPhoneType:        iphoneType,
		ItemModel:         itemModel,
		ItemColor:         itemColor,
		ItemSerial:        itemSerial,
		ItemType:          itemType,
		ServiceType:       serviceType,
		ServiceName:       serviceName,
		Description:       req.Description,
		Complaint:         req.Complaint,
		AppointmentDate:   appointmentDate,
		AppointmentTime:   appointmentTime,
		PickupAddress:     pickupAddress,
		PickupLocation:    pickupLocation,
		PickupLatitude:    pickupLatitude,
		PickupLongitude:   pickupLongitude,
		ServiceLocation:   serviceLocation,
		IsOnDemand:        isOnDemand,
		Status:            initialStatus,
		EstimatedCost:     estimatedCost,
		ActualCost:        0,
		EstimatedDuration: estimatedDuration,
		ActualDuration:    0,
		Metadata:          req.Metadata,
	}

	// Set alias fields for backward compatibility
	order.SetAliasFields()

	// Save to database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Reload with relations
	order, err = s.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	// Return response with populated data
	response := order.ToResponse()
	return &response, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*model.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	response := order.ToResponse()
	return &response, nil
}

// GetOrderByNumber retrieves an order by order number
func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*model.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	response := order.ToResponse()
	return &response, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, req *model.UpdateOrderStatusRequest) (*model.ServiceOrderResponse, error) {
	// Get existing order
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrOrderNotFound
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
func (s *OrderService) AssignTechnician(ctx context.Context, orderID, technicianID uuid.UUID) (*model.ServiceOrderResponse, error) {
	// Validate technician exists and has correct role
	technician, err := s.userRepo.GetByID(ctx, technicianID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	if string(technician.Role) != string(model.RoleTeknisi) {
		return nil, errors.New("user is not a technician")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	// Assign technician
	if err := s.orderRepo.AssignTechnician(ctx, orderID, technicianID); err != nil {
		return nil, err
	}

	// Update order status if needed
	if order.Status == model.StatusPendingPickup {
		order.Status = model.StatusInService
		s.orderRepo.Update(ctx, order)
	}

	response := order.ToResponse()
	return &response, nil
}

// AssignCourier assigns a courier to an order
func (s *OrderService) AssignCourier(ctx context.Context, orderID, courierID uuid.UUID) (*model.ServiceOrderResponse, error) {
	// Validate courier exists and has correct role
	courier, err := s.userRepo.GetByID(ctx, courierID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	if string(courier.Role) != string(model.RoleKurir) {
		return nil, errors.New("user is not a courier")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	// Assign courier
	if err := s.orderRepo.AssignCourier(ctx, orderID, courierID); err != nil {
		return nil, err
	}

	response := order.ToResponse()
	return &response, nil
}

// ListOrders retrieves orders with pagination and filters
func (s *OrderService) ListOrders(ctx context.Context, page, limit int, filters *repository.ServiceOrderFilters) (*model.PaginatedResponse, error) {
	offset := (page - 1) * limit

	orders, total, err := s.orderRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var responses []model.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	// Calculate pagination
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	pagination := model.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return &model.PaginatedResponse{
		Status:     "success",
		Data:       responses,
		Pagination: pagination,
		Message:    "Orders retrieved successfully",
		Timestamp:  model.GetCurrentTimestamp(),
	}, nil
}

// GetOrdersByCustomer retrieves orders for a specific customer
func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]model.ServiceOrderResponse, error) {
	orders, err := s.orderRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	var responses []model.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, nil
}

// GetOrdersByBranch retrieves orders for a specific branch with pagination
func (s *OrderService) GetOrdersByBranch(ctx context.Context, branchID uuid.UUID, page, limit int) ([]model.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{BranchID: &branchID})
	if err != nil {
		return nil, 0, err
	}

	var responses []model.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, total, nil
}

// GetOrdersByStatus retrieves orders by status
func (s *OrderService) GetOrdersByStatus(ctx context.Context, status model.OrderStatus) ([]model.ServiceOrderResponse, error) {
	orders, err := s.orderRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	var responses []model.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, nil
}

// GetOrdersByTechnician retrieves orders assigned to a technician with pagination
func (s *OrderService) GetOrdersByTechnician(ctx context.Context, technicianID uuid.UUID, page, limit int) ([]model.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{TechnicianID: &technicianID})
	if err != nil {
		return nil, 0, err
	}

	var responses []model.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, total, nil
}

// GetOrdersByCourier retrieves orders assigned to a courier with pagination
func (s *OrderService) GetOrdersByCourier(ctx context.Context, courierID uuid.UUID, page, limit int) ([]model.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{CourierID: &courierID})
	if err != nil {
		return nil, 0, err
	}

	var responses []model.ServiceOrderResponse
	for _, order := range orders {
		responses = append(responses, order.ToResponse())
	}

	return responses, total, nil
}

// GetOrders retrieves orders for a user with pagination and optional filters
func (s *OrderService) GetOrders(ctx context.Context, userID uuid.UUID, page, limit int, status string, branchID *uuid.UUID) ([]model.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	filters := &repository.ServiceOrderFilters{}
	if status != "" {
		st := model.OrderStatus(status)
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

	var responses []model.ServiceOrderResponse
	for _, o := range orders {
		responses = append(responses, o.ToResponse())
	}

	return responses, total, nil
}

// GetAllOrders retrieves all orders (admin) with pagination and optional filters
func (s *OrderService) GetAllOrders(ctx context.Context, page, limit int, status string, branchID *uuid.UUID) ([]model.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	filters := &repository.ServiceOrderFilters{}
	if status != "" {
		st := model.OrderStatus(status)
		filters.Status = &st
	}
	if branchID != nil {
		filters.BranchID = branchID
	}

	orders, total, err := s.orderRepo.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.ServiceOrderResponse
	for _, o := range orders {
		responses = append(responses, o.ToResponse())
	}

	return responses, total, nil
}

// UpdateOrder updates an order's details (admin)
func (s *OrderService) UpdateOrder(ctx context.Context, id uuid.UUID, req *model.ServiceOrderRequest) (*model.ServiceOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrOrderNotFound
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
		return model.ErrOrderNotFound
	}
	return s.orderRepo.Delete(ctx, id)
}

// GetOrdersByBranchID retrieves orders for a specific branch with pagination
func (s *OrderService) GetOrdersByBranchID(ctx context.Context, branchID uuid.UUID, page, limit int) ([]model.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{BranchID: &branchID})
	if err != nil {
		return nil, 0, err
	}
	var responses []model.ServiceOrderResponse
	for _, o := range orders {
		responses = append(responses, o.ToResponse())
	}
	return responses, total, nil
}

// GetAvailableJobs returns orders which are available for couriers (pending and not assigned)
func (s *OrderService) GetAvailableJobs(ctx context.Context, page, limit int) ([]model.ServiceOrderResponse, int64, error) {
	offset := (page - 1) * limit
	status := model.StatusPendingPickup
	orders, total, err := s.orderRepo.List(ctx, offset, limit, &repository.ServiceOrderFilters{Status: &status})
	if err != nil {
		return nil, 0, err
	}
	var responses []model.ServiceOrderResponse
	for _, o := range orders {
		if o.CourierID == nil {
			responses = append(responses, o.ToResponse())
		}
	}
	return responses, total, nil
}

// UpdateOrderCost updates the cost information of an order
func (s *OrderService) UpdateOrderCost(ctx context.Context, id uuid.UUID, estimatedCost, actualCost float64, estimatedDuration, actualDuration int) (*model.ServiceOrderResponse, error) {
	// Get existing order
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrOrderNotFound
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
