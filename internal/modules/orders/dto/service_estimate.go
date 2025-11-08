package dto

import (
	orderEntity "service-go/internal/modules/orders/entity"
)

// ServiceEstimate represents service cost and time estimation
type ServiceEstimate struct {
	MinPrice      int64 `json:"min_price"`      // Harga minimum
	MaxPrice      int64 `json:"max_price"`      // Harga maksimum
	EstimatedDays int   `json:"estimated_days"` // Estimasi hari pengerjaan
	Warranty      int   `json:"warranty"`       // Garansi dalam hari
}

// GetServiceEstimate returns estimated cost and time for a service type
func GetServiceEstimate(serviceType orderEntity.ServiceType) ServiceEstimate {
	estimates := map[orderEntity.ServiceType]ServiceEstimate{
		orderEntity.ServiceTypeScreenRepair: {
			MinPrice:      500000,
			MaxPrice:      2000000,
			EstimatedDays: 2,
			Warranty:      30,
		},
		orderEntity.ServiceTypeBatteryReplacement: {
			MinPrice:      300000,
			MaxPrice:      800000,
			EstimatedDays: 1,
			Warranty:      90,
		},
		orderEntity.ServiceTypeWaterDamage: {
			MinPrice:      800000,
			MaxPrice:      3000000,
			EstimatedDays: 3,
			Warranty:      30,
		},
		orderEntity.ServiceTypeSoftwareIssue: {
			MinPrice:      150000,
			MaxPrice:      500000,
			EstimatedDays: 1,
			Warranty:      14,
		},
		orderEntity.ServiceTypeHardwareRepair: {
			MinPrice:      1000000,
			MaxPrice:      5000000,
			EstimatedDays: 5,
			Warranty:      60,
		},
		orderEntity.ServiceTypeOther: {
			MinPrice:      200000,
			MaxPrice:      1000000,
			EstimatedDays: 2,
			Warranty:      30,
		},
	}

	if estimate, ok := estimates[serviceType]; ok {
		return estimate
	}
	// Default estimate
	return ServiceEstimate{
		MinPrice:      200000,
		MaxPrice:      1000000,
		EstimatedDays: 2,
		Warranty:      30,
	}
}

