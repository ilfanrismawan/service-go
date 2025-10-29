package service

import (
	"context"
	"errors"
	"service/internal/core"
	"service/internal/repository"
	"service/internal/utils"

	"github.com/google/uuid"
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	orderRepo   *repository.ServiceOrderRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService() *PaymentService {
	return &PaymentService{
		paymentRepo: repository.NewPaymentRepository(),
		orderRepo:   repository.NewServiceOrderRepository(),
	}
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, req *core.PaymentRequest) (*core.PaymentResponse, error) {
	// Validate order exists
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	_, err = s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, core.ErrOrderNotFound
	}

	// Generate unique invoice number
	invoiceNumber := utils.GenerateInvoiceNumber()

	// Check if invoice number already exists (very unlikely but safety check)
	exists, err := s.paymentRepo.CheckInvoiceExists(ctx, invoiceNumber, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		// Regenerate if exists
		invoiceNumber = utils.GenerateInvoiceNumber()
	}

	// Create payment entity
	payment := &core.Payment{
		OrderID:       orderID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        core.PaymentStatusPending,
		InvoiceNumber: invoiceNumber,
		Notes:         req.Notes,
	}

	// Save to database
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Return response with populated data
	response := payment.ToResponse()
	return &response, nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*core.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrPaymentNotFound
	}

	response := payment.ToResponse()
	return &response, nil
}

// GetPaymentByInvoice retrieves a payment by invoice number
func (s *PaymentService) GetPaymentByInvoice(ctx context.Context, invoiceNumber string) (*core.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetByInvoiceNumber(ctx, invoiceNumber)
	if err != nil {
		return nil, core.ErrPaymentNotFound
	}

	response := payment.ToResponse()
	return &response, nil
}

// UpdatePaymentStatus updates the status of a payment
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status core.PaymentStatus, transactionID string) (*core.PaymentResponse, error) {
	// Get existing payment
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrPaymentNotFound
	}

	// Update status
	payment.Status = status
	if transactionID != "" {
		payment.TransactionID = transactionID
	}

	// Set paid_at if status is paid
	if status == core.PaymentStatusPaid {
		now := core.GetCurrentTimestamp()
		payment.PaidAt = &now
	}

	// Save changes
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	response := payment.ToResponse()
	return &response, nil
}

// ProcessMidtransPayment processes payment through Midtrans
func (s *PaymentService) ProcessMidtransPayment(ctx context.Context, req *core.MidtransPaymentRequest) (*core.MidtransPaymentResponse, error) {
	// Validate order exists
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	_, err = s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, core.ErrOrderNotFound
	}

	// Generate invoice number
	invoiceNumber := utils.GenerateInvoiceNumber()

	// Create payment record
	payment := &core.Payment{
		OrderID:       orderID,
		Amount:        req.Amount,
		PaymentMethod: core.PaymentMethodMidtrans,
		Status:        core.PaymentStatusPending,
		InvoiceNumber: invoiceNumber,
	}

	// Save payment
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// TODO: Integrate with actual Midtrans API
	// For now, return mock response
	response := &core.MidtransPaymentResponse{
		Token:         "mock-token-" + payment.ID.String(),
		RedirectURL:   "https://app.midtrans.com/snap/v2/vtweb/" + payment.ID.String(),
		StatusCode:    "201",
		StatusMessage: "Success, transaction is created",
	}

	return response, nil
}

// ListPayments retrieves payments with pagination and filters
func (s *PaymentService) ListPayments(ctx context.Context, page, limit int, filters *PaymentFilters) (*core.PaginatedResponse, error) {
	offset := (page - 1) * limit

	// Convert service-level filters to repository-level filters
	var repoFilters *repository.PaymentFilters
	if filters != nil {
		rf := repository.PaymentFilters{
			OrderID:       filters.OrderID,
			Status:        filters.Status,
			PaymentMethod: filters.PaymentMethod,
			DateFrom:      filters.DateFrom,
			DateTo:        filters.DateTo,
		}
		repoFilters = &rf
	}

	payments, total, err := s.paymentRepo.List(ctx, offset, limit, repoFilters)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var responses []core.PaymentResponse
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	// Calculate pagination
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	pagination := core.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return &core.PaginatedResponse{
		Status:     "success",
		Data:       responses,
		Pagination: pagination,
		Message:    "Payments retrieved successfully",
		Timestamp:  core.GetCurrentTimestamp(),
	}, nil
}

// GetPaymentsByOrder retrieves payments for a specific order
func (s *PaymentService) GetPaymentsByOrder(ctx context.Context, orderID uuid.UUID) ([]core.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var responses []core.PaymentResponse
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	return responses, nil
}

// GetPaymentsByStatus retrieves payments by status
func (s *PaymentService) GetPaymentsByStatus(ctx context.Context, status core.PaymentStatus) ([]core.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	var responses []core.PaymentResponse
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	return responses, nil
}

// PaymentFilters represents filters for payment queries
type PaymentFilters struct {
	OrderID       *uuid.UUID
	Status        *core.PaymentStatus
	PaymentMethod *core.PaymentMethod
	DateFrom      *string
	DateTo        *string
}
