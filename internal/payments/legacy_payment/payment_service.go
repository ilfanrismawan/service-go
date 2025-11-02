package legacy_payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"service/internal/shared/config"
	"service/internal/payments/dto"
	"service/internal/payments/repository"
	orderDTO "service/internal/orders/dto"
	orderRepo "service/internal/orders/repository"
	userRepo "service/internal/users/repository"
	"service/internal/shared/utils"
	"time"

	"github.com/google/uuid"
)

// MidtransService handles Midtrans payment integration
type MidtransService struct {
	serverKey    string
	clientKey    string
	isProduction bool
	baseURL      string
}

// NewMidtransService creates a new Midtrans service
func NewMidtransService() *MidtransService {
	baseURL := "https://api.sandbox.midtrans.com"
	if config.Config.MidtransIsProduction {
		baseURL = "https://api.midtrans.com"
	}

	return &MidtransService{
		serverKey:    config.Config.MidtransServerKey,
		clientKey:    config.Config.MidtransClientKey,
		isProduction: config.Config.MidtransIsProduction,
		baseURL:      baseURL,
	}
}

// MidtransPaymentRequest represents Midtrans payment request
type MidtransPaymentRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	PaymentType        string             `json:"payment_type"`
	BankTransfer       *BankTransfer      `json:"bank_transfer,omitempty"`
	Echannel           *Echannel          `json:"echannel,omitempty"`
	Gopay              *Gopay             `json:"gopay,omitempty"`
	OVO                *OVO               `json:"ovo,omitempty"`
	Dana               *Dana              `json:"dana,omitempty"`
	ShopeePay          *ShopeePay         `json:"shopeepay,omitempty"`
	Qris               *Qris              `json:"qris,omitempty"`
	Alfamart           *Alfamart          `json:"alfamart,omitempty"`
	Indomaret          *Indomaret         `json:"indomaret,omitempty"`
	CreditCard         *CreditCard         `json:"credit_card,omitempty"`
	CustomExpiry       *CustomExpiry      `json:"custom_expiry,omitempty"`
}

// TransactionDetails represents transaction details
type TransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

// BankTransfer represents bank transfer payment
type BankTransfer struct {
	Bank string `json:"bank"`
}

// Echannel represents Mandiri e-channel payment
type Echannel struct {
	BillInfo1 string `json:"bill_info1"`
	BillInfo2 string `json:"bill_info2"`
}

// Gopay represents GoPay payment
type Gopay struct {
	EnableCallback bool   `json:"enable_callback"`
	CallbackURL    string `json:"callback_url,omitempty"`
}

// Qris represents QRIS payment
type Qris struct {
	Acquirer string `json:"acquirer"`
}

// OVO represents OVO payment
type OVO struct {
	EnableCallback bool   `json:"enable_callback"`
	CallbackURL    string `json:"callback_url,omitempty"`
}

// Dana represents Dana payment
type Dana struct {
	EnableCallback bool   `json:"enable_callback"`
	CallbackURL    string `json:"callback_url,omitempty"`
}

// ShopeePay represents ShopeePay payment
type ShopeePay struct {
	EnableCallback bool   `json:"enable_callback"`
	CallbackURL    string `json:"callback_url,omitempty"`
}

// Alfamart represents Alfamart payment
type Alfamart struct {
	Store   string `json:"store,omitempty"`
	Message string `json:"message,omitempty"`
}

// Indomaret represents Indomaret payment
type Indomaret struct {
	Store   string `json:"store,omitempty"`
	Message string `json:"message,omitempty"`
}

// CreditCard represents Credit Card payment
type CreditCard struct {
	TokenID    string `json:"token_id,omitempty"`
	Secure     bool   `json:"secure,omitempty"`
	SaveToken  bool   `json:"save_token_id,omitempty"`
	Bank       string `json:"bank,omitempty"`
	Installment bool  `json:"installment,omitempty"`
}

// CustomExpiry represents custom expiry settings
type CustomExpiry struct {
	OrderTime      string `json:"order_time"`
	ExpiryDuration int    `json:"expiry_duration"`
	Unit           string `json:"unit"`
}

// MidtransPaymentResponse represents Midtrans payment response
type MidtransPaymentResponse struct {
	StatusCode             string                 `json:"status_code"`
	StatusMessage          string                 `json:"status_message"`
	TransactionID          string                 `json:"transaction_id"`
	OrderID                string                 `json:"order_id"`
	MerchantID             string                 `json:"merchant_id"`
	GrossAmount            string                 `json:"gross_amount"`
	Currency               string                 `json:"currency"`
	PaymentType            string                 `json:"payment_type"`
	TransactionTime        string                 `json:"transaction_time"`
	TransactionStatus      string                 `json:"transaction_status"`
	FraudStatus            string                 `json:"fraud_status"`
	Actions                []Action               `json:"actions"`
	ChannelResponseCode    string                 `json:"channel_response_code"`
	ChannelResponseMessage string                 `json:"channel_response_message"`
	VaNumbers              []VaNumber             `json:"va_numbers"`
	PaymentCode            string                 `json:"payment_code"`
	QRString               string                 `json:"qr_string"`
	RedirectURL            string                 `json:"redirect_url"`
	AdditionalData         map[string]interface{} `json:"additional_data"`
}

// Action represents payment action
type Action struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

// VaNumber represents virtual account number
type VaNumber struct {
	Bank     string `json:"bank"`
	VaNumber string `json:"va_number"`
}

// PaymentService handles payment business logic
type PaymentService struct {
	midtransService *MidtransService
	paymentRepo     *repository.PaymentRepository
	orderRepo       *orderRepo.ServiceOrderRepository
	userRepo        *userRepo.UserRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService() *PaymentService {
	return &PaymentService{
		midtransService: NewMidtransService(),
		paymentRepo:     repository.NewPaymentRepository(),
		orderRepo:       orderRepo.NewServiceOrderRepository(),
		userRepo:        userRepo.NewUserRepository(),
	}
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, orderID uuid.UUID, paymentMethod dto.PaymentMethod, amount float64) (*dto.Payment, error) {
	// Get order details
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Create payment record
	payment := &dto.Payment{
		ID:            uuid.New(),
		OrderID:       orderID,
		UserID:        order.UserID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        dto.PaymentStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save payment to database
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return payment, nil
}

// ProcessPayment processes payment through Midtrans
func (s *PaymentService) ProcessPayment(ctx context.Context, paymentID uuid.UUID, paymentMethod dto.PaymentMethod) (*MidtransPaymentResponse, error) {
	// Get payment details
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Get order details
	order, err := s.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// For cash payments, mark as paid immediately
	if paymentMethod == dto.PaymentMethodCash {
		payment.Status = dto.PaymentStatusPaid
		payment.TransactionID = fmt.Sprintf("CASH-%d", time.Now().Unix())
		payment.UpdatedAt = time.Now()

		if err := s.paymentRepo.Update(ctx, payment); err != nil {
			return nil, fmt.Errorf("failed to update cash payment: %w", err)
		}

		// Return mock response for cash payment
		return &MidtransPaymentResponse{
			StatusCode:        "200",
			StatusMessage:     "Cash payment processed successfully",
			TransactionID:     payment.TransactionID,
			OrderID:           payment.ID.String(),
			GrossAmount:       fmt.Sprintf("%.0f", payment.Amount),
			PaymentType:       "cash",
			TransactionStatus: "settlement",
		}, nil
	}

	// Create Midtrans payment request for online payments
	midtransReq := &MidtransPaymentRequest{
		TransactionDetails: TransactionDetails{
			OrderID:     payment.ID.String(),
			GrossAmount: int64(payment.Amount * 100), // Convert to cents
		},
		PaymentType: string(paymentMethod),
		CustomExpiry: &CustomExpiry{
			OrderTime:      time.Now().Format("2006-01-02 15:04:05"),
			ExpiryDuration: 24,
			Unit:           "hour",
		},
	}

	// Set payment method specific parameters
	switch paymentMethod {
	case dto.PaymentMethodBankTransfer, dto.PaymentMethodBCAVA:
		midtransReq.BankTransfer = &BankTransfer{
			Bank: "bca",
		}
		midtransReq.PaymentType = "bank_transfer"
	case dto.PaymentMethodBNIVA:
		midtransReq.BankTransfer = &BankTransfer{
			Bank: "bni",
		}
		midtransReq.PaymentType = "bank_transfer"
	case dto.PaymentMethodBRIVA:
		midtransReq.BankTransfer = &BankTransfer{
			Bank: "bri",
		}
		midtransReq.PaymentType = "bank_transfer"
	case dto.PaymentMethodPermataVA:
		midtransReq.BankTransfer = &BankTransfer{
			Bank: "permata",
		}
		midtransReq.PaymentType = "bank_transfer"
	case dto.PaymentMethodMandiriEchannel:
		midtransReq.Echannel = &Echannel{
			BillInfo1: fmt.Sprintf("Order #%s", order.OrderNumber),
			BillInfo2: "iPhone Service",
		}
		midtransReq.PaymentType = "echannel"
	case dto.PaymentMethodGopay:
		midtransReq.Gopay = &Gopay{
			EnableCallback: true,
			CallbackURL:    fmt.Sprintf("%s/api/v1/payments/midtrans/callback", config.Config.BaseURL),
		}
		midtransReq.PaymentType = "gopay"
	case dto.PaymentMethodOVO:
		midtransReq.OVO = &OVO{
			EnableCallback: true,
			CallbackURL:    fmt.Sprintf("%s/api/v1/payments/midtrans/callback", config.Config.BaseURL),
		}
		midtransReq.PaymentType = "ovo"
	case dto.PaymentMethodDana:
		midtransReq.Dana = &Dana{
			EnableCallback: true,
			CallbackURL:    fmt.Sprintf("%s/api/v1/payments/midtrans/callback", config.Config.BaseURL),
		}
		midtransReq.PaymentType = "dana"
	case dto.PaymentMethodShopeePay:
		midtransReq.ShopeePay = &ShopeePay{
			EnableCallback: true,
			CallbackURL:    fmt.Sprintf("%s/api/v1/payments/midtrans/callback", config.Config.BaseURL),
		}
		midtransReq.PaymentType = "shopeepay"
	case dto.PaymentMethodQRIS:
		midtransReq.Qris = &Qris{
			Acquirer: "gopay",
		}
		midtransReq.PaymentType = "qris"
	case dto.PaymentMethodAlfamart:
		midtransReq.Alfamart = &Alfamart{
			Store:   "alfamart",
			Message: fmt.Sprintf("Bayar pesanan #%s", order.OrderNumber),
		}
		midtransReq.PaymentType = "cstore"
	case dto.PaymentMethodIndomaret:
		midtransReq.Indomaret = &Indomaret{
			Store:   "indomaret",
			Message: fmt.Sprintf("Bayar pesanan #%s", order.OrderNumber),
		}
		midtransReq.PaymentType = "cstore"
	case dto.PaymentMethodCreditCard:
		midtransReq.CreditCard = &CreditCard{
			Secure:    true,
			SaveToken: false,
		}
		midtransReq.PaymentType = "credit_card"
	default:
		// Default to bank_transfer
		midtransReq.BankTransfer = &BankTransfer{
			Bank: "bca",
		}
		midtransReq.PaymentType = "bank_transfer"
	}

	// Send request to Midtrans (mock implementation)
	response, err := s.midtransService.CreatePayment(ctx, midtransReq)
	if err != nil {
		// Update payment status to failed
		payment.Status = dto.PaymentStatusFailed
		payment.UpdatedAt = time.Now()
		s.paymentRepo.Update(ctx, payment)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Update payment with transaction details
	payment.TransactionID = response.TransactionID
	payment.Status = dto.PaymentStatusPending
	payment.UpdatedAt = time.Now()

	// Save payment URL based on payment method
	switch paymentMethod {
	case dto.PaymentMethodBankTransfer:
		if len(response.VaNumbers) > 0 {
			payment.InvoiceURL = fmt.Sprintf("Virtual Account: %s", response.VaNumbers[0].VaNumber)
		}
	case dto.PaymentMethodMandiriEchannel:
		payment.InvoiceURL = fmt.Sprintf("Bill Key: %s", response.PaymentCode)
	case dto.PaymentMethodGopay:
		payment.InvoiceURL = response.RedirectURL
	case dto.PaymentMethodQris:
		payment.InvoiceURL = response.QRString
	}

	// Update payment in database
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return response, nil
}

// CreatePayment creates payment through Midtrans API
func (ms *MidtransService) CreatePayment(ctx context.Context, req *MidtransPaymentRequest) (*MidtransPaymentResponse, error) {
	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/v2/charge", ms.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+utils.EncodeBase64(ms.serverKey+":"))

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("midtrans API error: %s", string(body))
	}

	// Parse response
	var midtransResp MidtransPaymentResponse
	if err := json.Unmarshal(body, &midtransResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &midtransResp, nil
}

// GetPaymentStatus gets payment status from Midtrans
func (ms *MidtransService) GetPaymentStatus(ctx context.Context, transactionID string) (*MidtransPaymentResponse, error) {
	// Create HTTP request
	url := fmt.Sprintf("%s/v2/%s/status", ms.baseURL, transactionID)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+utils.EncodeBase64(ms.serverKey+":"))

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("midtrans API error: %s", string(body))
	}

	// Parse response
	var midtransResp MidtransPaymentResponse
	if err := json.Unmarshal(body, &midtransResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &midtransResp, nil
}

// HandlePaymentCallback handles payment callback from Midtrans
func (s *PaymentService) HandlePaymentCallback(ctx context.Context, transactionID string, status string) error {
	// Get payment by transaction ID
	payment, err := s.paymentRepo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Update payment status based on Midtrans status
	var newStatus dto.PaymentStatus
	switch status {
	case "capture", "settlement":
		newStatus = dto.PaymentStatusPaid
	case "pending":
		newStatus = dto.PaymentStatusPending
	case "deny", "expire", "cancel":
		newStatus = dto.PaymentStatusFailed
	case "refund":
		newStatus = dto.PaymentStatusRefunded
	default:
		newStatus = dto.PaymentStatusPending
	}

	// Update payment status
	payment.Status = newStatus
	payment.UpdatedAt = time.Now()

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Update order status if payment is successful
	if newStatus == dto.PaymentStatusPaid {
		order, err := s.orderRepo.GetByID(ctx, payment.OrderID)
		if err == nil {
			// Update order status to confirmed
			order.Status = orderDTO.StatusPendingPickup
			order.UpdatedAt = time.Now()
			s.orderRepo.Update(ctx, order)
		}
	}

	return nil
}

// RefundPayment processes payment refund
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID uuid.UUID, amount float64, reason string) error {
	// Get payment details
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Check if payment is eligible for refund
	if payment.Status != dto.PaymentStatusPaid {
		return fmt.Errorf("payment is not eligible for refund")
	}

	// Process refund through Midtrans (mock implementation)
	// In production, implement actual Midtrans refund API
	log.Printf("Processing refund for payment %s: amount %.2f, reason: %s", paymentID.String(), amount, reason)

	// Update payment status
	payment.Status = dto.PaymentStatusRefunded
	payment.UpdatedAt = time.Now()

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

// GenerateInvoice generates invoice for payment
func (s *PaymentService) GenerateInvoice(ctx context.Context, paymentID uuid.UUID) (*dto.Invoice, error) {
	// Get payment details
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Get order details
	order, err := s.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get user details
	user, err := s.userRepo.GetByID(ctx, payment.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate invoice
	invoice := &dto.Invoice{
		InvoiceNumber: utils.GenerateInvoiceNumber(),
		OrderNumber:   order.OrderNumber,
		CustomerName:  user.Name,
		CustomerEmail: user.Email,
		CustomerPhone: user.Phone,
		Amount:        payment.Amount,
		PaymentMethod: string(payment.PaymentMethod),
		Status:        string(payment.Status),
		CreatedAt:     payment.CreatedAt,
		DueDate:       payment.CreatedAt.Add(24 * time.Hour),
	}

	return invoice, nil
}
