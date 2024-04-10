package server

import (
	"fmt"
	"github.com/GGmaz/wallet-arringo/config"
	"github.com/GGmaz/wallet-arringo/internal/db"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	v1 "github.com/GGmaz/wallet-arringo/internal/server/api"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"log"
)

type Server struct {
	config *config.Config
}

func dbMiddleware(gormDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("transaction", gormDB)
		c.Next()
	}
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	r := gin.Default()
	r.Use(gin.Logger())

	gormDB, err := db.Init(server.config.Db)
	if err != nil {
		log.Fatal("Could not connect to the database" + err.Error())
		return
	}

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Use(dbMiddleware(gormDB))
	r.Use(CORSMiddleware())

	//// Create a NATS connection
	////nc, err := nats.Connect(nats.DefaultURL)
	//nc, err := nats.Connect(server.config.NatsUrl)
	//if err != nil {
	//	log.Fatal("Failed to connect to NATS" + err.Error())
	//	return
	//}
	//startNats(nc, gormDB)
	//defer nc.Close()
	//
	//go startKafkaConsumer(gormDB, server.config.KafkaUrl)

	v1.RegisterVersion(r)

	err = r.Run(":" + server.config.Port)

	if err != nil {
		log.Fatal("Could not start the server" + err.Error())
		return
	}

	println("Starting server on port: " + server.config.Port)
}

func startNats(nc *nats.Conn, gormDB *gorm.DB) {
	log.Println("Connected to " + nats.DefaultURL)

	_, err := nc.Subscribe("balance-request", func(m *nats.Msg) {
		// Parse the user email from the request parameters
		userEmail := string(m.Data)
		// Get users balance
		balance, err := repo.GetBalanceForUserMail(gormDB, userEmail)
		if err != nil {
			println("Could not get balance")
			return
		}

		balanceBytes := []byte(fmt.Sprintf("%.2f", balance))

		err = nc.Publish(m.Reply, balanceBytes)
		if err != nil {
			println("Could not publish balance")
			return
		}
	})
	if err != nil {
		println("Could not subscribe to balance-request")
		return
	}
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
