package main

import (
	"log"
	"service/internal/shared/config/config"
	"service/internal/shared/database"
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
	// Note: AutoMigrate will handle adding new columns and making existing columns nullable
	// for backward compatibility with multi-service refactoring
	if err := db.AutoMigrate(&model.ServiceOrder{}); err != nil {
		log.Fatalf("Failed to migrate ServiceOrder table: %v", err)
	}
	log.Println("✓ ServiceOrder table migrated")
	
	// Step 3.1: Migrate existing ServiceOrder data for backward compatibility
	migrateServiceOrderData(db)

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

	// Step 13: Create ServiceCategory table
	if err := db.AutoMigrate(&model.ServiceCategory{}); err != nil {
		log.Fatalf("Failed to migrate ServiceCategory table: %v", err)
	}
	log.Println("✓ ServiceCategory table migrated")

	// Step 14: Create ServiceCatalog table which depends on ServiceCategory
	if err := db.AutoMigrate(&model.ServiceCatalog{}); err != nil {
		log.Fatalf("Failed to migrate ServiceCatalog table: %v", err)
	}
	log.Println("✓ ServiceCatalog table migrated")

	// Step 15: Create ServiceProvider table which depends on User
	if err := db.AutoMigrate(&model.ServiceProvider{}); err != nil {
		log.Fatalf("Failed to migrate ServiceProvider table: %v", err)
	}
	log.Println("✓ ServiceProvider table migrated")

	// Step 16: Create ProviderService join table
	if err := db.AutoMigrate(&model.ProviderService{}); err != nil {
		log.Fatalf("Failed to migrate ProviderService table: %v", err)
	}
	log.Println("✓ ProviderService table migrated")

	// Step 17: Create LocationTracking table
	if err := db.AutoMigrate(&model.LocationTracking{}); err != nil {
		log.Fatalf("Failed to migrate LocationTracking table: %v", err)
	}
	log.Println("✓ LocationTracking table migrated")

	// Step 18: Create CurrentLocation table
	if err := db.AutoMigrate(&model.CurrentLocation{}); err != nil {
		log.Fatalf("Failed to migrate CurrentLocation table: %v", err)
	}
	log.Println("✓ CurrentLocation table migrated")

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
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_user_id ON service_orders(customer_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_branch_id ON service_orders(branch_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_courier_id ON service_orders(courier_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_status ON service_orders(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_created_at ON service_orders(created_at)")
	
	// New multi-service indexes for ServiceOrder
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_service_catalog_id ON service_orders(service_catalog_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_service_provider_id ON service_orders(service_provider_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_appointment_date ON service_orders(appointment_date)")

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

	// ServiceCategory indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_categories_is_active ON service_categories(is_active)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_categories_sort_order ON service_categories(sort_order)")

	// ServiceCatalog indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_catalogs_category_id ON service_catalogs(category_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_catalogs_is_active ON service_catalogs(is_active)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_catalogs_requires_appointment ON service_catalogs(requires_appointment)")

	// ServiceProvider indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_providers_user_id ON service_providers(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_providers_city ON service_providers(city)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_providers_province ON service_providers(province)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_providers_is_active ON service_providers(is_active)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_providers_is_verified ON service_providers(is_verified)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_providers_rating ON service_providers(rating)")

	// ProviderService indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_provider_services_provider_id ON provider_services(provider_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_provider_services_service_catalog_id ON provider_services(service_catalog_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_provider_services_is_active ON provider_services(is_active)")

	// LocationTracking indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_location_tracking_order_id ON location_tracking(order_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_location_tracking_user_id ON location_tracking(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_location_tracking_timestamp ON location_tracking(timestamp)")

	// CurrentLocation indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_current_locations_order_id ON current_locations(order_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_current_locations_user_id ON current_locations(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_current_locations_updated_at ON current_locations(updated_at)")

	// ServiceOrder new indexes for on-demand service
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_is_on_demand ON service_orders(is_on_demand)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_current_latitude ON service_orders(current_latitude)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_current_longitude ON service_orders(current_longitude)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_service_orders_last_location_update ON service_orders(last_location_update)")

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

// migrateServiceOrderData migrates existing ServiceOrder data for backward compatibility
// This function ensures that existing iPhone service orders continue to work
// by mapping iPhone fields to generic item fields
func migrateServiceOrderData(db *gorm.DB) {
	log.Println("Migrating existing ServiceOrder data for backward compatibility...")
	
	// Check if there are any existing orders
	var orderCount int64
	db.Model(&model.ServiceOrder{}).Count(&orderCount)
	if orderCount == 0 {
		log.Println("No existing orders to migrate")
		return
	}
	
	// Migrate iPhone fields to generic item fields for existing orders
	// This ensures backward compatibility while supporting new multi-service features
	result := db.Exec(`
		UPDATE service_orders 
		SET 
			item_model = COALESCE(item_model, iphone_model),
			item_color = COALESCE(item_color, iphone_color),
			item_serial = COALESCE(item_serial, iphone_imei),
			item_type = COALESCE(item_type, iphone_type, iphone_model),
			service_name = COALESCE(service_name, 'iPhone Service'),
			service_location = COALESCE(service_location, 
				(SELECT address FROM branches WHERE branches.id = service_orders.branch_id)
			)
		WHERE 
			(item_model IS NULL OR item_model = '') 
			AND (iphone_model IS NOT NULL AND iphone_model != '')
	`)
	
	if result.Error != nil {
		log.Printf("Warning: Error migrating ServiceOrder data: %v", result.Error)
	} else {
		log.Printf("✓ Migrated %d existing ServiceOrder records", result.RowsAffected)
	}
	
	// Set default service type for orders without service_catalog_id
	result = db.Exec(`
		UPDATE service_orders 
		SET service_type = COALESCE(service_type, 'other')
		WHERE service_type IS NULL OR service_type = ''
	`)
	
	if result.Error != nil {
		log.Printf("Warning: Error setting default service type: %v", result.Error)
	} else {
		log.Printf("✓ Set default service type for %d orders", result.RowsAffected)
	}
	
	log.Println("✓ ServiceOrder data migration completed")
}
