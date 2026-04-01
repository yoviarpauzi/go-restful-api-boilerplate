package database

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(config *viper.Viper, log *zap.Logger) *gorm.DB {
	username := config.GetString("DB_USER")
	password := config.GetString("DB_PASSWORD")
	host := config.GetString("DB_HOST")
	port := config.GetInt("DB_PORT")
	database := config.GetString("DB_NAME")
	sslMode := config.GetString("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", host, username, password, database, port, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(config.GetInt("DB_MAX_IDLE_CONNS"))
	sqlDB.SetMaxOpenConns(config.GetInt("DB_MAX_OPEN_CONNS"))
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("DB_CONN_MAX_LIFETIME")) * time.Minute)

	return db
}
