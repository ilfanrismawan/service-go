package core

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

const (
	PaymentMethodCash            PaymentMethod = "cash"
	PaymentMethodMidtrans        PaymentMethod = "midtrans"
	PaymentMethodGopay           PaymentMethod = "gopay"
	PaymentMethodQRIS            PaymentMethod = "qris"
	PaymentMethodQris            PaymentMethod = "qris" // Alias for QRIS
	PaymentMethodTransfer        PaymentMethod = "transfer"
	PaymentMethodBankTransfer    PaymentMethod = "bank_transfer"
	PaymentMethodMandiriEchannel PaymentMethod = "mandiri_echannel"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// Payment represents a payment transaction
type Payment struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderID       uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	UserID        uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"` // Alias for Order.CustomerID
	Order         ServiceOrder   `json:"order" gorm:"foreignKey:OrderID"`
	Amount        float64        `json:"amount" gorm:"not null"`
	PaymentMethod PaymentMethod  `json:"payment_method" gorm:"not null"`
	Status        PaymentStatus  `json:"status" gorm:"not null;default:'pending'"`
	TransactionID string         `json:"transaction_id,omitempty"` // External payment gateway transaction ID
	InvoiceNumber string         `json:"invoice_number" gorm:"uniqueIndex;not null"`
	InvoiceURL    string         `json:"invoice_url,omitempty"` // URL to the invoice
	PaidAt        *time.Time     `json:"paid_at,omitempty"`
	Notes         string         `json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for Payment
func (Payment) TableName() string {
	return "payments"
}

// PaymentRequest represents the request payload for creating a payment
type PaymentRequest struct {
	OrderID       string        `json:"order_id" validate:"required"`
	Amount        float64       `json:"amount" validate:"required,gt=0"`
	PaymentMethod PaymentMethod `json:"payment_method" validate:"required"`
	Notes         string        `json:"notes,omitempty"`
}

// PaymentResponse represents the response payload for payment data
type PaymentResponse struct {
	ID            uuid.UUID            `json:"id"`
	OrderID       uuid.UUID            `json:"order_id"`
	Order         ServiceOrderResponse `json:"order"`
	Amount        float64              `json:"amount"`
	PaymentMethod PaymentMethod        `json:"payment_method"`
	Status        PaymentStatus        `json:"status"`
	TransactionID string               `json:"transaction_id,omitempty"`
	InvoiceNumber string               `json:"invoice_number"`
	PaidAt        *time.Time           `json:"paid_at,omitempty"`
	Notes         string               `json:"notes,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// ToResponse converts Payment to PaymentResponse
func (p *Payment) ToResponse() PaymentResponse {
	return PaymentResponse{
		ID:            p.ID,
		OrderID:       p.OrderID,
		Order:         p.Order.ToResponse(),
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

// SetUserID sets the UserID field from the Order's CustomerID
func (p *Payment) SetUserID() {
	p.UserID = p.Order.CustomerID
}

// MidtransPaymentRequest represents the request payload for Midtrans payment
type MidtransPaymentRequest struct {
	OrderID       string  `json:"order_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	CustomerEmail string  `json:"customer_email" validate:"required,email"`
	CustomerPhone string  `json:"customer_phone" validate:"required"`
}

// MidtransPaymentResponse represents the response from Midtrans API
type MidtransPaymentResponse struct {
	Token         string `json:"token"`
	RedirectURL   string `json:"redirect_url"`
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
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
