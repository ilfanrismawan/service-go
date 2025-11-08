package dto

import (
	"time"

	"github.com/google/uuid"

	paymentEntity "service-go/internal/modules/payments/entity"
	orderDto "service-go/internal/modules/orders/dto"
)

// PaymentRequest represents the request payload for creating a payment
type PaymentRequest struct {
	OrderID       string                  `json:"order_id" validate:"required"`
	Amount        float64                 `json:"amount" validate:"required,gt=0"`
	PaymentMethod paymentEntity.PaymentMethod `json:"payment_method" validate:"required"`
	Notes         string                  `json:"notes,omitempty"`
}

// PaymentResponse represents the response payload for payment data
type PaymentResponse struct {
	ID            uuid.UUID                    `json:"id"`
	OrderID       uuid.UUID                    `json:"order_id"`
	Order         orderDto.ServiceOrderResponse `json:"order"`
	Amount        float64                      `json:"amount"`
	PaymentMethod paymentEntity.PaymentMethod  `json:"payment_method"`
	Status        paymentEntity.PaymentStatus   `json:"status"`
	TransactionID string                       `json:"transaction_id,omitempty"`
	InvoiceNumber string                       `json:"invoice_number"`
	PaidAt        *time.Time                   `json:"paid_at,omitempty"`
	Notes         string                       `json:"notes,omitempty"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
}

// MidtransPaymentRequest represents the request payload for Midtrans payment
type MidtransPaymentRequest struct {
	OrderID       string  `json:"order_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	CustomerEmail string  `json:"customer_email" validate:"required,email"`
	CustomerPhone string  `json:"customer_phone" validate:"required,phone"`
}

// MidtransPaymentResponse represents the response from Midtrans API
type MidtransPaymentResponse struct {
	Token         string `json:"token"`
	RedirectURL   string `json:"redirect_url"`
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// MidtransCallbackPayload represents Midtrans callback/webhook payload
type MidtransCallbackPayload struct {
	OrderID           string `json:"order_id"`
	TransactionID     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
}

// Invoice represents an invoice
type Invoice struct {
	InvoiceNumber string    `json:"invoice_number"`
	OrderNumber   string    `json:"order_number"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	CustomerPhone string    `json:"customer_phone"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	DueDate       time.Time `json:"due_date"`
}

// ToPaymentResponse converts Payment entity to PaymentResponse DTO
func ToPaymentResponse(p *paymentEntity.Payment) PaymentResponse {
	return PaymentResponse{
		ID:            p.ID,
		OrderID:       p.OrderID,
		Amount:        p.Amount,
		PaymentMethod: p.PaymentMethod,
		Status:        p.Status,
		TransactionID: p.TransactionID,
		InvoiceNumber: p.InvoiceNumber,
		PaidAt:        p.PaidAt,
		Notes:         p.Notes,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

