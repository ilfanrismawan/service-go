package service

import (
	"context"
	"errors"
	orderRepo "service/internal/modules/orders/repository"
	serviceRepo "service/internal/modules/services/repository"
	trackingRepo "service/internal/modules/tracking/repository"
	"service/internal/shared/model"
	"time"

	"github.com/google/uuid"
)

// LocationTrackingService handles location tracking business logic
type LocationTrackingService struct {
	trackingRepo     *trackingRepo.LocationTrackingRepository
	currentLocRepo   *trackingRepo.CurrentLocationRepository
	orderRepo        *orderRepo.ServiceOrderRepository
	providerRepo     *serviceRepo.ServiceProviderRepository
}

// NewLocationTrackingService creates a new location tracking service
func NewLocationTrackingService() *LocationTrackingService {
	return &LocationTrackingService{
		trackingRepo:   trackingRepo.NewLocationTrackingRepository(),
		currentLocRepo: trackingRepo.NewCurrentLocationRepository(),
		orderRepo:      orderRepo.NewServiceOrderRepository(),
		providerRepo:   serviceRepo.NewServiceProviderRepository(),
	}
}

// UpdateLocation updates location for an order (called by courier/provider)
func (s *LocationTrackingService) UpdateLocation(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, req *model.LocationUpdateRequest) (*model.LocationUpdateResponse, error) {
	// Validate order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	// Validate user is courier or provider for this order
	isAuthorized := false
	
	// Check if user is courier
	if order.CourierID != nil && *order.CourierID == userID {
		isAuthorized = true
	}
	
	// Check if user is provider
	if !isAuthorized && order.ServiceProviderID != nil {
		provider, err := s.providerRepo.GetByID(ctx, *order.ServiceProviderID)
		if err == nil && provider.UserID == userID {
			isAuthorized = true
		}
	}
	
	if !isAuthorized {
		return nil, errors.New("user is not authorized to update location for this order")
	}

	// Calculate ETA if destination is known
	eta := 0
	distance := 0.0
	if order.PickupLatitude != nil && order.PickupLongitude != nil {
		distance = model.CalculateDistance(
			req.Latitude,
			req.Longitude,
			*order.PickupLatitude,
			*order.PickupLongitude,
		)
		
		// Calculate ETA based on distance and speed
		if req.Speed > 0 {
			eta = int((distance / req.Speed) * 60) // Convert to minutes
		} else {
			// Default speed: 30 km/h for city driving
			eta = int((distance / 30) * 60)
		}
	}

	// Create location tracking history
	tracking := &model.LocationTracking{
		OrderID:   orderID,
		UserID:    userID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Accuracy:  req.Accuracy,
		Speed:     req.Speed,
		Heading:   req.Heading,
		Timestamp: time.Now(),
	}

	if err := s.trackingRepo.CreateLocationHistory(ctx, tracking); err != nil {
		return nil, err
	}

	// Update current location
	currentLocation := &model.CurrentLocation{
		OrderID:   orderID,
		UserID:    userID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Accuracy:  req.Accuracy,
		Speed:     req.Speed,
		Heading:   req.Heading,
		ETA:       eta,
		Distance:  distance,
		UpdatedAt: time.Now(),
	}

	if err := s.currentLocRepo.UpsertCurrentLocation(ctx, currentLocation); err != nil {
		return nil, err
	}

	// Update order with current location and ETA
	order.CurrentLatitude = &req.Latitude
	order.CurrentLongitude = &req.Longitude
	order.ETA = &eta
	now := time.Now()
	order.LastLocationUpdate = &now

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	response := currentLocation.ToResponse()
	return &response, nil
}

// GetCurrentLocation retrieves current location for an order (called by customer)
func (s *LocationTrackingService) GetCurrentLocation(ctx context.Context, orderID uuid.UUID) (*model.LocationUpdateResponse, error) {
	// Validate order exists
	_, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	location, err := s.currentLocRepo.GetCurrentLocation(ctx, orderID)
	if err != nil {
		return nil, errors.New("location not found")
	}

	response := location.ToResponse()
	return &response, nil
}

// GetLocationHistory retrieves location history for an order
func (s *LocationTrackingService) GetLocationHistory(ctx context.Context, orderID uuid.UUID, limit int) ([]model.LocationHistoryResponse, error) {
	// Validate order exists
	_, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	history, err := s.trackingRepo.GetLocationHistory(ctx, orderID, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]model.LocationHistoryResponse, len(history))
	for i, h := range history {
		responses[i] = h.ToHistoryResponse()
	}

	return responses, nil
}

// CalculateETA calculates ETA based on current location and destination
func (s *LocationTrackingService) CalculateETA(ctx context.Context, orderID uuid.UUID, req *model.ETACalculationRequest) (*model.ETAResponse, error) {
	// Validate order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	// Get destination from order
	var destLat, destLon float64
	if order.PickupLatitude != nil && order.PickupLongitude != nil {
		destLat = *order.PickupLatitude
		destLon = *order.PickupLongitude
	} else if order.ServiceProviderID != nil {
		// Get provider location
		provider, err := s.providerRepo.GetByID(ctx, *order.ServiceProviderID)
		if err != nil {
			return nil, errors.New("destination location not found")
		}
		destLat = provider.Latitude
		destLon = provider.Longitude
	} else if order.BranchID != nil {
		// Get branch location (for legacy orders)
		// Note: Need to import branch repository if needed
		return nil, errors.New("destination location not found - branch location not implemented")
	} else {
		return nil, errors.New("destination location not found")
	}

	// Calculate distance
	distance := model.CalculateDistance(
		req.CurrentLatitude,
		req.CurrentLongitude,
		destLat,
		destLon,
	)

	// Calculate ETA
	speed := req.Speed
	if speed == 0 {
		speed = 30 // Default: 30 km/h for city driving
	}

	etaMinutes := int((distance / speed) * 60)
	estimatedArrival := time.Now().Add(time.Duration(etaMinutes) * time.Minute)

	return &model.ETAResponse{
		ETA:              etaMinutes,
		Distance:         distance,
		EstimatedArrival: estimatedArrival,
	}, nil
}

