package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startKafkaConsumer(gormDB *gorm.DB, kafkaUrl string) {
	// Kafka broker address
	broker := kafkaUrl

	// Kafka topic to consume from
	topic := "user-created"

	// Kafka consumer group
	groupID := "my-group"

	// Initialize a new Kafka reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  groupID,
		Topic:    topic,
		MaxWait:  500 * time.Millisecond,
		MaxBytes: 10e6,
		Dialer:   &kafka.Dialer{Timeout: 10 * time.Second},
		//Logger:   kafka.LoggerFunc(fmt.Printf),
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			fmt.Fprintf(os.Stderr, msg+"\n", args...)
		}),
	})

	// Create a context and channel for handling termination signals
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("Received termination signal. Shutting down...")
		cancel()
	}()

	// Start consuming messages
	for {
		select {
		case <-ctx.Done():
			// Context canceled, exit the loop
			r.Close()
			return
		default:
			// Read the next message
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				fmt.Printf("Error reading message: %v\n", err)
				if errors.Is(err, context.Canceled) {
					fmt.Println("Context canceled. Shutting down...")
					return
				}
				continue
			}

			// Process the message payload
			user := model.User{}

			if err := json.Unmarshal(msg.Value, &user); err != nil {
				fmt.Printf("Error unmarshalling JSON payload: %v\n", err)
				return
			}

			res := repo.Create(gormDB, &user)
			if res.Error != nil {
				println("error saving user from kafka: " + res.Error.Error())
			} else {
				// Commit the offset manually
				if err := r.CommitMessages(ctx, msg); err != nil {
					fmt.Printf("Error committing offset: %v\n", err)
				}
			}
		}
	}
}
