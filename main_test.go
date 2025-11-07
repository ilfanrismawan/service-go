package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	branchHandler "service/internal/modules/branches/handler"
	orderHandler "service/internal/modules/orders/handler"
	paymentHandler "service/internal/modules/payments/handler"
	authHandler "service/internal/modules/users/handler"
	"service/internal/shared/handlers"
	"service/internal/shared/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.New()
	r.GET("/health", handlers.NewHealthHandler().HealthCheck)

	// Create a test request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "iPhone Service API is running", response["message"])
}

func TestAuthRegister(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.New()
	authHandler := authHandler.NewAuthHandler()
	r.POST("/auth/register", authHandler.Register)

	// Create test user data
	userData := model.UserRequest{
		FullName: "Test User",
		Email:    "test@example.com",
		Phone:    "081234567890",
		Password: "password123",
		Role:     model.RolePelanggan,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(userData)
	assert.NoError(t, err)

	// Create a test request
	req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.APIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "User registered successfully", response.Message)
}

func TestAuthLogin(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.New()
	authHandler := authHandler.NewAuthHandler()
	r.POST("/auth/login", authHandler.Login)

	// Create test login data
	loginData := model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(loginData)
	assert.NoError(t, err)

	// Create a test request
	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.APIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Login successful", response.Message)
}

func TestBranchList(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.New()
	branchHandler := branchHandler.NewBranchHandler()
	r.GET("/branches", branchHandler.GetBranches)

	// Create a test request
	req, err := http.NewRequest("GET", "/branches", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.PaginatedResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response.Status)
}

func TestOrderCreate(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.New()
	orderHandler := orderHandler.NewOrderHandler()
	r.POST("/orders", orderHandler.CreateOrder)

	// Create test order data
	orderData := model.ServiceOrderRequest{
		IPhoneType:        "iPhone 14 Pro",
		Complaint:         "Screen cracked",
		PickupLocation:    "Jakarta",
		BranchID:          "550e8400-e29b-41d4-a716-446655440000",
		EstimatedCost:     500000,
		EstimatedDuration: 3,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(orderData)
	assert.NoError(t, err)

	// Create a test request
	req, err := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.APIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Order created successfully", response.Message)
}

func TestPaymentCreateInvoice(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.New()
	paymentHandler := paymentHandler.NewPaymentHandler()
	r.POST("/payments/create-invoice", paymentHandler.CreateInvoice)

	// Create test payment data
	paymentData := model.PaymentRequest{
		OrderID:       "550e8400-e29b-41d4-a716-446655440000",
		Amount:        500000,
		PaymentMethod: model.PaymentMethodMidtrans,
		Notes:         "Payment for iPhone repair",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(paymentData)
	assert.NoError(t, err)

	// Create a test request
	req, err := http.NewRequest("POST", "/payments/create-invoice", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.APIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Invoice created successfully", response.Message)
}

// Benchmark tests
func BenchmarkHealthCheck(b *testing.B) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", handlers.NewHealthHandler().HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkAuthLogin(b *testing.B) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	authHandler := authHandler.NewAuthHandler()
	r.POST("/auth/login", authHandler.Login)

	loginData := model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}
