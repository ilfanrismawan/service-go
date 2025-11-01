package core

import (
	"time"
)

// BranchKPI represents KPI metrics for a branch
type BranchKPI struct {
	BranchID          string    `json:"branch_id"`
	BranchName        string    `json:"branch_name"`
	Revenue           float64   `json:"revenue"`
	TotalOrders       int64     `json:"total_orders"`
	CompletedOrders   int64     `json:"completed_orders"`
	PendingOrders     int64     `json:"pending_orders"`
	AverageHandlingTime float64 `json:"average_handling_time"` // in hours
	CustomerSatisfaction float64 `json:"customer_satisfaction"` // average rating 1-5
	Period            string    `json:"period"` // "daily", "weekly", "monthly"
	Date              time.Time `json:"date"`
}

// TechnicianPerformance represents performance metrics for a technician
type TechnicianPerformance struct {
	TechnicianID      string    `json:"technician_id"`
	TechnicianName    string    `json:"technician_name"`
	CompletedOrders   int64     `json:"completed_orders"`
	AverageRating     float64   `json:"average_rating"`
	AverageTime       float64   `json:"average_time"` // in hours
	OnTimeCompletion  float64   `json:"on_time_completion"` // percentage
	Period            string    `json:"period"`
	Date              time.Time `json:"date"`
}

// DashboardStats represents comprehensive dashboard statistics
type DashboardStats struct {
	// Overall stats
	TotalRevenue      float64   `json:"total_revenue"`
	TotalOrders       int64     `json:"total_orders"`
	ActiveCustomers   int64     `json:"active_customers"`
	ActiveBranches    int64     `json:"active_branches"`
	
	// Period comparisons
	RevenueGrowth     float64   `json:"revenue_growth"` // percentage
	OrdersGrowth      float64   `json:"orders_growth"` // percentage
	
	// Branch KPIs
	BranchKPIs        []BranchKPI `json:"branch_kpis"`
	
	// Technician performance
	TopTechnicians    []TechnicianPerformance `json:"top_technicians"`
	
	// Service type distribution
	ServiceTypeStats  []ServiceTypeStat `json:"service_type_stats"`
	
	// Payment method distribution
	PaymentMethodStats []PaymentMethodStat `json:"payment_method_stats"`
}

// ServiceTypeStat represents statistics for service types
type ServiceTypeStat struct {
	ServiceType string  `json:"service_type"`
	Count       int64   `json:"count"`
	Revenue     float64 `json:"revenue"`
	Percentage  float64 `json:"percentage"`
}

// PaymentMethodStat represents statistics for payment methods
type PaymentMethodStat struct {
	PaymentMethod string  `json:"payment_method"`
	Count         int64   `json:"count"`
	Revenue       float64 `json:"revenue"`
	Percentage    float64 `json:"percentage"`
}

