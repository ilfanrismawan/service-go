package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	orderEntity "service-go/internal/modules/orders/entity"
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
	PaymentMethodShopeePay       PaymentMethod = "shopeepay"
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
	Order         orderEntity.ServiceOrder `json:"order" gorm:"foreignKey:OrderID"`
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

// SetUserID sets the UserID field from the Order's CustomerID
func (p *Payment) SetUserID() {
	p.UserID = p.Order.CustomerID
}

