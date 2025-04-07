package database

import (
	"fmt"

	"github.com/0xirvan/url-shortener/config"
	"github.com/0xirvan/url-shortener/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName,
	)

	dialect := mysql.Open(dsn)

	db, err := gorm.Open(dialect, &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})

	if err != nil {
		utils.Log.Errorf("Failed to connect to database: %+v", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		utils.Log.Errorf("Failed to get database instance: %+v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0)

	return db
}
