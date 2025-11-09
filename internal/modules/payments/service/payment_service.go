package service

import (
	"context"
	"errors"
	"fmt"
	"service/internal/modules/orders/repository"
	pay "service/internal/modules/payments/legacy_payment"
	repo "service/internal/modules/payments/repository"
	"service/internal/shared/model"
	"service/internal/shared/utils"
	"time"

	"github.com/google/uuid"
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo    *repo.PaymentRepository
	orderRepo      *repository.ServiceOrderRepository
	midtransService *pay.MidtransService
}

// NewPaymentService creates a new payment service
func NewPaymentService() *PaymentService {
	return &PaymentService{
		paymentRepo:     repo.NewPaymentRepository(),
		orderRepo:       repository.NewServiceOrderRepository(),
		midtransService: pay.NewMidtransService(),
	}
}

// HandleMidtransCallback verifies payload and updates payment/order state
func (s *PaymentService) HandleMidtransCallback(ctx context.Context, cb *model.MidtransCallbackPayload, serverKey string) error {
	// Verify signature: sha512(order_id+status_code+gross_amount+server_key)
	// Note: Midtrans signature format may vary, adjust if needed
	if serverKey == "" {
		return errors.New("server key is required for signature verification")
	}
	expected := utils.SHA512Hex(cb.OrderID + cb.StatusCode + cb.GrossAmount + serverKey)
	if expected != cb.SignatureKey {
		return errors.New("invalid signature - callback may be from unauthorized source")
	}

	// Find payment by transaction ID first, fallback to invoice/order_id equals payment.ID if used as order_id
	var payment *model.Payment
	// Try by transaction ID
	if cb.TransactionID != "" {
		if p, err := s.paymentRepo.GetByTransactionID(ctx, cb.TransactionID); err == nil {
			payment = p
		}
	}
	if payment == nil {
		// Try by invoice number or payment ID from order_id field
		if p, err := s.paymentRepo.GetByInvoiceNumber(ctx, cb.OrderID); err == nil {
			payment = p
		}
	}
	if payment == nil {
		return model.ErrPaymentNotFound
	}

	// Idempotency: if status already mapped, no-op
	mapped := model.PaymentStatusPending
	switch cb.TransactionStatus {
	case "capture", "settlement":
		mapped = model.PaymentStatusPaid
	case "pending":
		mapped = model.PaymentStatusPending
	case "deny", "expire", "cancel":
		mapped = model.PaymentStatusFailed
	case "refund", "partial_refund":
		mapped = model.PaymentStatusRefunded
	default:
		mapped = model.PaymentStatusPending
	}
	if payment.Status == mapped {
		return nil
	}

	// Update payment
	payment.Status = mapped
	if cb.TransactionID != "" {
		payment.TransactionID = cb.TransactionID
	}
	if mapped == model.PaymentStatusPaid {
		now := model.GetCurrentTimestamp()
		payment.PaidAt = &now
	}
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}

	// Update order on success
	if mapped == model.PaymentStatusPaid {
		if order, err := s.orderRepo.GetByID(ctx, payment.OrderID); err == nil {
			order.Status = model.StatusReady
			_ = s.orderRepo.Update(ctx, order)
		}
	}
	return nil
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, req *model.PaymentRequest) (*model.PaymentResponse, error) {
	// Validate order exists
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
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
	payment := &model.Payment{
		OrderID:       orderID,
		UserID:        order.CustomerID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        model.PaymentStatusPending,
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
func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*model.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrPaymentNotFound
	}

	response := payment.ToResponse()
	return &response, nil
}

// GetPaymentByInvoice retrieves a payment by invoice number
func (s *PaymentService) GetPaymentByInvoice(ctx context.Context, invoiceNumber string) (*model.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetByInvoiceNumber(ctx, invoiceNumber)
	if err != nil {
		return nil, model.ErrPaymentNotFound
	}

	response := payment.ToResponse()
	return &response, nil
}

// UpdatePaymentStatus updates the status of a payment
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status model.PaymentStatus, transactionID string) (*model.PaymentResponse, error) {
	// Get existing payment
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrPaymentNotFound
	}

	// Update status
	payment.Status = status
	if transactionID != "" {
		payment.TransactionID = transactionID
	}

	// Set paid_at if status is paid
	if status == model.PaymentStatusPaid {
		now := model.GetCurrentTimestamp()
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
func (s *PaymentService) ProcessMidtransPayment(ctx context.Context, req *model.MidtransPaymentRequest) (*model.MidtransPaymentResponse, error) {
	// Validate order exists
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Generate invoice number
	invoiceNumber := utils.GenerateInvoiceNumber()

	// Check if invoice number already exists
	exists, err := s.paymentRepo.CheckInvoiceExists(ctx, invoiceNumber, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		// Regenerate if exists
		invoiceNumber = utils.GenerateInvoiceNumber()
	}

	// Create payment record
	payment := &model.Payment{
		OrderID:       orderID,
		UserID:        order.CustomerID,
		Amount:        req.Amount,
		PaymentMethod: model.PaymentMethodMidtrans,
		Status:        model.PaymentStatusPending,
		InvoiceNumber: invoiceNumber,
	}

	// Save payment
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Create Midtrans payment request
	midtransReq := &pay.MidtransPaymentRequest{
		TransactionDetails: pay.TransactionDetails{
			OrderID:     payment.ID.String(), // Use payment ID as order_id for Midtrans
			GrossAmount: int64(req.Amount * 100), // Convert to cents
		},
		PaymentType: "credit_card", // Default to credit_card, can be extended
		CustomExpiry: &pay.CustomExpiry{
			OrderTime:      time.Now().Format("2006-01-02 15:04:05"),
			ExpiryDuration: 24,
			Unit:           "hour",
		},
	}

	// Call Midtrans API
	midtransResp, err := s.midtransService.CreatePayment(ctx, midtransReq)
	if err != nil {
		// Update payment status to failed
		payment.Status = model.PaymentStatusFailed
		s.paymentRepo.Update(ctx, payment)
		return nil, fmt.Errorf("failed to create Midtrans payment: %w", err)
	}

	// Update payment with transaction ID
	payment.TransactionID = midtransResp.TransactionID
	if midtransResp.RedirectURL != "" {
		payment.InvoiceURL = midtransResp.RedirectURL
	}

	// Update payment
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	// Map Midtrans response to model response
	response := &model.MidtransPaymentResponse{
		Token:         midtransResp.TransactionID, // Use transaction ID as token
		RedirectURL:   midtransResp.RedirectURL,
		StatusCode:    midtransResp.StatusCode,
		StatusMessage: midtransResp.StatusMessage,
	}

	return response, nil
}

// ReconcilePendingPayments checks pending Midtrans payments and updates their status
func (s *PaymentService) ReconcilePendingPayments(ctx context.Context) error {
	// Get pending payments
	pendingList, err := s.paymentRepo.GetByStatus(ctx, model.PaymentStatusPending)
	if err != nil {
		return err
	}
	if len(pendingList) == 0 {
		return nil
	}
	ms := pay.NewMidtransService()
	for _, p := range pendingList {
		// Only reconcile online payments
		if p.PaymentMethod != model.PaymentMethodMidtrans && p.PaymentMethod != model.PaymentMethodGopay && p.PaymentMethod != model.PaymentMethodQris && p.PaymentMethod != model.PaymentMethodBankTransfer && p.PaymentMethod != model.PaymentMethodMandiriEchannel {
			continue
		}
		if p.TransactionID == "" {
			continue
		}
		resp, err := ms.GetPaymentStatus(ctx, p.TransactionID)
		if err != nil {
			continue
		}
		// Map Midtrans status
		newStatus := model.PaymentStatusPending
		switch resp.TransactionStatus {
		case "capture", "settlement":
			newStatus = model.PaymentStatusPaid
		case "pending":
			newStatus = model.PaymentStatusPending
		case "deny", "expire", "cancel":
			newStatus = model.PaymentStatusFailed
		case "refund", "partial_refund":
			newStatus = model.PaymentStatusRefunded
		}
		if p.Status == newStatus {
			continue
		}
		p.Status = newStatus
		if newStatus == model.PaymentStatusPaid {
			now := model.GetCurrentTimestamp()
			p.PaidAt = &now
		}
		_ = s.paymentRepo.Update(ctx, p)
		if newStatus == model.PaymentStatusPaid {
			if order, err := s.orderRepo.GetByID(ctx, p.OrderID); err == nil {
				order.Status = model.StatusReady
				_ = s.orderRepo.Update(ctx, order)
			}
		}
	}
	return nil
}

// ListPayments retrieves payments with pagination and filters
func (s *PaymentService) ListPayments(ctx context.Context, page, limit int, filters *PaymentFilters) (*model.PaginatedResponse, error) {
	offset := (page - 1) * limit

	// Convert service-level filters to repository-level filters
	var repoFilters *repo.PaymentFilters
	if filters != nil {
		rf := repo.PaymentFilters{
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
	var responses []model.PaymentResponse
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	// Calculate pagination
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	pagination := model.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return &model.PaginatedResponse{
		Status:     "success",
		Data:       responses,
		Pagination: pagination,
		Message:    "Payments retrieved successfully",
		Timestamp:  model.GetCurrentTimestamp(),
	}, nil
}

// GetPaymentsByOrder retrieves payments for a specific order
func (s *PaymentService) GetPaymentsByOrder(ctx context.Context, orderID uuid.UUID) ([]model.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var responses []model.PaymentResponse
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	return responses, nil
}

// GetPaymentsByStatus retrieves payments by status
func (s *PaymentService) GetPaymentsByStatus(ctx context.Context, status model.PaymentStatus) ([]model.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	var responses []model.PaymentResponse
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	return responses, nil
}

// PaymentFilters represents filters for payment queries
type PaymentFilters struct {
	OrderID       *uuid.UUID
	Status        *model.PaymentStatus
	PaymentMethod *model.PaymentMethod
	DateFrom      *string
	DateTo        *string
}
