package database

import (
	"fmt"
	"pdf-parser/internal/config"
	"pdf-parser/internal/model"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitMySQL(cfg *config.DatabaseConfig) error {
	var err error
	once.Do(func() {
		db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			err = fmt.Errorf("failed to get database instance: %w", err)
			return
		}

		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

		if err = autoMigrate(); err != nil {
			err = fmt.Errorf("failed to migrate database: %w", err)
			return
		}
	})
	return err
}

func autoMigrate() error {
	return db.AutoMigrate(
		&model.ParseTask{},
		&model.ParseResult{},
		&model.User{},
		&model.Session{},
		&model.CodeExchange{},
	)
}

func AutoMigrate() error {
	return autoMigrate()
}

func GetDB() *gorm.DB {
	return db
}

func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
