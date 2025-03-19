package util

import (
	"fmt"
	"log"
	"time"

	"github.com/m7medVision/crime-management-system/internal/auth"
	"github.com/m7medVision/crime-management-system/internal/config"
	"github.com/m7medVision/crime-management-system/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(config *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// Set connection pool parameters
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate the schema
	err = db.AutoMigrate(
		&model.User{},
		&model.Case{},
		&model.Report{},
		&model.Evidence{},
		&model.Suspect{},
		&model.Victim{},
		&model.Witness{},
		&model.AuditLog{},
		&model.Comment{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	// Create default admin user if not exists
	var adminCount int64
	db.Model(&model.User{}).Where("role = ?", model.RoleAdmin).Count(&adminCount)
	if adminCount == 0 {
		adminPassword, err := auth.HashPassword("admin123") // In production, use a secure password
		if err != nil {
			log.Fatalf("Failed to hash default admin password: %v", err)
		}

		adminUser := model.User{
			Username:       "admin",
			Password:       adminPassword,
			Email:          "admin@districtcore.gov",
			FullName:       "System Administrator",
			Role:           model.RoleAdmin,
			ClearanceLevel: model.ClearanceCritical,
			IsActive:       true,
		}
		result := db.Create(&adminUser)
		if result.Error != nil {
			log.Fatalf("Failed to create default admin user: %v", result.Error)
		}
		log.Println("Default admin user created successfully")
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
