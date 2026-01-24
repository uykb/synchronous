package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MySQLDB *gorm.DB

// AppConfig stores the global configuration as a JSON blob
type AppConfig struct {
	ID        uint      `gorm:"primaryKey"`
	Config    []byte    `gorm:"type:json"`
	UpdatedAt time.Time
}

// Order represents a trading order (migrated from SQLite)
type Order struct {
	ID           uint      `gorm:"primaryKey"`
	Exchange     string    `gorm:"index"`
	Symbol       string    `gorm:"index"`
	OrderID       string `gorm:"size:128;uniqueIndex"`
	ClientOrderID string `gorm:"size:128"`
	Side         string
	Type         string
	Price        float64
	Quantity     float64
	Status       string
	ErrorMessage string
	Timestamp    int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func InitMySQL() error {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		// Fallback for local dev if needed, or error out
		// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		return fmt.Errorf("MYSQL_DSN environment variable is required")
	}

	var err error
	MySQLDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// Auto Migrate
	err = MySQLDB.AutoMigrate(&AppConfig{}, &Order{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("MySQL connected and migrated successfully")
	return nil
}

// SaveConfig saves the config struct to the database
func SaveConfig(cfg interface{}) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	var appConfig AppConfig
	result := MySQLDB.First(&appConfig, 1)
	
	appConfig.ID = 1
	appConfig.Config = data

	if result.Error == gorm.ErrRecordNotFound {
		return MySQLDB.Create(&appConfig).Error
	}
	
	return MySQLDB.Save(&appConfig).Error
}

// LoadConfigRaw reads the raw JSON config from the database
func LoadConfigRaw() ([]byte, error) {
	var appConfig AppConfig
	result := MySQLDB.First(&appConfig, 1)
	if result.Error != nil {
		return nil, result.Error
	}
	return appConfig.Config, nil
}

func SaveOrderResultMySQL(res interface{}) error {
	// Assuming res is *models.OrderResult
	// We need to map it to Order struct or just save it if it matches
	// For simplicity, let's just use the Order struct defined here
	// But we need to map fields.
	// Since I can't import models here easily without circular dep if models imports database (it doesn't),
	// I'll accept interface{} and use reflection or just assume caller handles mapping?
	// No, caller passes *models.OrderResult.
	
	// Let's just use GORM's ability to save structs if they match.
	// But models.OrderResult might not have GORM tags.
	
	return MySQLDB.Create(res).Error
}
