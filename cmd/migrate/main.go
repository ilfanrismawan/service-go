package main

import (
	"log"
	"service/internal/shared/config/config"
	"service/internal/shared/database/database"
	"service/internal/shared/model"

	"gorm.io/gorm"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitPostgres()

	// Run migrations
	runMigrations()
}

func runMigrations() {
	db := database.DB

	log.Println("Running migrations...")

	// Step 1: Create Branch table first since it has no dependencies
	if err := db.AutoMigrate(&model.Branch{}); err != nil {
		log.Fatalf("Failed to migrate Branch table: %v", err)
	}
	log.Println("✓ Branch table migrated")

	// Step 2: Create User table which depends on Branch
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to migrate User table: %v", err)
	}
	log.Println("✓ User table migrated")

	// Step 3: Create ServiceOrder table which depends on User and Branch
	if err := db.AutoMigrate(&model.ServiceOrder{}); err != nil {
		log.Fatalf("Failed to migrate ServiceOrder table: %v", err)
	}
	log.Println("✓ ServiceOrder table migrated")

	// Step 4: Create Payment which depends on ServiceOrder
	if err := db.AutoMigrate(&model.Payment{}); err != nil {
		log.Fatalf("Failed to migrate Payment table: %v", err)
	}
	log.Println("✓ Payment table migrated")

	// Step 5: Create Notification which depends on User and ServiceOrder
	if err := db.AutoMigrate(&model.Notification{}); err != nil {
		log.Fatalf("Failed to migrate Notification table: %v", err)
	}
	log.Println("✓ Notification table migrated")

	// Step 6: Create Membership which depends on User
	if err := db.AutoMigrate(&model.Membership{}); err != nil {
		log.Fatalf("Failed to migrate Membership table: %v", err)
	}
	log.Println("✓ Membership table migrated")

	// Step 7: Create ChatMessage which depends on User and ServiceOrder
	if err := db.AutoMigrate(&model.ChatMessage{}); err != nil {
		log.Fatalf("Failed to migrate ChatMessage table: %v", err)
	}
	log.Println("✓ ChatMessage table migrated")

	// Step 8: Create Queue table
	if err := db.AutoMigrate(&model.Queue{}); err != nil {
		log.Fatalf("Failed to migrate Queue table: %v", err)
	}
	log.Println("✓ Queue table migrated")

	// Step 9: Create Warranty table
	if err := db.AutoMigrate(&model.Warranty{}); err != nil {
		log.Fatalf("Failed to migrate Warranty table: %v", err)
	}
	log.Println("✓ Warranty table migrated")

	// Step 10: Create SparePartInventory table
	if err := db.AutoMigrate(&model.SparePartInventory{}); err != nil {
		log.Fatalf("Failed to migrate SparePartInventory table: %v", err)
	}
	log.Println("✓ SparePartInventory table migrated")

	// Step 11: Create Rating table
	if err := db.AutoMigrate(&model.Rating{}); err != nil {
		log.Fatalf("Failed to migrate Rating table: %v", err)
	}
	log.Println("✓ Rating table migrated")

	// Step 12: Create AuditTrail table
	if err := db.AutoMigrate(&model.AuditTrail{}); err != nil {
		log.Fatalf("Failed to migrate AuditTrail table: %v", err)
	}
	log.Println("✓ AuditTrail table migrated")

	// Create indexes
	createIndexes(db)

	// Seed initial data
	seedInitialData(db)

	log.Println("✅ Database migrations completed successfully")
}

func createIndexes(db *gorm.DB) {
	// User indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_branch_id ON users(branch_id)")

	// Branch indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_branches_city ON branches(city)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_branches_province ON branches(province)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_branches_is_active ON branches(is_active)")

	// Service order indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_user_id ON service_orders(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_branch_id ON service_orders(branch_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_courier_id ON service_orders(courier_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_status ON service_orders(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_created_at ON service_orders(created_at)")

	// Payment indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payments_payment_method ON payments(payment_method)")

	// Notification indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_order_id ON notifications(order_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read)")

	// Membership indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_memberships_user_id ON memberships(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_memberships_tier ON memberships(tier)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_memberships_status ON memberships(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_memberships_total_spent ON memberships(total_spent)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_memberships_joined_at ON memberships(joined_at)")

	// Queue indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_queues_branch_id ON queues(branch_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_queues_customer_id ON queues(customer_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_queues_status ON queues(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_queues_created_at ON queues(created_at)")

	// Warranty indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_warranties_order_id ON warranties(order_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_warranties_end_date ON warranties(end_date)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_warranties_is_active ON warranties(is_active)")

	// SparePartInventory indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_spare_parts_branch_id ON spare_part_inventory(branch_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_spare_parts_part_code ON spare_part_inventory(part_code)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_spare_parts_stock ON spare_part_inventory(stock)")

	// Rating indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ratings_order_id ON ratings(order_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ratings_customer_id ON ratings(customer_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ratings_branch_id ON ratings(branch_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ratings_technician_id ON ratings(technician_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ratings_rating ON ratings(rating)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ratings_created_at ON ratings(created_at)")

	// AuditTrail indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_trails_user_id ON audit_trails(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_trails_resource ON audit_trails(resource)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_trails_resource_id ON audit_trails(resource_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_trails_action ON audit_trails(action)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_trails_created_at ON audit_trails(created_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_trails_ip_address ON audit_trails(ip_address)")

	log.Println("Database indexes created successfully")
}

func seedInitialData(db *gorm.DB) {
	// Check if data already exists
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Initial data already exists, skipping seed")
		return
	}

	// Create default admin user
	adminUser := &model.User{
		FullName: "Admin",
		Email:    "admin@iphoneservice.com",
		Phone:    "081234567890",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Role:     model.RoleAdminPusat,
		IsActive: true,
	}

	if err := db.Create(adminUser).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
	}

	// Create sample branches
	branches := []*model.Branch{
		{
			Name:      "Jakarta Central",
			Address:   "Jl. Sudirman No. 123",
			City:      "Jakarta",
			Province:  "DKI Jakarta",
			Phone:     "021-12345678",
			Latitude:  -6.2088,
			Longitude: 106.8456,
			IsActive:  true,
		},
		{
			Name:      "Surabaya Branch",
			Address:   "Jl. Tunjungan No. 456",
			City:      "Surabaya",
			Province:  "Jawa Timur",
			Phone:     "031-87654321",
			Latitude:  -7.2575,
			Longitude: 112.7521,
			IsActive:  true,
		},
		{
			Name:      "Bandung Branch",
			Address:   "Jl. Asia Afrika No. 789",
			City:      "Bandung",
			Province:  "Jawa Barat",
			Phone:     "022-11223344",
			Latitude:  -6.9175,
			Longitude: 107.6191,
			IsActive:  true,
		},
	}

	for _, branch := range branches {
		if err := db.Create(branch).Error; err != nil {
			log.Printf("Failed to create branch %s: %v", branch.Name, err)
		}
	}

	log.Println("Initial data seeded successfully")
}
