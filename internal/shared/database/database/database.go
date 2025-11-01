package database

import (
	"context"
	"log"
	"service/internal/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

// InitPostgres initializes PostgreSQL connection
func InitPostgres() {
	var err error

	log.Printf("Attempting to connect to PostgreSQL with URL: %s\n", config.Config.DatabaseURL)

	// Configure GORM logger based on environment
	var gormLogger logger.Interface
	if config.Config.Environment == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Connect to PostgreSQL with enhanced configuration
	DB, err = gorm.Open(postgres.Open(config.Config.DatabaseURL), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DryRun:                                   false,
	})
	if err != nil {
		log.Printf("❌ Failed to connect to PostgreSQL with error: %v\n", err)
		log.Fatal("Database connection failed")
	}

	// Verify connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("❌ Failed to get underlying *sql.DB: %v\n", err)
		log.Fatal("Database connection verification failed")
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("❌ Failed to ping database: %v\n", err)
		log.Fatal("Database ping failed")
	}

	log.Println("✅ Connected to PostgreSQL successfully")

	// Create uuid-ossp extension if it doesn't exist
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Printf("❌ Failed to create uuid-ossp extension: %v\n", err)
		log.Fatal("Extension creation failed")
	}
	log.Println("✅ UUID extension enabled")

	// Database connection successful
	log.Println("✅ Database connected and initialized successfully")
}

// InitRedis initializes Redis connection
func InitRedis() {
	log.Printf("Attempting to connect to Redis with URL: %s\n", config.Config.RedisURL)

	opt, err := redis.ParseURL(config.Config.RedisURL)
	if err != nil {
		log.Printf("❌ Failed to parse Redis URL: %v\n", err)
		log.Fatal("Redis URL parsing failed")
	}

	Redis = redis.NewClient(opt)

	// Test Redis connection with timeout context
	ctx := context.Background()
	_, err = Redis.Ping(ctx).Result()
	if err != nil {
		log.Printf("❌ Failed to ping Redis: %v\n", err)
		log.Fatal("Redis connection failed")
	}

	// Try setting and getting a test value
	testKey := "connection_test"
	err = Redis.Set(ctx, testKey, "ok", 0).Err()
	if err != nil {
		log.Printf("❌ Failed to set test value in Redis: %v\n", err)
		log.Fatal("Redis write test failed")
	}

	_, err = Redis.Get(ctx, testKey).Result()
	if err != nil {
		log.Printf("❌ Failed to get test value from Redis: %v\n", err)
		log.Fatal("Redis read test failed")
	}

	Redis.Del(ctx, testKey)
	log.Println("✅ Connected to Redis successfully")
}

// CloseDatabase closes database connections
func CloseDatabase() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	if Redis != nil {
		Redis.Close()
	}
}
