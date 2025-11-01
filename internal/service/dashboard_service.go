package service

import (
	"context"
	"fmt"
	"service/internal/core"
	"service/internal/repository"
	"time"

	"github.com/google/uuid"
)

// Helper function for bool pointer
func boolPtr(b bool) *bool {
	return &b
}

// DashboardService handles dashboard business logic
type DashboardService struct {
	orderRepo  *repository.ServiceOrderRepository
	userRepo   *repository.UserRepository
	branchRepo *repository.BranchRepository
	paymentRepo *repository.PaymentRepository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService() *DashboardService {
	return &DashboardService{
		orderRepo:  repository.NewServiceOrderRepository(),
		userRepo:   repository.NewUserRepository(),
		branchRepo: repository.NewBranchRepository(),
		paymentRepo: repository.NewPaymentRepository(),
	}
}

// GetDashboardStats retrieves dashboard statistics
func (s *DashboardService) GetDashboardStats(ctx context.Context, userID *uuid.UUID, branchID *uuid.UUID) (*core.DashboardStats, error) {
	// Get user role if userID is provided
	var userRole *core.UserRole
	if userID != nil {
		user, err := s.userRepo.GetByID(ctx, *userID)
		if err != nil {
			return nil, err
		}
		userRole = &user.Role
	}

	// Set filters based on user role
	filters := &repository.ServiceOrderFilters{}
	if branchID != nil {
		filters.BranchID = branchID
	} else if userRole != nil && *userRole != core.RoleAdminPusat {
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
	stats := &core.DashboardStats{
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
		case core.StatusPendingPickup, core.StatusOnPickup, core.StatusInService, core.StatusReady, core.StatusDelivered:
			stats.PendingOrders++
		case core.StatusCompleted:
			stats.CompletedOrders++
		}
	}

	// Get payment statistics
	paymentFilters := &repository.PaymentFilters{}
	if branchID != nil {
		paymentFilters.OrderID = branchID // This would need to be adjusted based on your payment structure
	}

	payments, _, err := s.paymentRepo.List(ctx, 0, 1000, paymentFilters)
	if err != nil {
		return nil, err
	}

	// Calculate total revenue
	for _, payment := range payments {
		if payment.Status == core.PaymentStatusPaid {
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
		role := core.RolePelanggan
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
func (s *DashboardService) GetServiceStats(ctx context.Context, userID *uuid.UUID, branchID *uuid.UUID) ([]core.ServiceStats, error) {
	// Set filters based on user role
	filters := &repository.ServiceOrderFilters{}
	if branchID != nil {
		filters.BranchID = branchID
	} else if userID != nil {
		user, err := s.userRepo.GetByID(ctx, *userID)
		if err == nil && user.Role != core.RoleAdminPusat && user.BranchID != nil {
			filters.BranchID = user.BranchID
		}
	}

	// Get all orders
	orders, _, err := s.orderRepo.List(ctx, 0, 1000, filters)
	if err != nil {
		return nil, err
	}

	// Group by service type
	serviceStatsMap := make(map[core.ServiceType]*core.ServiceStats)
	for _, order := range orders {
		if order.Status == core.StatusCompleted {
			if stats, exists := serviceStatsMap[order.ServiceType]; exists {
				stats.Count++
				stats.Revenue += order.ActualCost
			} else {
				serviceStatsMap[order.ServiceType] = &core.ServiceStats{
					ServiceType: string(order.ServiceType),
					Count:       1,
					Revenue:     order.ActualCost,
				}
			}
		}
	}

	// Convert map to slice
	var serviceStats []core.ServiceStats
	for _, stats := range serviceStatsMap {
		serviceStats = append(serviceStats, *stats)
	}

	return serviceStats, nil
}

// GetRevenueReport retrieves revenue report for a date range
func (s *DashboardService) GetRevenueReport(ctx context.Context, userID *uuid.UUID, branchID *uuid.UUID, dateFrom, dateTo string) ([]core.RevenueReport, error) {
	// Set filters based on user role
	filters := &repository.ServiceOrderFilters{}
	if branchID != nil {
		filters.BranchID = branchID
	} else if userID != nil {
		user, err := s.userRepo.GetByID(ctx, *userID)
		if err == nil && user.Role != core.RoleAdminPusat && user.BranchID != nil {
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
	revenueMap := make(map[string]*core.RevenueReport)
	for _, order := range orders {
		if order.Status == core.StatusCompleted {
			date := order.CreatedAt.Format("2006-01-02")
			if report, exists := revenueMap[date]; exists {
				report.Revenue += order.ActualCost
				report.Orders++
			} else {
				revenueMap[date] = &core.RevenueReport{
					Date:    date,
					Revenue: order.ActualCost,
					Orders:  1,
				}
			}
		}
	}

	// Convert map to slice
	var revenueReport []core.RevenueReport
	for _, report := range revenueMap {
		revenueReport = append(revenueReport, *report)
	}

	return revenueReport, nil
}

// GetBranchStats retrieves statistics for a specific branch
func (s *DashboardService) GetBranchStats(ctx context.Context, branchID uuid.UUID) (*core.DashboardStats, error) {
	return s.GetDashboardStats(ctx, nil, &branchID)
}

// GetUserStats retrieves statistics for a specific user
func (s *DashboardService) GetUserStats(ctx context.Context, userID uuid.UUID) (*core.DashboardStats, error) {
	return s.GetDashboardStats(ctx, &userID, nil)
}

// GetBranchKPIs retrieves KPI metrics for all branches or a specific branch
func (s *DashboardService) GetBranchKPIs(ctx context.Context, branchID *uuid.UUID, period string) ([]core.BranchKPI, error) {
	var branches []core.Branch
	var err error

	if branchID != nil {
		branch, err := s.branchRepo.GetByID(ctx, *branchID)
		if err != nil {
			return nil, err
		}
		branches = []core.Branch{*branch}
	} else {
		branches, err = s.branchRepo.GetActiveBranches(ctx)
		if err != nil {
			return nil, err
		}
	}

	var kpis []core.BranchKPI
	for _, branch := range branches {
		kpi, err := s.calculateBranchKPI(ctx, branch.ID, period)
		if err != nil {
			continue // Skip branches with errors
		}
		kpis = append(kpis, *kpi)
	}

	return kpis, nil
}

// calculateBranchKPI calculates KPI for a specific branch
func (s *DashboardService) calculateBranchKPI(ctx context.Context, branchID uuid.UUID, period string) (*core.BranchKPI, error) {
	branch, err := s.branchRepo.GetByID(ctx, branchID)
	if err != nil {
		return nil, err
	}

	// Set date range based on period
	now := time.Now()
	var startDate, endDate time.Time

	switch period {
	case "daily":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = now
	case "weekly":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}
		startDate = now.AddDate(0, 0, -weekday+1)
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
		endDate = now
	case "monthly":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = now
	default:
		period = "monthly"
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = now
	}

	// Get orders for this branch in period
	filters := &repository.ServiceOrderFilters{
		BranchID: &branchID,
	}
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")
	filters.DateFrom = &startDateStr
	filters.DateTo = &endDateStr

	orders, _, err := s.orderRepo.List(ctx, 0, 10000, filters)
	if err != nil {
		return nil, err
	}

	// Calculate metrics
	var revenue float64
	var totalOrders int64
	var completedOrders int64
	var pendingOrders int64
	var totalHandlingTime float64
	var completedCount int64

	for _, order := range orders {
		totalOrders++
		switch order.Status {
		case core.StatusCompleted:
			completedOrders++
			revenue += order.ActualCost
			// Calculate handling time
			if !order.CreatedAt.IsZero() && !order.UpdatedAt.IsZero() {
				duration := order.UpdatedAt.Sub(order.CreatedAt)
				totalHandlingTime += duration.Hours()
				completedCount++
			}
		case core.StatusPendingPickup, core.StatusOnPickup, core.StatusInService, core.StatusReady, core.StatusDelivered:
			pendingOrders++
		}
	}

	avgHandlingTime := float64(0)
	if completedCount > 0 {
		avgHandlingTime = totalHandlingTime / float64(completedCount)
	}

	// Get average customer satisfaction (rating)
	ratingRepo := repository.NewRatingRepository()
	ratingFilters := &repository.RatingFilters{
		BranchID: &branchID,
		IsPublic: boolPtr(true),
	}
	avgRating, err := ratingRepo.GetAverageRating(ctx, ratingFilters)
	if err == nil && avgRating != nil && avgRating.TotalRatings > 0 {
		customerSatisfaction = avgRating.AverageRating
	}

	return &core.BranchKPI{
		BranchID:            branch.ID.String(),
		BranchName:          branch.Name,
		Revenue:             revenue,
		TotalOrders:         totalOrders,
		CompletedOrders:     completedOrders,
		PendingOrders:       pendingOrders,
		AverageHandlingTime: avgHandlingTime,
		CustomerSatisfaction: customerSatisfaction,
		Period:             period,
		Date:               now,
	}, nil
}

// GetTechnicianPerformance retrieves performance metrics for technicians
func (s *DashboardService) GetTechnicianPerformance(ctx context.Context, technicianID *uuid.UUID, period string) ([]core.TechnicianPerformance, error) {
	// Get all technicians or specific one
	var technicians []core.User
	var err error

	if technicianID != nil {
		tech, err := s.userRepo.GetByID(ctx, *technicianID)
		if err != nil {
			return nil, err
		}
		if tech.Role != core.RoleTeknisi {
			return nil, fmt.Errorf("user is not a technician")
		}
		technicians = []core.User{*tech}
	} else {
		role := core.RoleTeknisi
		technicians, _, err = s.userRepo.List(ctx, 0, 1000, &role, nil)
		if err != nil {
			return nil, err
		}
	}

	var performances []core.TechnicianPerformance
	for _, tech := range technicians {
		perf, err := s.calculateTechnicianPerformance(ctx, tech.ID, period)
		if err != nil {
			continue
		}
		performances = append(performances, *perf)
	}

	return performances, nil
}

// calculateTechnicianPerformance calculates performance for a specific technician
func (s *DashboardService) calculateTechnicianPerformance(ctx context.Context, technicianID uuid.UUID, period string) (*core.TechnicianPerformance, error) {
	tech, err := s.userRepo.GetByID(ctx, technicianID)
	if err != nil {
		return nil, err
	}

	// Set date range
	now := time.Now()
	var startDate time.Time
	switch period {
	case "daily":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "weekly":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startDate = now.AddDate(0, 0, -weekday+1)
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	case "monthly":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	default:
		period = "monthly"
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	// Get orders assigned to this technician
	filters := &repository.ServiceOrderFilters{
		TechnicianID: &technicianID,
	}
	startDateStr := startDate.Format("2006-01-02")
	filters.DateFrom = &startDateStr

	orders, _, err := s.orderRepo.List(ctx, 0, 10000, filters)
	if err != nil {
		return nil, err
	}

	// Calculate metrics
	var completedOrders int64
	var totalTime float64
	var onTimeCount int64

	for _, order := range orders {
		if order.Status == core.StatusCompleted {
			completedOrders++
			if !order.CreatedAt.IsZero() && !order.UpdatedAt.IsZero() {
				duration := order.UpdatedAt.Sub(order.CreatedAt)
				totalTime += duration.Hours()
				// Simple on-time check: completed within 3 days
				if duration.Hours() <= 72 {
					onTimeCount++
				}
			}
		}
	}

	avgTime := float64(0)
	if completedOrders > 0 {
		avgTime = totalTime / float64(completedOrders)
	}

	onTimePercentage := float64(0)
	if completedOrders > 0 {
		onTimePercentage = (float64(onTimeCount) / float64(completedOrders)) * 100
	}

	// Get average rating from rating repository
	ratingRepo := repository.NewRatingRepository()
	ratingFilters := &repository.RatingFilters{
		TechnicianID: &technicianID,
		IsPublic:     boolPtr(true),
	}
	avgRatingData, err := ratingRepo.GetAverageRating(ctx, ratingFilters)
	if err == nil && avgRatingData != nil && avgRatingData.TotalRatings > 0 {
		avgRating = avgRatingData.AverageRating
	}

	return &core.TechnicianPerformance{
		TechnicianID:     tech.ID.String(),
		TechnicianName:   tech.FullName,
		CompletedOrders:  completedOrders,
		AverageRating:    avgRating,
		AverageTime:      avgTime,
		OnTimeCompletion: onTimePercentage,
		Period:           period,
		Date:             now,
	}, nil
}