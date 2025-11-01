package main

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Test Database
	fmt.Printf("Testing PostgreSQL connection using: %s\n", os.Getenv("DATABASE_URL"))
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		fmt.Printf("❌ PostgreSQL connection error: %v\n", err)
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("❌ PostgreSQL DB() error: %v\n", err)
		} else {
			err = sqlDB.Ping()
			if err != nil {
				fmt.Printf("❌ PostgreSQL ping error: %v\n", err)
			} else {
				fmt.Println("✅ PostgreSQL connection successful")
			}
		}
	}

	// Test Redis
	fmt.Printf("Testing Redis connection using: %s\n", os.Getenv("REDIS_URL"))
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		fmt.Printf("❌ Redis URL parse error: %v\n", err)
		return
	}

	rdb := redis.NewClient(opt)
	defer rdb.Close()

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Printf("❌ Redis ping error: %v\n", err)
	} else {
		fmt.Println("✅ Redis connection successful")
	}
}
