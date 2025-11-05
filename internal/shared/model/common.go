package model

import (
	"errors"
	"math"
	// userDTO "service/internal/users/dto"
	"time"
)

// APIResponse represents the standard API response format
type APIResponse struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorResponse represents the error response format
type ErrorResponse struct {
	Status    string      `json:"status"`
	Error     string      `json:"error"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// PaginationResponse represents pagination information
type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Status     string             `json:"status"`
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
	Message    string             `json:"message"`
	Timestamp  time.Time          `json:"timestamp"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
	ExpiresIn    int64        `json:"expires_in"`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordRequest represents the change password request payload
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// ForgotPasswordRequest represents the forgot password request payload
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents the reset password request payload
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// ServiceStats represents service statistics
type ServiceStats struct {
	ServiceType string  `json:"service_type"`
	Count       int64   `json:"count"`
	Revenue     float64 `json:"revenue"`
}

// RevenueReport represents revenue report data
type RevenueReport struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

// Common errors
var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrNotFound          = errors.New("not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInternalError     = errors.New("internal server error")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrOrderNotFound     = errors.New("order not found")
	ErrBranchNotFound    = errors.New("branch not found")
	ErrPaymentNotFound   = errors.New("payment not found")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrEmailExists       = errors.New("email already exists")
	ErrPhoneExists       = errors.New("phone number already exists")
	ErrOrderNumberExists = errors.New("order number already exists")
	ErrInvoiceExists     = errors.New("invoice number already exists")
)

// SuccessResponse creates a success response
func SuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Status:    "success",
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// ErrorResponse creates an error response
func CreateErrorResponse(err string, message string, details interface{}) ErrorResponse {
	return ErrorResponse{
		Status:    "error",
		Error:     err,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// PaginatedSuccessResponse creates a paginated success response
func PaginatedSuccessResponse(data interface{}, pagination PaginationResponse, message string) PaginatedResponse {
	return PaginatedResponse{
		Status:     "success",
		Data:       data,
		Pagination: pagination,
		Message:    message,
		Timestamp:  time.Now(),
	}
}

// GetCurrentTimestamp returns current timestamp
func GetCurrentTimestamp() time.Time {
	return time.Now()
}

// CalculateDistance calculates distance between two coordinates using Haversine formula
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
