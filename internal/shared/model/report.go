package model

// ServiceTypeStats represents statistics for a service type
type ServiceTypeStats struct {
	ServiceType string  `json:"service_type"`
	Count       int64   `json:"count"`
	Revenue     float64 `json:"revenue"`
	Percentage  float64 `json:"percentage"`
}

// BranchStats represents statistics for a branch
type BranchStats struct {
	BranchID   string  `json:"branch_id"`
	BranchName string  `json:"branch_name"`
	Orders     int64   `json:"orders"`
	Revenue    float64 `json:"revenue"`
	Percentage float64 `json:"percentage"`
}

// GrowthMetrics represents growth metrics
type GrowthMetrics struct {
	OrderGrowth     float64 `json:"order_growth"`
	RevenueGrowth   float64 `json:"revenue_growth"`
	PreviousOrders  int64   `json:"previous_orders"`
	PreviousRevenue float64 `json:"previous_revenue"`
}

// MonthlyReportData represents monthly report data
type MonthlyReportData struct {
	Month   string  `json:"month"`
	Orders  int64   `json:"orders"`
	Revenue float64 `json:"revenue"`
}

// YearlyReport represents a yearly report
type YearlyReport struct {
	Year         int                 `json:"year"`
	TotalOrders  int64               `json:"total_orders"`
	TotalRevenue float64             `json:"total_revenue"`
	MonthlyData  []MonthlyReportData `json:"monthly_data"`
}
