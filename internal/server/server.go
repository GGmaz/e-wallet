package server

import (
	"bytes"
	"github.com/GGmaz/wallet-arringo/config"
	"github.com/GGmaz/wallet-arringo/docs"
	"github.com/GGmaz/wallet-arringo/internal/db"
	"github.com/GGmaz/wallet-arringo/internal/helper"
	"github.com/GGmaz/wallet-arringo/internal/scheduler"
	v1 "github.com/GGmaz/wallet-arringo/internal/server/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

type Server struct {
	config *config.Config
}

func pgMiddleware(gormDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("transaction", gormDB)
		c.Next()
	}
}

func redisMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", redisClient)
		scheduler.StartDataCollector(c, redisClient)
		c.Next()
	}
}

func hmacValidationMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if request url starts with /swagger
		if c.Request.URL.Path[:len("/swagger")] == "/swagger" || c.Request.ContentLength == 0 {
			c.Next()
			return
		}

		// Read request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer c.Request.Body.Close()

		// Calculate hash of the request body
		calculatedHash := helper.Sha512HmacHash(body, secretKey)

		// Retrieve hash from the header
		receivedHash := c.GetHeader("X-Authorization-Sign")
		if receivedHash == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing X-Authorization-Sign header"})
			return
		}

		// Compare calculated hash with received hash
		if calculatedHash != receivedHash {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Hash verification failed"})
			return
		}

		// Set the request body reader back for further processing
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		// Continue processing the request
		c.Next()
	}
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

// @title e-wallet API
// @version 1.0
// @description This is a project for Arringo company.
// @host localhost:8082
// @BasePath /api/v1
func (server *Server) Start() {
	r := gin.Default()
	r.Use(gin.Logger())

	gormDB, err := db.InitPg(server.config.DbPg)
	if err != nil {
		log.Fatal("Could not connect to Postgres" + err.Error())
		return
	}

	redisClient, err := db.InitRedis(server.config.DbRedis)
	if err != nil {
		log.Fatal("Could not connect to Redis" + err.Error())
		return
	}

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Use(pgMiddleware(gormDB), redisMiddleware(redisClient))
	r.Use(hmacValidationMiddleware(server.config.SecretKey))
	r.Use(CORSMiddleware())

	v1.RegisterVersion(r)

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = r.Run(":" + server.config.Port)

	if err != nil {
		log.Fatal("Could not start the server" + err.Error())
		return
	}

	println("Starting server on port: " + server.config.Port)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
