package database

import (
	"fmt"
	"github.com/genryusaishigikuni/ch4nge/config"
	"github.com/genryusaishigikuni/ch4nge/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// Get database configuration
	dsn := config.AppConfig.GetDatabaseDSN()

	// Configure GORM logger based on environment
	var logLevel logger.LogLevel
	if config.AppConfig.Server.Env == "production" {
		logLevel = logger.Error
	} else {
		logLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// Connect to database with retry logic
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 5)
		}
	}

	if err != nil {
		log.Fatal("Failed to connect to database after retries:", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")

	// Auto migrate tables
	err = DB.AutoMigrate(
		&models.User{}, &models.UserFriend{}, &models.Achievement{}, &models.UserAchievement{},
		&models.MiniChallenge{}, &models.UserMiniChallenge{}, &models.WeeklyChallenge{},
		&models.UserWeeklyChallenge{}, &models.Activity{}, &models.Action{}, &models.Post{},
		&models.PostLike{}, &models.PostShare{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// Create admin user if not exists
	createAdminUser()
}

func createAdminUser() {
	var count int64
	DB.Model(&models.User{}).Where("is_admin = ?", true).Count(&count)
	if count == 0 {
		// Use environment variable for admin password in production
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" {
			adminPassword = "admin123"
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash admin password:", err)
		}

		admin := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			IsAdmin:  true,
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Fatal("Failed to create admin user:", err)
		}

		fmt.Printf("Admin user created: admin@example.com / %s\n", adminPassword)
	}
}
