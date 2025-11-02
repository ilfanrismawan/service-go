package service

import (
	"context"
	"errors"
<<<<<<< HEAD
=======
<<<<<<<< HEAD:internal/service/payment_service.go
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	"service/internal/core"
    pay "service/internal/payments"
	"service/internal/orders/repository"
	"service/internal/utils"
    "time"
<<<<<<< HEAD
=======
========
	orderDTO "service/internal/orders/dto"
	orderRepo "service/internal/orders/repository"
	"service/internal/payments/dto"
	pay "service/internal/payments/legacy_payment"
	"service/internal/payments/repository"
	"service/internal/shared/model"
	"service/internal/shared/utils"
>>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b:internal/payments/service/payment_service.go
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b

	"github.com/google/uuid"
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
<<<<<<< HEAD
	orderRepo   *repository.ServiceOrderRepository
=======
	orderRepo   *orderRepo.ServiceOrderRepository
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
}

// NewPaymentService creates a new payment service
func NewPaymentService() *PaymentService {
	return &PaymentService{
		paymentRepo: repository.NewPaymentRepository(),
<<<<<<< HEAD
		orderRepo:   repository.NewServiceOrderRepository(),
=======
		orderRepo:   orderRepo.NewServiceOrderRepository(),
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}
}

// HandleMidtransCallback verifies payload and updates payment/order state
<<<<<<< HEAD
func (s *PaymentService) HandleMidtransCallback(ctx context.Context, cb *core.MidtransCallbackPayload, serverKey string) error {
    // Verify signature: sha512(order_id+status_code+gross_amount+server_key)
    expected := utils.SHA512Hex(cb.OrderID + cb.StatusCode + cb.GrossAmount + serverKey)
    if expected != cb.SignatureKey {
        return errors.New("invalid signature")
    }

    // Find payment by transaction ID first, fallback to invoice/order_id equals payment.ID if used as order_id
    var payment *core.Payment
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
        return core.ErrPaymentNotFound
    }

    // Idempotency: if status already mapped, no-op
    mapped := core.PaymentStatusPending
    switch cb.TransactionStatus {
    case "capture", "settlement":
        mapped = core.PaymentStatusPaid
    case "pending":
        mapped = core.PaymentStatusPending
    case "deny", "expire", "cancel":
        mapped = core.PaymentStatusFailed
    case "refund", "partial_refund":
        mapped = core.PaymentStatusRefunded
    default:
        mapped = core.PaymentStatusPending
    }
    if payment.Status == mapped {
        return nil
    }

    // Update payment
    payment.Status = mapped
    if cb.TransactionID != "" {
        payment.TransactionID = cb.TransactionID
    }
    if mapped == core.PaymentStatusPaid {
        now := core.GetCurrentTimestamp()
        payment.PaidAt = &now
    }
    if err := s.paymentRepo.Update(ctx, payment); err != nil {
        return err
    }

    // Update order on success
    if mapped == core.PaymentStatusPaid {
        if order, err := s.orderRepo.GetByID(ctx, payment.OrderID); err == nil {
            order.Status = core.StatusReady
            _ = s.orderRepo.Update(ctx, order)
        }
    }
    return nil
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, req *core.PaymentRequest) (*core.PaymentResponse, error) {
=======
func (s *PaymentService) HandleMidtransCallback(ctx context.Context, cb *dto.MidtransCallbackPayload, serverKey string) error {
	// Verify signature: sha512(order_id+status_code+gross_amount+server_key)
	expected := utils.SHA512Hex(cb.OrderID + cb.StatusCode + cb.GrossAmount + serverKey)
	if expected != cb.SignatureKey {
		return errors.New("invalid signature")
	}

	// Find payment by transaction ID first, fallback to invoice/order_id equals payment.ID if used as order_id
	var payment *dto.Payment
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
	mapped := dto.PaymentStatusPending
	switch cb.TransactionStatus {
	case "capture", "settlement":
		mapped = dto.PaymentStatusPaid
	case "pending":
		mapped = dto.PaymentStatusPending
	case "deny", "expire", "cancel":
		mapped = dto.PaymentStatusFailed
	case "refund", "partial_refund":
		mapped = dto.PaymentStatusRefunded
	default:
		mapped = dto.PaymentStatusPending
	}
	if payment.Status == mapped {
		return nil
	}

	// Update payment
	payment.Status = mapped
	if cb.TransactionID != "" {
		payment.TransactionID = cb.TransactionID
	}
	if mapped == dto.PaymentStatusPaid {
		now := model.GetCurrentTimestamp()
		payment.PaidAt = &now
	}
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}

	// Update order on success
	if mapped == dto.PaymentStatusPaid {
		if order, err := s.orderRepo.GetByID(ctx, payment.OrderID); err == nil {
			order.Status = orderDTO.StatusReady
			_ = s.orderRepo.Update(ctx, order)
		}
	}
	return nil
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, req *dto.PaymentRequest) (*dto.PaymentResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	// Validate order exists
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	_, err = s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
<<<<<<< HEAD
		return nil, core.ErrOrderNotFound
=======
		return nil, model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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

<<<<<<< HEAD
	// Create payment entity
	payment := &core.Payment{
		OrderID:       orderID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        core.PaymentStatusPending,
=======
	// Calculate PPN 11% (Indonesian tax)
	subtotal := req.Amount
	taxAmount := utils.CalculatePPN(subtotal)
	totalAmount := utils.CalculateAmountWithTax(subtotal)

	// Create payment entity
	payment := &dto.Payment{
		OrderID:       orderID,
		Subtotal:      subtotal,
		TaxAmount:     taxAmount,
		Amount:        totalAmount, // Total includes PPN
		PaymentMethod: req.PaymentMethod,
		Status:        dto.PaymentStatusPending,
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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
<<<<<<< HEAD
func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*core.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrPaymentNotFound
=======
func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrPaymentNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	response := payment.ToResponse()
	return &response, nil
}

// GetPaymentByInvoice retrieves a payment by invoice number
<<<<<<< HEAD
func (s *PaymentService) GetPaymentByInvoice(ctx context.Context, invoiceNumber string) (*core.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetByInvoiceNumber(ctx, invoiceNumber)
	if err != nil {
		return nil, core.ErrPaymentNotFound
=======
func (s *PaymentService) GetPaymentByInvoice(ctx context.Context, invoiceNumber string) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetByInvoiceNumber(ctx, invoiceNumber)
	if err != nil {
		return nil, model.ErrPaymentNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	response := payment.ToResponse()
	return &response, nil
}

// UpdatePaymentStatus updates the status of a payment
<<<<<<< HEAD
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status core.PaymentStatus, transactionID string) (*core.PaymentResponse, error) {
	// Get existing payment
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, core.ErrPaymentNotFound
=======
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status dto.PaymentStatus, transactionID string) (*dto.PaymentResponse, error) {
	// Get existing payment
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, model.ErrPaymentNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	// Update status
	payment.Status = status
	if transactionID != "" {
		payment.TransactionID = transactionID
	}

	// Set paid_at if status is paid
<<<<<<< HEAD
	if status == core.PaymentStatusPaid {
		now := core.GetCurrentTimestamp()
=======
	if status == dto.PaymentStatusPaid {
		now := model.GetCurrentTimestamp()
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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
<<<<<<< HEAD
func (s *PaymentService) ProcessMidtransPayment(ctx context.Context, req *core.MidtransPaymentRequest) (*core.MidtransPaymentResponse, error) {
=======
func (s *PaymentService) ProcessMidtransPayment(ctx context.Context, req *dto.MidtransPaymentRequest) (*dto.MidtransPaymentResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	// Validate order exists
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	_, err = s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
<<<<<<< HEAD
		return nil, core.ErrOrderNotFound
=======
		return nil, model.ErrOrderNotFound
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}

	// Generate invoice number
	invoiceNumber := utils.GenerateInvoiceNumber()

	// Create payment record
<<<<<<< HEAD
	payment := &core.Payment{
		OrderID:       orderID,
		Amount:        req.Amount,
		PaymentMethod: core.PaymentMethodMidtrans,
		Status:        core.PaymentStatusPending,
=======
	payment := &dto.Payment{
		OrderID:       orderID,
		Amount:        req.Amount,
		PaymentMethod: dto.PaymentMethodMidtrans,
		Status:        dto.PaymentStatusPending,
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		InvoiceNumber: invoiceNumber,
	}

	// Save payment
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// TODO: Integrate with actual Midtrans API
	// For now, return mock response
<<<<<<< HEAD
	response := &core.MidtransPaymentResponse{
=======
	response := &dto.MidtransPaymentResponse{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		Token:         "mock-token-" + payment.ID.String(),
		RedirectURL:   "https://app.midtrans.com/snap/v2/vtweb/" + payment.ID.String(),
		StatusCode:    "201",
		StatusMessage: "Success, transaction is created",
	}

	return response, nil
}

// ReconcilePendingPayments checks pending Midtrans payments and updates their status
func (s *PaymentService) ReconcilePendingPayments(ctx context.Context) error {
<<<<<<< HEAD
    // Get pending payments
    pendingList, err := s.paymentRepo.GetByStatus(ctx, core.PaymentStatusPending)
    if err != nil {
        return err
    }
    if len(pendingList) == 0 {
        return nil
    }
    ms := pay.NewMidtransService()
    for _, p := range pendingList {
        // Only reconcile online payments
        if p.PaymentMethod != core.PaymentMethodMidtrans && p.PaymentMethod != core.PaymentMethodGopay && p.PaymentMethod != core.PaymentMethodQris && p.PaymentMethod != core.PaymentMethodBankTransfer && p.PaymentMethod != core.PaymentMethodMandiriEchannel {
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
        newStatus := core.PaymentStatusPending
        switch resp.TransactionStatus {
        case "capture", "settlement":
            newStatus = core.PaymentStatusPaid
        case "pending":
            newStatus = core.PaymentStatusPending
        case "deny", "expire", "cancel":
            newStatus = core.PaymentStatusFailed
        case "refund", "partial_refund":
            newStatus = core.PaymentStatusRefunded
        }
        if p.Status == newStatus {
            continue
        }
        p.Status = newStatus
        if newStatus == core.PaymentStatusPaid {
            now := core.GetCurrentTimestamp()
            p.PaidAt = &now
        }
        _ = s.paymentRepo.Update(ctx, p)
        if newStatus == core.PaymentStatusPaid {
            if order, err := s.orderRepo.GetByID(ctx, p.OrderID); err == nil {
                order.Status = core.StatusReady
                _ = s.orderRepo.Update(ctx, order)
            }
        }
    }
    return nil
}

// ListPayments retrieves payments with pagination and filters
func (s *PaymentService) ListPayments(ctx context.Context, page, limit int, filters *PaymentFilters) (*core.PaginatedResponse, error) {
=======
	// Get pending payments
	pendingList, err := s.paymentRepo.GetByStatus(ctx, dto.PaymentStatusPending)
	if err != nil {
		return err
	}
	if len(pendingList) == 0 {
		return nil
	}
	ms := pay.NewMidtransService()
	for _, p := range pendingList {
		// Only reconcile online payments
		if p.PaymentMethod != dto.PaymentMethodMidtrans && p.PaymentMethod != dto.PaymentMethodGopay && p.PaymentMethod != dto.PaymentMethodQris && p.PaymentMethod != dto.PaymentMethodBankTransfer && p.PaymentMethod != dto.PaymentMethodMandiriEchannel {
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
		newStatus := dto.PaymentStatusPending
		switch resp.TransactionStatus {
		case "capture", "settlement":
			newStatus = dto.PaymentStatusPaid
		case "pending":
			newStatus = dto.PaymentStatusPending
		case "deny", "expire", "cancel":
			newStatus = dto.PaymentStatusFailed
		case "refund", "partial_refund":
			newStatus = dto.PaymentStatusRefunded
		}
		if p.Status == newStatus {
			continue
		}
		p.Status = newStatus
		if newStatus == dto.PaymentStatusPaid {
			now := model.GetCurrentTimestamp()
			p.PaidAt = &now
		}
		_ = s.paymentRepo.Update(ctx, p)
		if newStatus == dto.PaymentStatusPaid {
			if order, err := s.orderRepo.GetByID(ctx, p.OrderID); err == nil {
				order.Status = orderDTO.StatusReady
				_ = s.orderRepo.Update(ctx, order)
			}
		}
	}
	return nil
}

// ListPayments retrieves payments with pagination and filters
func (s *PaymentService) ListPayments(ctx context.Context, page, limit int, filters *PaymentFilters) (*model.PaginatedResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
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
<<<<<<< HEAD
	var responses []core.PaymentResponse
=======
	var responses []dto.PaymentResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	// Calculate pagination
	totalPages := int((total + int64(limit) - 1) / int64(limit))
<<<<<<< HEAD
	pagination := core.PaginationResponse{
=======
	pagination := model.PaginationResponse{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

<<<<<<< HEAD
	return &core.PaginatedResponse{
=======
	return &model.PaginatedResponse{
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
		Status:     "success",
		Data:       responses,
		Pagination: pagination,
		Message:    "Payments retrieved successfully",
<<<<<<< HEAD
		Timestamp:  core.GetCurrentTimestamp(),
=======
		Timestamp:  model.GetCurrentTimestamp(),
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	}, nil
}

// GetPaymentsByOrder retrieves payments for a specific order
<<<<<<< HEAD
func (s *PaymentService) GetPaymentsByOrder(ctx context.Context, orderID uuid.UUID) ([]core.PaymentResponse, error) {
=======
func (s *PaymentService) GetPaymentsByOrder(ctx context.Context, orderID uuid.UUID) ([]dto.PaymentResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	payments, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	var responses []core.PaymentResponse
=======
	var responses []dto.PaymentResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	return responses, nil
}

// GetPaymentsByStatus retrieves payments by status
<<<<<<< HEAD
func (s *PaymentService) GetPaymentsByStatus(ctx context.Context, status core.PaymentStatus) ([]core.PaymentResponse, error) {
=======
func (s *PaymentService) GetPaymentsByStatus(ctx context.Context, status dto.PaymentStatus) ([]dto.PaymentResponse, error) {
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	payments, err := s.paymentRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	var responses []core.PaymentResponse
=======
	var responses []dto.PaymentResponse
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	for _, payment := range payments {
		responses = append(responses, payment.ToResponse())
	}

	return responses, nil
}

// PaymentFilters represents filters for payment queries
type PaymentFilters struct {
	OrderID       *uuid.UUID
<<<<<<< HEAD
	Status        *core.PaymentStatus
	PaymentMethod *core.PaymentMethod
=======
	Status        *dto.PaymentStatus
	PaymentMethod *dto.PaymentMethod
>>>>>>> 62e28be2ad1dcbf35e27144a7b44a87f6b0a371b
	DateFrom      *string
	DateTo        *string
}
