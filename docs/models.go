package docs

// @title iPhone Service POS API
// @version 1.0
// @description Backend API for iPhone Service Point of Sales system supporting 50 branches across Indonesia
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.iphoneservice.com/support
// @contact.email support@iphoneservice.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name auth
// @tag.description Authentication endpoints for user registration, login, and profile management

// @tag.name branches
// @tag.description Branch management endpoints for managing service outlets across Indonesia

// @tag.name orders
// @tag.description Order management endpoints for iPhone service orders

// @tag.name payments
// @tag.description Payment processing endpoints with Midtrans integration

// @tag.name notifications
// @tag.description Notification endpoints for order updates and system messages

// @tag.name files
// @tag.description File upload endpoints for photos and documents

// @tag.name chat
// @tag.description Chat endpoints for customer-technician communication

// @tag.name dashboard
// @tag.description Dashboard analytics endpoints for business insights

// @tag.name health
// @tag.description Health check endpoints for monitoring

// @tag.name admin
// @tag.description Admin endpoints for system management

// @tag.name cashier
// @tag.description Cashier endpoints for order processing

// @tag.name technician
// @tag.description Technician endpoints for service management

// @tag.name courier
// @tag.description Courier endpoints for pickup and delivery

// Response models
type APIResponse struct {
	Status    string      `json:"status" example:"success"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message" example:"Operation completed successfully"`
	Timestamp string      `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}

// Alias for swagger compatibility (core.APIResponse -> docs.APIResponse)
type CoreAPIResponse = APIResponse

type ErrorResponse struct {
	Status    string      `json:"status" example:"error"`
	Error     string      `json:"error" example:"validation_error"`
	Message   string      `json:"message" example:"Invalid request data"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp string      `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}

// Alias for swagger compatibility (core.ErrorResponse -> docs.ErrorResponse)
type CoreErrorResponse = ErrorResponse

type PaginatedResponse struct {
	Status     string      `json:"status" example:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Message    string      `json:"message" example:"Data retrieved successfully"`
	Timestamp  string      `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}

type Pagination struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}

// User models
type UserResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string `json:"name" example:"John Doe"`
	Email     string `json:"email" example:"john.doe@example.com"`
	Phone     string `json:"phone" example:"081234567890"`
	Role      string `json:"role" example:"customer"`
	IsActive  bool   `json:"is_active" example:"true"`
	CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type UserRequest struct {
	Name     string `json:"name" example:"John Doe" validate:"required,min=2,max=100"`
	Email    string `json:"email" example:"john.doe@example.com" validate:"required,email"`
	Phone    string `json:"phone" example:"081234567890" validate:"required,min=10,max=15"`
	Password string `json:"password" example:"password123" validate:"required,min=6"`
	Role     string `json:"role" example:"customer" validate:"required,oneof=customer cashier technician courier branch_admin central_admin"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"john.doe@example.com" validate:"required,email"`
	Password string `json:"password" example:"password123" validate:"required"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int          `json:"expires_in" example:"86400"`
}

// Branch models
type BranchResponse struct {
	ID        string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string  `json:"name" example:"Jakarta Central"`
	Address   string  `json:"address" example:"Jl. Sudirman No. 123"`
	City      string  `json:"city" example:"Jakarta"`
	Province  string  `json:"province" example:"DKI Jakarta"`
	Phone     string  `json:"phone" example:"021-12345678"`
	Latitude  float64 `json:"latitude" example:"-6.2088"`
	Longitude float64 `json:"longitude" example:"106.8456"`
	IsActive  bool    `json:"is_active" example:"true"`
	CreatedAt string  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type BranchRequest struct {
	Name      string  `json:"name" example:"Jakarta Central" validate:"required,min=2,max=100"`
	Address   string  `json:"address" example:"Jl. Sudirman No. 123" validate:"required,min=10,max=200"`
	City      string  `json:"city" example:"Jakarta" validate:"required,min=2,max=50"`
	Province  string  `json:"province" example:"DKI Jakarta" validate:"required,min=2,max=50"`
	Phone     string  `json:"phone" example:"021-12345678" validate:"required,min=8,max=15"`
	Latitude  float64 `json:"latitude" example:"-6.2088" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" example:"106.8456" validate:"required,min=-180,max=180"`
}

type BranchDistance struct {
	Branch   BranchResponse `json:"branch"`
	Distance float64        `json:"distance" example:"2.5"`
}

// Order models
type ServiceOrderResponse struct {
	ID                string   `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderNumber       string   `json:"order_number" example:"ORD-20240101-001"`
	UserID            string   `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	BranchID          string   `json:"branch_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CourierID         *string  `json:"courier_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	TechnicianID      *string  `json:"technician_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	iPhoneType        string   `json:"iphone_type" example:"iPhone 14 Pro"`
	Complaint         string   `json:"complaint" example:"Screen cracked, needs replacement"`
	PickupLocation    string   `json:"pickup_location" example:"Jakarta Selatan"`
	Status            string   `json:"status" example:"in_service"`
	ServiceCost       float64  `json:"service_cost" example:"500000"`
	EstimatedDuration int      `json:"estimated_duration" example:"3"`
	PhotoURLs         []string `json:"photo_urls" example:"http://example.com/photo1.jpg"`
	CreatedAt         string   `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt         string   `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type ServiceOrderRequest struct {
	iPhoneType        string  `json:"iphone_type" example:"iPhone 14 Pro" validate:"required,min=2,max=50"`
	Complaint         string  `json:"complaint" example:"Screen cracked, needs replacement" validate:"required,min=10,max=500"`
	PickupLocation    string  `json:"pickup_location" example:"Jakarta Selatan" validate:"required,min=5,max=100"`
	BranchID          string  `json:"branch_id" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required,uuid"`
	EstimatedCost     float64 `json:"estimated_cost" example:"500000" validate:"required,min=0"`
	EstimatedDuration int     `json:"estimated_duration" example:"3" validate:"required,min=1,max=30"`
}

// Payment models
type PaymentResponse struct {
	ID            string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderID       string  `json:"order_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount        float64 `json:"amount" example:"500000"`
	PaymentMethod string  `json:"payment_method" example:"midtrans"`
	Status        string  `json:"status" example:"paid"`
	TransactionID string  `json:"transaction_id,omitempty" example:"TXN-001"`
	InvoiceNumber string  `json:"invoice_number" example:"INV-20240101-001"`
	PaidAt        *string `json:"paid_at,omitempty" example:"2024-01-01T00:00:00Z"`
	CreatedAt     string  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     string  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type PaymentRequest struct {
	OrderID       string  `json:"order_id" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required,uuid"`
	Amount        float64 `json:"amount" example:"500000" validate:"required,min=0"`
	PaymentMethod string  `json:"payment_method" example:"midtrans" validate:"required,oneof=cash midtrans gopay qris transfer"`
	Notes         string  `json:"notes,omitempty" example:"Payment for iPhone repair"`
}

type MidtransPaymentResponse struct {
	Token         string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RedirectURL   string `json:"redirect_url" example:"https://app.midtrans.com/snap/v2/vtweb/..."`
	StatusCode    string `json:"status_code" example:"201"`
	StatusMessage string `json:"status_message" example:"Success"`
}

// Notification models
type NotificationResponse struct {
	ID        string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderID   *string `json:"order_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type      string  `json:"type" example:"order_update"`
	Message   string  `json:"message" example:"Your order is now in service"`
	IsRead    bool    `json:"is_read" example:"false"`
	CreatedAt string  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type NotificationRequest struct {
	UserID  string  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required,uuid"`
	OrderID *string `json:"order_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000" validate:"omitempty,uuid"`
	Type    string  `json:"type" example:"order_update" validate:"required,oneof=order_update order_ready order_delivered order_completed order_cancelled payment_received payment_failed welcome promotion system"`
	Message string  `json:"message" example:"Your order is now in service" validate:"required,min=5,max=500"`
}

// File upload models
type FileUploadResponse struct {
	Filename     string `json:"filename" example:"orders_550e8400-e29b-41d4-a716-446655440000_pickup.jpg"`
	OriginalName string `json:"original_name" example:"iphone_damage.jpg"`
	Size         int64  `json:"size" example:"1024000"`
	ContentType  string `json:"content_type" example:"image/jpeg"`
	URL          string `json:"url" example:"http://localhost:9000/iphone-service/orders/550e8400-e29b-41d4-a716-446655440000/pickup/orders_550e8400-e29b-41d4-a716-446655440000_pickup.jpg"`
	UploadedAt   string `json:"uploaded_at" example:"2024-01-01T00:00:00Z"`
}

// Chat models
type ChatMessageResponse struct {
	ID         string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderID    string `json:"order_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	SenderID   string `json:"sender_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ReceiverID string `json:"receiver_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Message    string `json:"message" example:"How is the repair progress?"`
	CreatedAt  string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type ChatMessageRequest struct {
	OrderID    string `json:"order_id" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required,uuid"`
	ReceiverID string `json:"receiver_id" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required,uuid"`
	Message    string `json:"message" example:"How is the repair progress?" validate:"required,min=1,max=1000"`
}

// Dashboard models
type DashboardOverview struct {
	TotalOrders     int64   `json:"total_orders" example:"150"`
	TotalRevenue    float64 `json:"total_revenue" example:"75000000"`
	TotalCustomers  int64   `json:"total_customers" example:"1200"`
	ActiveBranches  int64   `json:"active_branches" example:"50"`
	PendingOrders   int64   `json:"pending_orders" example:"25"`
	CompletedOrders int64   `json:"completed_orders" example:"125"`
}

type OrderStats struct {
	TotalOrders     int64   `json:"total_orders" example:"150"`
	PendingOrders   int64   `json:"pending_orders" example:"25"`
	InServiceOrders int64   `json:"in_service_orders" example:"50"`
	ReadyOrders     int64   `json:"ready_orders" example:"30"`
	CompletedOrders int64   `json:"completed_orders" example:"125"`
	CancelledOrders int64   `json:"cancelled_orders" example:"5"`
	TotalRevenue    float64 `json:"total_revenue" example:"75000000"`
}

type RevenueStats struct {
	TotalRevenue      float64            `json:"total_revenue" example:"75000000"`
	MonthlyRevenue    float64            `json:"monthly_revenue" example:"25000000"`
	DailyRevenue      float64            `json:"daily_revenue" example:"800000"`
	AverageOrderValue float64            `json:"average_order_value" example:"500000"`
	PaymentMethods    map[string]float64 `json:"payment_methods" example:"{\"cash\": 30000000, \"midtrans\": 25000000, \"gopay\": 20000000}"`
}

type BranchStats struct {
	TotalBranches    int64              `json:"total_branches" example:"50"`
	ActiveBranches   int64              `json:"active_branches" example:"48"`
	InactiveBranches int64              `json:"inactive_branches" example:"2"`
	OrdersByBranch   map[string]int64   `json:"orders_by_branch" example:"{\"Jakarta Central\": 25, \"Surabaya Branch\": 20}"`
	RevenueByBranch  map[string]float64 `json:"revenue_by_branch" example:"{\"Jakarta Central\": 12500000, \"Surabaya Branch\": 10000000}"`
}

// Health check models
type HealthResponse struct {
	Status    string                 `json:"status" example:"healthy"`
	Timestamp string                 `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Version   string                 `json:"version" example:"1.0.0"`
	Services  map[string]interface{} `json:"services"`
}

type ServiceStatus struct {
	Service string                 `json:"service" example:"postgresql"`
	Healthy bool                   `json:"healthy" example:"true"`
	Error   *string                `json:"error,omitempty" example:"Connection timeout"`
	Stats   map[string]interface{} `json:"stats,omitempty"`
}
