package db

import (
	"context"
	"fmt"
	"github.com/GGmaz/wallet-arringo/config"
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPg(cfg config.DBConfigPg) (*gorm.DB, error) {
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
		&model.Transaction{},
		&model.User{},
		&model.Account{},
	} // Add all tables here
	gormDB.AutoMigrate(tables...)
	return gormDB, err
}

func InitRedis(cfg config.DBConfigRedis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Host + ":" + strconv.Itoa(cfg.Port),
		DB:   0, // use default DB
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
