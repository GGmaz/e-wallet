package db

import (
	"context"
	"fmt"
	"github.com/GGmaz/wallet-arringo/config"
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg config.DBConfig) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)
	sqlDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, _ := sqlDB.DB()
	// Set connection pool options
	db.SetMaxOpenConns(20) // Maximum number of open connections
	db.SetMaxIdleConns(10) // Maximum number of idle connections

	newLogger.Info(context.Background(), "aaa")
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{Logger: newLogger})
	fmt.Println("Created ", gormDB)

	tables := []interface{}{
		//TODO: Add all tables here
		&model.Transaction{},
		&model.User{},
	} // Add all tables here
	gormDB.AutoMigrate(tables...)
	return gormDB, err
}
