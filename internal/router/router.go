package router

import (
	"net/http"
	branchHandler "service/internal/branches/handler"
	orderHandler "service/internal/orders/handler"
	paymentHandler "service/internal/payments/handler"
	sharedHandlers "service/internal/shared/handlers"
	"service/internal/shared/middleware"
	"service/internal/shared/model"
	"service/internal/users/handler"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all API routes
func SetupRoutes(r *gin.Engine) {
	// Initialize domain handlers
	authHandler := handler.NewAuthHandler()
	branchHdlr := branchHandler.NewBranchHandler()
	orderHdlr := orderHandler.NewOrderHandler()
	paymentHdlr := paymentHandler.NewPaymentHandler()

	// Initialize shared handlers
	notificationHandler := sharedHandlers.NewNotificationHandler()
	chatHandler := sharedHandlers.NewChatHandler()
	dashboardHandler := sharedHandlers.NewDashboardHandler()
	fileHandler := sharedHandlers.NewFileHandler()
	healthHandler := sharedHandlers.NewHealthHandler()
	wsHandler := sharedHandlers.NewWebSocketHandler()
	swaggerHandler := sharedHandlers.NewSwaggerHandler()
	membershipHandler := sharedHandlers.NewMembershipHandler()
	reportHandler := sharedHandlers.NewReportHandler()
	ratingHandler := sharedHandlers.NewRatingHandler()

	// Setup Swagger documentation routes
	swaggerHandler.SetupSwaggerRoutes(r)

	// Health check endpoint (no auth required)
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/health/live", healthHandler.LivenessCheck)
	r.GET("/health/ready", healthHandler.ReadinessCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := v1.Group("/")
		{
			// Authentication routes (rate limited)
			authPublic := public.Group("/auth")
			authPublic.Use(middleware.RateLimitMiddleware())
			{
				authPublic.POST("/register", authHandler.Register)
				authPublic.POST("/login", authHandler.Login)
				authPublic.POST("/refresh", authHandler.RefreshToken)
				authPublic.POST("/logout", authHandler.Logout)
				authPublic.POST("/forgot-password", authHandler.ForgotPassword)
				authPublic.POST("/reset-password", authHandler.ResetPassword)
			}

			// Payment callbacks (public, signature verified in handler)
			public.POST("/payments/midtrans/callback", paymentHdlr.MidtransCallback)

			// Public branch information
			public.GET("/branches", branchHdlr.GetBranches)
			public.GET("/branches/nearest", branchHdlr.GetNearestBranches)
			public.GET("/branches/:id", branchHdlr.GetBranch)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile routes
			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.PUT("/auth/profile", authHandler.UpdateProfile)
			protected.POST("/auth/change-password", authHandler.ChangePassword)
			protected.PUT("/auth/fcm-token", authHandler.UpdateFCMToken)

			// Order routes
			protected.POST("/orders", orderHdlr.CreateOrder)
			protected.GET("/orders", orderHdlr.GetOrders)
			protected.GET("/orders/:id", orderHdlr.GetOrder)
			protected.PUT("/orders/:id/status", orderHdlr.UpdateOrderStatus)
			protected.PUT("/orders/:id/assign-courier", orderHdlr.AssignCourier)
			protected.PUT("/orders/:id/assign-technician", orderHdlr.AssignTechnician)

			// Payment routes
			protected.POST("/payments/create-invoice", paymentHdlr.CreateInvoice)
			protected.POST("/payments/process", paymentHdlr.ProcessPayment)
			protected.GET("/payments/:id", paymentHdlr.GetPayment)
			protected.GET("/payments/order/:orderId", paymentHdlr.GetPaymentsByOrder)

			// Notification routes
			protected.GET("/notifications", notificationHandler.GetNotifications)
			protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
			protected.POST("/notifications", notificationHandler.SendNotification)
			protected.POST("/notifications/order/:orderId/status", notificationHandler.SendOrderStatusNotification)
			protected.POST("/notifications/order/:orderId/payment", notificationHandler.SendPaymentNotification)

			// File upload routes
			protected.POST("/files/upload", fileHandler.UploadFile)
			protected.POST("/files/orders/photo", fileHandler.UploadOrderPhoto)
			protected.POST("/files/users/avatar", fileHandler.UploadUserAvatar)
			protected.GET("/files/url", fileHandler.GetFileURL)
			protected.GET("/files/list", fileHandler.ListFiles)
			protected.DELETE("/files/delete", fileHandler.DeleteFile)

			// Chat routes
			protected.GET("/chat/orders/:orderId", chatHandler.GetChatMessages)
			protected.POST("/chat/orders/:orderId", chatHandler.SendMessage)

			// Dashboard routes
			protected.GET("/dashboard/overview", dashboardHandler.GetOverview)
			protected.GET("/dashboard/orders", dashboardHandler.GetOrderStats)
			protected.GET("/dashboard/revenue", dashboardHandler.GetRevenueStats)
			protected.GET("/dashboard/branches", dashboardHandler.GetBranchStats)

			// Membership routes
			protected.GET("/membership", membershipHandler.GetMembership)
			protected.POST("/membership", membershipHandler.CreateMembership)
			protected.PUT("/membership", membershipHandler.UpdateMembership)
			protected.POST("/membership/redeem-points", membershipHandler.RedeemPoints)
			protected.POST("/membership/subscribe", membershipHandler.SubscribeToMembership)
			protected.POST("/membership/cancel", membershipHandler.CancelSubscription)
			protected.POST("/membership/trial", membershipHandler.StartTrial)
			protected.GET("/membership/tiers", membershipHandler.GetMembershipTiers)
			protected.POST("/membership/upgrade", membershipHandler.UpgradeMembership)
			protected.GET("/membership/usage", membershipHandler.GetMembershipUsage)

			// Report routes
			protected.GET("/reports/current-month", reportHandler.GetCurrentMonthReport)
			protected.GET("/reports/monthly", reportHandler.GetMonthlyReport)
			protected.GET("/reports/yearly", reportHandler.GetYearlyReport)
			protected.GET("/reports/summary", reportHandler.GetReportSummary)

			// Rating routes
			protected.POST("/ratings", ratingHandler.CreateRating)
			protected.GET("/ratings", ratingHandler.ListRatings)
			protected.GET("/ratings/average", ratingHandler.GetAverageRating)
			protected.GET("/ratings/:id", ratingHandler.GetRating)
			protected.PUT("/ratings/:id", ratingHandler.UpdateRating)
			protected.DELETE("/ratings/:id", ratingHandler.DeleteRating)
		}

		// Admin routes (admin role required)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.RoleMiddleware(model.RoleAdminPusat, model.RoleAdminCabang))
		{
			// Branch management
			admin.POST("/branches", branchHdlr.CreateBranch)
			admin.PUT("/branches/:id", branchHdlr.UpdateBranch)
			admin.DELETE("/branches/:id", branchHdlr.DeleteBranch)
			admin.GET("/branches", branchHdlr.GetBranches)

			// User management
			admin.GET("/users", authHandler.GetUsers)
			admin.GET("/users/:id", authHandler.GetUser)
			admin.PUT("/users/:id", authHandler.UpdateUser)
			admin.DELETE("/users/:id", authHandler.DeleteUser)

			// Order management
			admin.GET("/orders", orderHdlr.GetAllOrders)
			admin.PUT("/orders/:id", orderHdlr.UpdateOrder)
			admin.DELETE("/orders/:id", orderHdlr.DeleteOrder)

			// Payment management
			admin.GET("/payments", paymentHdlr.GetAllPayments)
			admin.PUT("/payments/:id", paymentHdlr.UpdatePayment)

			// Dashboard admin
			admin.GET("/dashboard", dashboardHandler.GetAdminDashboard)

			// Membership management
			admin.GET("/membership/list", membershipHandler.ListMemberships)
			admin.GET("/membership/stats", membershipHandler.GetMembershipStats)
			admin.GET("/membership/top-spenders", membershipHandler.GetTopSpenders)
		}

		// Cashier routes (kasir role required)
		cashier := v1.Group("/cashier")
		cashier.Use(middleware.AuthMiddleware())
		cashier.Use(middleware.RoleMiddleware(model.RoleKasir))
		{
			// Order processing
			cashier.GET("/orders", orderHdlr.GetCashierOrders)
			cashier.PUT("/orders/:id/status", orderHdlr.UpdateOrderStatus)
			cashier.POST("/orders/:id/payment", paymentHdlr.ProcessPayment)

			// Branch orders
			cashier.GET("/branches/:id/orders", orderHdlr.GetBranchOrders)
		}

		// Technician routes (teknisi role required)
		technician := v1.Group("/technician")
		technician.Use(middleware.AuthMiddleware())
		technician.Use(middleware.RoleMiddleware(model.RoleTeknisi))
		{
			// Order management
			technician.GET("/orders", orderHdlr.GetTechnicianOrders)
			technician.PUT("/orders/:id/status", orderHdlr.UpdateOrderStatus)
			technician.POST("/orders/:id/photo", fileHandler.UploadOrderPhoto)

			// Chat
			technician.GET("/chat/orders/:orderId", chatHandler.GetChatMessages)
			technician.POST("/chat/orders/:orderId", chatHandler.SendMessage)
		}

		// Courier routes (kurir role required)
		courier := v1.Group("/courier")
		courier.Use(middleware.AuthMiddleware())
		courier.Use(middleware.RoleMiddleware(model.RoleKurir))
		{
			// Order management
			courier.GET("/orders", orderHdlr.GetCourierOrders)
			courier.PUT("/orders/:id/status", orderHdlr.UpdateOrderStatus)
			courier.POST("/orders/:id/photo", fileHandler.UploadOrderPhoto)

			// Available jobs
			courier.GET("/jobs", orderHdlr.GetAvailableJobs)
			courier.POST("/jobs/:id/accept", orderHdlr.AcceptJob)
		}
	}

	// WebSocket routes
	r.GET("/ws/chat", wsHandler.HandleWebSocket)

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.CreateErrorResponse(
			"not_found",
			"Endpoint not found",
			nil,
		))
	})
}
