package service

import (
	"context"
	dashboardDTO "service-go/internal/modules/admin/dto"
	branchRepository "service-go/internal/modules/branches/repository"
	orderEntity "service-go/internal/modules/orders/entity"
	orderRepository "service-go/internal/modules/orders/repository"
	paymentEntity "service-go/internal/modules/payments/entity"
	paymentRepository "service-go/internal/modules/payments/repository"
	userEntity "service-go/internal/modules/users/entity"
	userRepository "service-go/internal/modules/users/repository"
	"service-go/internal/shared/model"

	"github.com/google/uuid"
)

// DashboardService handles dashboard business logic
type DashboardService struct {
	orderRepo   *orderRepository.ServiceOrderRepository
	userRepo    *userRepository.UserRepository
	branchRepo  *branchRepository.BranchRepository
	paymentRepo *paymentRepository.PaymentRepository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService() *DashboardService {
	return &DashboardService{
		orderRepo:   orderRepository.NewServiceOrderRepository(),
		userRepo:    userRepository.NewUserRepository(),
		branchRepo:  branchRepository.NewBranchRepository(),
		paymentRepo: paymentRepository.NewPaymentRepository(),
	}
}

// GetDashboardStats retrieves dashboard statistics
func (s *DashboardService) GetDashboardStats(ctx context.Context, userID *uuid.UUID, branchID *uuid.UUID) (*dashboardDTO.DashboardStats, error) {
	// Get user role if userID is provided
	var userRole *userEntity.UserRole
	if userID != nil {
		user, err := s.userRepo.GetByID(ctx, *userID)
		if err != nil {
			return nil, err
		}
		userRole = &user.Role
	}

	// Set filters based on user role
	filters := &orderRepository.ServiceOrderFilters{}
	if branchID != nil {
		filters.BranchID = branchID
	} else if userRole != nil && *userRole != userEntity.RoleAdminPusat {
		// Non-admin users can only see their branch data
		if userRole != nil {
			user, _ := s.userRepo.GetByID(ctx, *userID)
			if user.BranchID != nil {
				filters.BranchID = user.BranchID
			}
		}
	}

	// Get order statistics
	orders, _, err := s.orderRepo.List(ctx, 0, 1000, filters)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	stats := &dashboardDTO.DashboardStats{
		TotalOrders:     int64(len(orders)),
		TotalRevenue:    0,
		PendingOrders:   0,
		CompletedOrders: 0,
		ActiveCustomers: 0,
		ActiveBranches:  0,
	}

	// Count orders by status
	for _, order := range orders {
		switch order.Status {
		case orderEntity.StatusPendingPickup, orderEntity.StatusOnPickup, orderEntity.StatusInService, orderEntity.StatusReady, orderEntity.StatusDelivered:
			stats.PendingOrders++
		case orderEntity.StatusCompleted:
			stats.CompletedOrders++
		}
	}

	// Get payment statistics
	paymentFilters := &paymentRepository.PaymentFilters{}
	if branchID != nil {
		paymentFilters.OrderID = branchID // This would need to be adjusted based on your payment structure
	}

	payments, _, err := s.paymentRepo.List(ctx, 0, 1000, paymentFilters)
	if err != nil {
		return nil, err
	}

	// Calculate total revenue
	for _, payment := range payments {
		if payment.Status == paymentEntity.PaymentStatusPaid {
			stats.TotalRevenue += payment.Amount
		}
	}

	// Get active customers count
	if branchID != nil {
		// Count customers for specific branch
		customers, err := s.userRepo.GetByBranchID(ctx, *branchID)
		if err == nil {
			stats.ActiveCustomers = int64(len(customers))
		}
	} else {
		// Count all active customers
		role := userEntity.RolePelanggan
		customers, _, err := s.userRepo.List(ctx, 0, 1000, &role, nil)
		if err == nil {
			stats.ActiveCustomers = int64(len(customers))
		}
	}

	// Get active branches count
	branches, err := s.branchRepo.GetActiveBranches(ctx)
	if err == nil {
		stats.ActiveBranches = int64(len(branches))
	}

	return stats, nil
}

// GetServiceStats retrieves service type statistics
func (s *DashboardService) GetServiceStats(ctx context.Context, userID *uuid.UUID, branchID *uuid.UUID) ([]model.ServiceStats, error) {
	// Set filters based on user role
	filters := &orderRepository.ServiceOrderFilters{}
	if branchID != nil {
		filters.BranchID = branchID
	} else if userID != nil {
		user, err := s.userRepo.GetByID(ctx, *userID)
		if err == nil && user.Role != userEntity.RoleAdminPusat && user.BranchID != nil {
			filters.BranchID = user.BranchID
		}
	}

	// Get all orders
	orders, _, err := s.orderRepo.List(ctx, 0, 1000, filters)
	if err != nil {
		return nil, err
	}

	// Group by service type
	serviceStatsMap := make(map[orderEntity.ServiceType]*model.ServiceStats)
	for _, order := range orders {
		if order.Status == orderEntity.StatusCompleted {
			if stats, exists := serviceStatsMap[order.ServiceType]; exists {
				stats.Count++
				stats.Revenue += order.ActualCost
			} else {
				serviceStatsMap[order.ServiceType] = &model.ServiceStats{
					ServiceType: string(order.ServiceType),
					Count:       1,
					Revenue:     order.ActualCost,
				}
			}
		}
	}

	// Convert map to slice
	var serviceStats []model.ServiceStats
	for _, stats := range serviceStatsMap {
		serviceStats = append(serviceStats, *stats)
	}

	return serviceStats, nil
}

// GetRevenueReport retrieves revenue report for a date range
func (s *DashboardService) GetRevenueReport(ctx context.Context, userID *uuid.UUID, branchID *uuid.UUID, dateFrom, dateTo string) ([]model.RevenueReport, error) {
	// Set filters based on user role
	filters := &orderRepository.ServiceOrderFilters{}
	if branchID != nil {
		filters.BranchID = branchID
	} else if userID != nil {
		user, err := s.userRepo.GetByID(ctx, *userID)
		if err == nil && user.Role != userEntity.RoleAdminPusat && user.BranchID != nil {
			filters.BranchID = user.BranchID
		}
	}

	filters.DateFrom = &dateFrom
	filters.DateTo = &dateTo

	// Get orders in date range
	orders, _, err := s.orderRepo.List(ctx, 0, 1000, filters)
	if err != nil {
		return nil, err
	}

	// Group by date
	revenueMap := make(map[string]*model.RevenueReport)
	for _, order := range orders {
		if order.Status == orderEntity.StatusCompleted {
			date := order.CreatedAt.Format("2006-01-02")
			if report, exists := revenueMap[date]; exists {
				report.Revenue += order.ActualCost
				report.Orders++
			} else {
				revenueMap[date] = &model.RevenueReport{
					Date:    date,
					Revenue: order.ActualCost,
					Orders:  1,
				}
			}
		}
	}

	// Convert map to slice
	var revenueReport []model.RevenueReport
	for _, report := range revenueMap {
		revenueReport = append(revenueReport, *report)
	}

	return revenueReport, nil
}

// GetBranchStats retrieves statistics for a specific branch
func (s *DashboardService) GetBranchStats(ctx context.Context, branchID uuid.UUID) (*dashboardDTO.DashboardStats, error) {
	return s.GetDashboardStats(ctx, nil, &branchID)
}

// GetUserStats retrieves statistics for a specific user
func (s *DashboardService) GetUserStats(ctx context.Context, userID uuid.UUID) (*dashboardDTO.DashboardStats, error) {
	return s.GetDashboardStats(ctx, &userID, nil)
}
