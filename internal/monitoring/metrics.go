package monitoring

import (
	"service/internal/core"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPRequestsInFlight  *prometheus.GaugeVec

	// Business metrics
	OrdersTotal          *prometheus.CounterVec
	OrdersByStatus       *prometheus.GaugeVec
	OrdersByBranch       *prometheus.GaugeVec
	OrdersByTechnician   *prometheus.GaugeVec
	OrdersByCourier      *prometheus.GaugeVec

	// Payment metrics
	PaymentsTotal        *prometheus.CounterVec
	PaymentsByMethod     *prometheus.CounterVec
	PaymentsByStatus     *prometheus.GaugeVec
	PaymentAmount        *prometheus.HistogramVec

	// User metrics
	UsersTotal           *prometheus.GaugeVec
	UsersByRole          *prometheus.GaugeVec
	UsersByBranch        *prometheus.GaugeVec
	ActiveUsers          *prometheus.GaugeVec

	// Branch metrics
	BranchesTotal        *prometheus.GaugeVec
	BranchesByCity       *prometheus.GaugeVec
	BranchesByProvince   *prometheus.GaugeVec
	ActiveBranches       *prometheus.GaugeVec

	// Notification metrics
	NotificationsTotal   *prometheus.CounterVec
	NotificationsByType  *prometheus.CounterVec
	NotificationsByUser  *prometheus.CounterVec
	UnreadNotifications  *prometheus.GaugeVec

	// System metrics
	DatabaseConnections  *prometheus.GaugeVec
	RedisConnections     *prometheus.GaugeVec
	CacheHitRate         *prometheus.GaugeVec
	CacheMissRate        *prometheus.GaugeVec
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		// HTTP metrics
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		HTTPRequestsInFlight: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed",
			},
			[]string{"method", "endpoint"},
		),

		// Business metrics
		OrdersTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "orders_total",
				Help: "Total number of orders created",
			},
			[]string{"branch_id", "user_id"},
		),
		OrdersByStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "orders_by_status",
				Help: "Number of orders by status",
			},
			[]string{"status"},
		),
		OrdersByBranch: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "orders_by_branch",
				Help: "Number of orders by branch",
			},
			[]string{"branch_id", "branch_name"},
		),
		OrdersByTechnician: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "orders_by_technician",
				Help: "Number of orders by technician",
			},
			[]string{"technician_id", "technician_name"},
		),
		OrdersByCourier: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "orders_by_courier",
				Help: "Number of orders by courier",
			},
			[]string{"courier_id", "courier_name"},
		),

		// Payment metrics
		PaymentsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payments_total",
				Help: "Total number of payments processed",
			},
			[]string{"payment_method", "status"},
		),
		PaymentsByMethod: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payments_by_method_total",
				Help: "Total number of payments by method",
			},
			[]string{"payment_method"},
		),
		PaymentsByStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "payments_by_status",
				Help: "Number of payments by status",
			},
			[]string{"status"},
		),
		PaymentAmount: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "payment_amount",
				Help:    "Payment amount distribution",
				Buckets: []float64{100000, 250000, 500000, 750000, 1000000, 2000000, 5000000},
			},
			[]string{"payment_method"},
		),

		// User metrics
		UsersTotal: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "users_total",
				Help: "Total number of users",
			},
			[]string{"role"},
		),
		UsersByRole: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "users_by_role",
				Help: "Number of users by role",
			},
			[]string{"role"},
		),
		UsersByBranch: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "users_by_branch",
				Help: "Number of users by branch",
			},
			[]string{"branch_id", "branch_name"},
		),
		ActiveUsers: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "active_users",
				Help: "Number of active users",
			},
			[]string{"role"},
		),

		// Branch metrics
		BranchesTotal: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "branches_total",
				Help: "Total number of branches",
			},
			[]string{"city", "province"},
		),
		BranchesByCity: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "branches_by_city",
				Help: "Number of branches by city",
			},
			[]string{"city"},
		),
		BranchesByProvince: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "branches_by_province",
				Help: "Number of branches by province",
			},
			[]string{"province"},
		),
		ActiveBranches: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "active_branches",
				Help: "Number of active branches",
			},
			[]string{"city", "province"},
		),

		// Notification metrics
		NotificationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "notifications_total",
				Help: "Total number of notifications sent",
			},
			[]string{"type", "channel"},
		),
		NotificationsByType: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "notifications_by_type_total",
				Help: "Total number of notifications by type",
			},
			[]string{"type"},
		),
		NotificationsByUser: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "notifications_by_user_total",
				Help: "Total number of notifications by user",
			},
			[]string{"user_id", "user_role"},
		),
		UnreadNotifications: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "unread_notifications",
				Help: "Number of unread notifications",
			},
			[]string{"user_id", "user_role"},
		),

		// System metrics
		DatabaseConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "database_connections",
				Help: "Number of database connections",
			},
			[]string{"state"},
		),
		RedisConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "redis_connections",
				Help: "Number of Redis connections",
			},
			[]string{"state"},
		),
		CacheHitRate: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cache_hit_rate",
				Help: "Cache hit rate",
			},
			[]string{"cache_type"},
		),
		CacheMissRate: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cache_miss_rate",
				Help: "Cache miss rate",
			},
			[]string{"cache_type"},
		),
	}
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(method, endpoint, statusCode string, duration time.Duration) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordOrderCreated records order creation metrics
func (m *Metrics) RecordOrderCreated(order *core.ServiceOrder) {
	m.OrdersTotal.WithLabelValues(order.BranchID.String(), order.CustomerID.String()).Inc()
	m.OrdersByStatus.WithLabelValues(string(order.Status)).Inc()
	m.OrdersByBranch.WithLabelValues(order.BranchID.String(), "").Inc()
}

// RecordOrderStatusChange records order status change metrics
func (m *Metrics) RecordOrderStatusChange(order *core.ServiceOrder, oldStatus, newStatus core.OrderStatus) {
	m.OrdersByStatus.WithLabelValues(string(oldStatus)).Dec()
	m.OrdersByStatus.WithLabelValues(string(newStatus)).Inc()
}

// RecordPaymentProcessed records payment processing metrics
func (m *Metrics) RecordPaymentProcessed(payment *core.Payment) {
	m.PaymentsTotal.WithLabelValues(string(payment.PaymentMethod), string(payment.Status)).Inc()
	m.PaymentsByMethod.WithLabelValues(string(payment.PaymentMethod)).Inc()
	m.PaymentsByStatus.WithLabelValues(string(payment.Status)).Inc()
	m.PaymentAmount.WithLabelValues(string(payment.PaymentMethod)).Observe(payment.Amount)
}

// RecordNotificationSent records notification sending metrics
func (m *Metrics) RecordNotificationSent(notification *core.Notification, channel string) {
	m.NotificationsTotal.WithLabelValues(string(notification.Type), channel).Inc()
	m.NotificationsByType.WithLabelValues(string(notification.Type)).Inc()
	m.NotificationsByUser.WithLabelValues(notification.UserID.String(), "").Inc()
}

// RecordUserRegistration records user registration metrics
func (m *Metrics) RecordUserRegistration(user *core.User) {
	m.UsersTotal.WithLabelValues(string(user.Role)).Inc()
	m.UsersByRole.WithLabelValues(string(user.Role)).Inc()
	if user.BranchID != nil {
		m.UsersByBranch.WithLabelValues(user.BranchID.String(), "").Inc()
	}
}

// RecordBranchCreated records branch creation metrics
func (m *Metrics) RecordBranchCreated(branch *core.Branch) {
	m.BranchesTotal.WithLabelValues(branch.City, branch.Province).Inc()
	m.BranchesByCity.WithLabelValues(branch.City).Inc()
	m.BranchesByProvince.WithLabelValues(branch.Province).Inc()
	if branch.IsActive {
		m.ActiveBranches.WithLabelValues(branch.City, branch.Province).Inc()
	}
}

// RecordDatabaseConnection records database connection metrics
func (m *Metrics) RecordDatabaseConnection(state string, count float64) {
	m.DatabaseConnections.WithLabelValues(state).Set(count)
}

// RecordRedisConnection records Redis connection metrics
func (m *Metrics) RecordRedisConnection(state string, count float64) {
	m.RedisConnections.WithLabelValues(state).Set(count)
}

// RecordCacheHit records cache hit metrics
func (m *Metrics) RecordCacheHit(cacheType string, hitRate float64) {
	m.CacheHitRate.WithLabelValues(cacheType).Set(hitRate)
}

// RecordCacheMiss records cache miss metrics
func (m *Metrics) RecordCacheMiss(cacheType string, missRate float64) {
	m.CacheMissRate.WithLabelValues(cacheType).Set(missRate)
}
