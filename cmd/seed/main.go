package main

import (
	"log"
	"service/internal/shared/config/config"
	"service/internal/shared/database"
	"service/internal/shared/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitPostgres()

	// Run seeder
	runSeeder()
}

func runSeeder() {
	db := database.DB

	// Clear existing data
	clearData(db)

	// Seed users
	seedUsers(db)

	// Seed branches
	seedBranches(db)

	// Seed service orders
	seedServiceOrders(db)

	// Seed payments
	seedPayments(db)

	// Seed notifications
	seedNotifications(db)

	log.Println("Database seeding completed successfully")
}

func clearData(db *gorm.DB) {
	// Delete in reverse order of dependencies
	db.Exec("DELETE FROM notifications")
	db.Exec("DELETE FROM payments")
	db.Exec("DELETE FROM service_orders")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM branches")

	log.Println("Existing data cleared")
}

func seedUsers(db *gorm.DB) {
	users := []*model.User{
		{
			ID:       uuid.New(),
			Name:     "Admin Central",
			Email:    "admin@iphoneservice.com",
			Phone:    "081234567890",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     model.RoleAdminPusat,
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Name:     "Branch Admin Jakarta",
			Email:    "admin.jakarta@iphoneservice.com",
			Phone:    "081234567891",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     model.RoleAdminCabang,
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Name:     "Cashier Jakarta",
			Email:    "cashier.jakarta@iphoneservice.com",
			Phone:    "081234567892",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     model.RoleKasir,
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Name:     "Technician Jakarta",
			Email:    "technician.jakarta@iphoneservice.com",
			Phone:    "081234567893",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     model.RoleTeknisi,
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Name:     "Courier Jakarta",
			Email:    "courier.jakarta@iphoneservice.com",
			Phone:    "081234567894",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     model.RoleKurir,
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Name:     "John Doe",
			Email:    "john.doe@example.com",
			Phone:    "081234567895",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     model.RolePelanggan,
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Name:     "Jane Smith",
			Email:    "jane.smith@example.com",
			Phone:    "081234567896",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     model.RolePelanggan,
			IsActive: true,
		},
	}

	for _, user := range users {
		if err := db.Create(user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Name, err)
		}
	}

	log.Println("Users seeded successfully")
}

func seedBranches(db *gorm.DB) {
	branches := []*model.Branch{
		{
			ID:        uuid.New(),
			Name:      "Jakarta Central",
			Address:   "Jl. Sudirman No. 123, Jakarta Selatan",
			City:      "Jakarta",
			Province:  "DKI Jakarta",
			Phone:     "021-12345678",
			Latitude:  -6.2088,
			Longitude: 106.8456,
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Name:      "Surabaya Branch",
			Address:   "Jl. Tunjungan No. 456, Surabaya",
			City:      "Surabaya",
			Province:  "Jawa Timur",
			Phone:     "031-87654321",
			Latitude:  -7.2575,
			Longitude: 112.7521,
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Name:      "Bandung Branch",
			Address:   "Jl. Asia Afrika No. 789, Bandung",
			City:      "Bandung",
			Province:  "Jawa Barat",
			Phone:     "022-11223344",
			Latitude:  -6.9175,
			Longitude: 107.6191,
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Name:      "Medan Branch",
			Address:   "Jl. Imam Bonjol No. 321, Medan",
			City:      "Medan",
			Province:  "Sumatera Utara",
			Phone:     "061-55667788",
			Latitude:  3.5952,
			Longitude: 98.6722,
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Name:      "Semarang Branch",
			Address:   "Jl. Pemuda No. 654, Semarang",
			City:      "Semarang",
			Province:  "Jawa Tengah",
			Phone:     "024-99887766",
			Latitude:  -6.9667,
			Longitude: 110.4167,
			IsActive:  true,
		},
	}

	for _, branch := range branches {
		if err := db.Create(branch).Error; err != nil {
			log.Printf("Failed to create branch %s: %v", branch.Name, err)
		}
	}

	log.Println("Branches seeded successfully")
}

func seedServiceOrders(db *gorm.DB) {
	// Get user and branch IDs
	var users []model.User
	var branches []model.Branch
	db.Find(&users)
	db.Find(&branches)

	if len(users) == 0 || len(branches) == 0 {
		log.Println("No users or branches found, skipping service orders")
		return
	}

	iphoneType1 := "iPhone 14 Pro"
	iphoneType2 := "iPhone 13"
	iphoneType3 := "iPhone 12"
	pickupLoc1 := "Jakarta Selatan"
	pickupLoc2 := "Surabaya"
	pickupLoc3 := "Jakarta Pusat"
	orders := []*model.ServiceOrder{
		{
			ID:                uuid.New(),
			CustomerID:        users[5].ID,    // John Doe
			BranchID:          &branches[0].ID, // Jakarta Central
			OrderNumber:       "ORD-20240101-001",
			IPhoneType:        &iphoneType1,
			Complaint:         "Screen cracked, needs replacement",
			PickupLocation:    &pickupLoc1,
			Status:            model.StatusInService,
			ServiceCost:       500000,
			EstimatedDuration: 3,
			CreatedAt:         time.Now().AddDate(0, 0, -5),
			UpdatedAt:         time.Now().AddDate(0, 0, -2),
		},
		{
			ID:                uuid.New(),
			CustomerID:        users[6].ID,    // Jane Smith
			BranchID:          &branches[1].ID, // Surabaya Branch
			OrderNumber:       "ORD-20240102-002",
			IPhoneType:        &iphoneType2,
			Complaint:         "Battery draining fast, needs replacement",
			PickupLocation:    &pickupLoc2,
			Status:            model.StatusReady,
			ServiceCost:       300000,
			EstimatedDuration: 2,
			CreatedAt:         time.Now().AddDate(0, 0, -3),
			UpdatedAt:         time.Now().AddDate(0, 0, -1),
		},
		{
			ID:                uuid.New(),
			CustomerID:        users[5].ID,    // John Doe
			BranchID:          &branches[0].ID, // Jakarta Central
			OrderNumber:       "ORD-20240103-003",
			IPhoneType:        &iphoneType3,
			Complaint:         "Camera not working, needs repair",
			PickupLocation:    &pickupLoc3,
			Status:            model.StatusCompleted,
			ServiceCost:       400000,
			EstimatedDuration: 4,
			CreatedAt:         time.Now().AddDate(0, 0, -7),
			UpdatedAt:         time.Now().AddDate(0, 0, -1),
		},
	}

	for _, order := range orders {
		if err := db.Create(order).Error; err != nil {
			log.Printf("Failed to create order %s: %v", order.OrderNumber, err)
		}
	}

	log.Println("Service orders seeded successfully")
}

func seedPayments(db *gorm.DB) {
	// Get order IDs
	var orders []model.ServiceOrder
	db.Find(&orders)

	if len(orders) == 0 {
		log.Println("No orders found, skipping payments")
		return
	}

	payments := []*model.Payment{
		{
			ID:            uuid.New(),
			OrderID:       orders[0].ID,
			UserID:        orders[0].CustomerID,
			Amount:        orders[0].ServiceCost,
			PaymentMethod: model.PaymentMethodMidtrans,
			Status:        model.PaymentStatusPaid,
			TransactionID: "TXN-001",
			InvoiceNumber: "INV-20240101-001",
			PaidAt:        &time.Time{},
			CreatedAt:     time.Now().AddDate(0, 0, -5),
			UpdatedAt:     time.Now().AddDate(0, 0, -4),
		},
		{
			ID:            uuid.New(),
			OrderID:       orders[1].ID,
			UserID:        orders[1].CustomerID,
			Amount:        orders[1].ServiceCost,
			PaymentMethod: model.PaymentMethodGopay,
			Status:        model.PaymentStatusPaid,
			TransactionID: "TXN-002",
			InvoiceNumber: "INV-20240102-002",
			PaidAt:        &time.Time{},
			CreatedAt:     time.Now().AddDate(0, 0, -3),
			UpdatedAt:     time.Now().AddDate(0, 0, -2),
		},
		{
			ID:            uuid.New(),
			OrderID:       orders[2].ID,
			UserID:        orders[2].CustomerID,
			Amount:        orders[2].ServiceCost,
			PaymentMethod: model.PaymentMethodCash,
			Status:        model.PaymentStatusPaid,
			TransactionID: "TXN-003",
			InvoiceNumber: "INV-20240103-003",
			PaidAt:        &time.Time{},
			CreatedAt:     time.Now().AddDate(0, 0, -7),
			UpdatedAt:     time.Now().AddDate(0, 0, -6),
		},
	}

	for _, payment := range payments {
		if err := db.Create(payment).Error; err != nil {
			log.Printf("Failed to create payment %s: %v", payment.InvoiceNumber, err)
		}
	}

	log.Println("Payments seeded successfully")
}

func seedNotifications(db *gorm.DB) {
	// Get user and order IDs
	var users []model.User
	var orders []model.ServiceOrder
	db.Find(&users)
	db.Find(&orders)

	if len(users) == 0 || len(orders) == 0 {
		log.Println("No users or orders found, skipping notifications")
		return
	}

	notifications := []*model.Notification{
		{
			ID:        uuid.New(),
			UserID:    users[5].ID, // John Doe
			OrderID:   &orders[0].ID,
			Type:      model.NotificationTypeOrderUpdate,
			Message:   "Your order ORD-20240101-001 is now in service",
			IsRead:    false,
			CreatedAt: time.Now().AddDate(0, 0, -2),
			UpdatedAt: time.Now().AddDate(0, 0, -2),
		},
		{
			ID:        uuid.New(),
			UserID:    users[6].ID, // Jane Smith
			OrderID:   &orders[1].ID,
			Type:      model.NotificationTypeOrderReady,
			Message:   "Your order ORD-20240102-002 is ready for pickup",
			IsRead:    true,
			CreatedAt: time.Now().AddDate(0, 0, -1),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		},
		{
			ID:        uuid.New(),
			UserID:    users[5].ID, // John Doe
			OrderID:   &orders[2].ID,
			Type:      model.NotificationTypeOrderCompleted,
			Message:   "Your order ORD-20240103-003 has been completed",
			IsRead:    true,
			CreatedAt: time.Now().AddDate(0, 0, -1),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		},
	}

	for _, notification := range notifications {
		if err := db.Create(notification).Error; err != nil {
			log.Printf("Failed to create notification: %v", err)
		}
	}

	log.Println("Notifications seeded successfully")
}
