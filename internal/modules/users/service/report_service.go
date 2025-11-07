package service

import (
	"context"
	"fmt"
	branchRepository "service/internal/modules/branches/repository"
	membershipRepository "service/internal/modules/membership/repository"
	orderRepository "service/internal/modules/orders/repository"
	paymentRepository "service/internal/modules/payments/repository"
	userRepository "service/internal/modules/users/repository"
	"service/internal/shared/model"
	"time"
)

// ReportService handles report business logic
type ReportService struct {
	orderRepo      *orderRepository.ServiceOrderRepository
	paymentRepo    *paymentRepository.PaymentRepository
	userRepo       *userRepository.UserRepository
	branchRepo     *branchRepository.BranchRepository
	membershipRepo *membershipRepository.MembershipRepository
}

// NewReportService creates a new report service
func NewReportService() *ReportService {
	return &ReportService{
		orderRepo:      orderRepository.NewServiceOrderRepository(),
		paymentRepo:    paymentRepository.NewPaymentRepository(),
		userRepo:       userRepository.NewUserRepository(),
		branchRepo:     branchRepository.NewBranchRepository(),
		membershipRepo: membershipRepository.NewMembershipRepository(),
	}
}

// MonthlyReport represents a monthly report
type MonthlyReport struct {
	Month           string                   `json:"month"`
	Year            int                      `json:"year"`
	TotalOrders     int64                    `json:"total_orders"`
	TotalRevenue    float64                  `json:"total_revenue"`
	TotalCustomers  int64                    `json:"total_customers"`
	NewCustomers    int64                    `json:"new_customers"`
	OrdersByStatus  map[string]int64         `json:"orders_by_status"`
	RevenueByBranch map[string]float64       `json:"revenue_by_branch"`
	OrdersByBranch  map[string]int64         `json:"orders_by_branch"`
	PaymentMethods  map[string]float64       `json:"payment_methods"`
	ServiceTypes    map[string]int64         `json:"service_types"`
	MembershipStats map[string]interface{}   `json:"membership_stats"`
	TopServices     []model.ServiceTypeStats `json:"top_services"`
	TopBranches     []model.BranchStats      `json:"top_branches"`
	GrowthMetrics   model.GrowthMetrics      `json:"growth_metrics"`
}

// GenerateMonthlyReport generates a comprehensive monthly report
func (s *ReportService) GenerateMonthlyReport(ctx context.Context, year int, month int) (*MonthlyReport, error) {
	// Validate month
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month: %d", month)
	}

	// Create date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// Get basic statistics (with fallback to 0 if error)
	totalOrders, err := s.orderRepo.CountOrdersByDateRange(ctx, startDate, endDate)
	if err != nil {
		totalOrders = 0 // Default to 0
	}

	totalRevenue, err := s.paymentRepo.GetTotalRevenueByDateRange(ctx, startDate, endDate)
	if err != nil {
		totalRevenue = 0 // Default to 0
	}

	totalCustomers, err := s.userRepo.CountCustomersByDateRange(ctx, startDate, endDate)
	if err != nil {
		totalCustomers = 0 // Default to 0
	}

	// Get new customers (registered in this month)
	newCustomers, err := s.userRepo.CountNewCustomersByDateRange(ctx, startDate, endDate)
	if err != nil {
		newCustomers = 0 // Default to 0
	}

	// Get orders by status (with fallback to empty map)
	ordersByStatus, err := s.orderRepo.GetOrdersByStatusInDateRange(ctx, startDate, endDate)
	if err != nil {
		ordersByStatus = make(map[string]int64)
	}

	// Get revenue by branch (with fallback to empty map)
	revenueByBranch, err := s.paymentRepo.GetRevenueByBranchInDateRange(ctx, startDate, endDate)
	if err != nil {
		revenueByBranch = make(map[string]float64)
	}

	// Get orders by branch (with fallback to empty map)
	ordersByBranch, err := s.orderRepo.GetOrdersByBranchInDateRange(ctx, startDate, endDate)
	if err != nil {
		ordersByBranch = make(map[string]int64)
	}

	// Get payment methods (with fallback to empty map)
	paymentMethods, err := s.paymentRepo.GetRevenueByPaymentMethodInDateRange(ctx, startDate, endDate)
	if err != nil {
		paymentMethods = make(map[string]float64)
	}

	// Get service types (with fallback to empty map)
	serviceTypes, err := s.orderRepo.GetOrdersByServiceTypeInDateRange(ctx, startDate, endDate)
	if err != nil {
		serviceTypes = make(map[string]int64)
	}

	// Get membership statistics (with fallback to empty map)
	membershipStats, err := s.membershipRepo.GetMembershipStats(ctx)
	if err != nil {
		membershipStats = make(map[string]interface{})
	}

	// Get top services (with fallback to empty slice)
	topServices, err := s.orderRepo.GetTopServiceTypesInDateRange(ctx, startDate, endDate, 10)
	if err != nil {
		topServices = []model.ServiceTypeStats{}
	}

	// Get top branches (with fallback to empty slice)
	topBranches, err := s.branchRepo.GetTopBranchesByRevenueInDateRange(ctx, startDate, endDate, 10)
	if err != nil {
		topBranches = []model.BranchStats{}
	}

	// Calculate growth metrics (with fallback to default)
	growthMetrics, err := s.calculateGrowthMetrics(ctx, year, month)
	if err != nil {
		growthMetrics = model.GrowthMetrics{
			OrderGrowth:     0,
			RevenueGrowth:   0,
			PreviousOrders:  0,
			PreviousRevenue: 0,
		}
	}

	report := &MonthlyReport{
		Month:           startDate.Format("January"),
		Year:            year,
		TotalOrders:     totalOrders,
		TotalRevenue:    totalRevenue,
		TotalCustomers:  totalCustomers,
		NewCustomers:    newCustomers,
		OrdersByStatus:  ordersByStatus,
		RevenueByBranch: revenueByBranch,
		OrdersByBranch:  ordersByBranch,
		PaymentMethods:  paymentMethods,
		ServiceTypes:    serviceTypes,
		MembershipStats: membershipStats,
		TopServices:     topServices,
		TopBranches:     topBranches,
		GrowthMetrics:   growthMetrics,
	}

	return report, nil
}

// calculateGrowthMetrics calculates growth metrics compared to previous month
func (s *ReportService) calculateGrowthMetrics(ctx context.Context, year int, month int) (model.GrowthMetrics, error) {
	// Get current month data
	currentStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	currentEnd := currentStart.AddDate(0, 1, 0).Add(-time.Second)

	currentOrders, err := s.orderRepo.CountOrdersByDateRange(ctx, currentStart, currentEnd)
	if err != nil {
		return model.GrowthMetrics{}, err
	}

	currentRevenue, err := s.paymentRepo.GetTotalRevenueByDateRange(ctx, currentStart, currentEnd)
	if err != nil {
		return model.GrowthMetrics{}, err
	}

	// Get previous month data
	prevStart := currentStart.AddDate(0, -1, 0)
	prevEnd := currentStart.Add(-time.Second)

	prevOrders, err := s.orderRepo.CountOrdersByDateRange(ctx, prevStart, prevEnd)
	if err != nil {
		return model.GrowthMetrics{}, err
	}

	prevRevenue, err := s.paymentRepo.GetTotalRevenueByDateRange(ctx, prevStart, prevEnd)
	if err != nil {
		return model.GrowthMetrics{}, err
	}

	// Calculate growth percentages
	orderGrowth := float64(0)
	revenueGrowth := float64(0)

	if prevOrders > 0 {
		orderGrowth = float64(currentOrders-prevOrders) / float64(prevOrders) * 100
	}

	if prevRevenue > 0 {
		revenueGrowth = (currentRevenue - prevRevenue) / prevRevenue * 100
	}

	return model.GrowthMetrics{
		OrderGrowth:     orderGrowth,
		RevenueGrowth:   revenueGrowth,
		PreviousOrders:  prevOrders,
		PreviousRevenue: prevRevenue,
	}, nil
}

// GetYearlyReport generates a yearly report
func (s *ReportService) GetYearlyReport(ctx context.Context, year int) (*model.YearlyReport, error) {
	// Get monthly data for the year
	var monthlyData []model.MonthlyReportData

	for month := 1; month <= 12; month++ {
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

		orders, err := s.orderRepo.CountOrdersByDateRange(ctx, startDate, endDate)
		if err != nil {
			continue
		}

		revenue, err := s.paymentRepo.GetTotalRevenueByDateRange(ctx, startDate, endDate)
		if err != nil {
			continue
		}

		monthlyData = append(monthlyData, model.MonthlyReportData{
			Month:   startDate.Format("January"),
			Orders:  orders,
			Revenue: revenue,
		})
	}

	// Calculate yearly totals
	var totalOrders int64
	var totalRevenue float64
	for _, data := range monthlyData {
		totalOrders += data.Orders
		totalRevenue += data.Revenue
	}

	return &model.YearlyReport{
		Year:         year,
		TotalOrders:  totalOrders,
		TotalRevenue: totalRevenue,
		MonthlyData:  monthlyData,
	}, nil
}
